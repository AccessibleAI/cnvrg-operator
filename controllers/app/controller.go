package app

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/dbs"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/registry"
	sso2 "github.com/AccessibleAI/cnvrg-operator/pkg/app/sso"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"gopkg.in/d4l3k/messagediff.v1"
	v1apps "k8s.io/api/apps/v1"
	v1batch "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"
)

const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"

type CnvrgAppReconciler struct {
	client.Client
	recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log = r.Log.WithValues("name", req.NamespacedName)
	log.Info("starting cnvrgapp reconciliation")

	// sync specs between actual and defaults
	equal, err := r.syncCnvrgAppSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		return ctrl.Result{Requeue: true}, nil // specs are not equals -> reconcile
	}

	// specs are synced, proceed reconcile
	cnvrgApp, err := r.getCnvrgAppSpec(req.NamespacedName)
	b, _ := json.Marshal(cnvrgApp)
	log.Info(string(b))
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgApp == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
		// add finalizer
		if !controllerutil.ContainsFinalizer(cnvrgApp, CnvrgappFinalizer) {
			controllerutil.AddFinalizer(cnvrgApp, CnvrgappFinalizer)
			if err := r.Update(ctx, cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else { // run finalizer on delete
		if controllerutil.ContainsFinalizer(cnvrgApp, CnvrgappFinalizer) {
			if err := r.cleanup(cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(cnvrgApp, CnvrgappFinalizer)
			if err := r.Update(ctx, cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	// check if enabled control plane workloads are all in ready status
	ready, percentageReady, stackReadiness, err := r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	// even if all control plane workloads are ready, let operator finish the full reconcile loop
	if percentageReady == 100 {
		percentageReady = 99
	}
	s := mlopsv1.Status{
		Status:         mlopsv1.StatusReconciling,
		Message:        fmt.Sprintf("reconciling... (%d%%)", percentageReady),
		Progress:       percentageReady,
		StackReadiness: stackReadiness}
	r.updateStatusMessage(s, cnvrgApp)

	// apply spec manifests
	if err := r.apply(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return ctrl.Result{}, err
	}

	// get control plan readiness
	ready, percentageReady, stackReadiness, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}
	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", percentageReady)
	log.Info(statusMsg)

	if ready { // u r done, done
		s := mlopsv1.Status{
			Status:         mlopsv1.StatusReady,
			Message:        statusMsg,
			Progress:       percentageReady,
			StackReadiness: stackReadiness}
		r.updateStatusMessage(s, cnvrgApp)
		log.Info("stack is ready!")
		r.recorder.Event(cnvrgApp, "Normal", "Created", fmt.Sprintf("cnvrgapp %s successfully deployed", req.NamespacedName))
		return ctrl.Result{}, nil
	} else { // reconcile again
		requeueAfter, _ := time.ParseDuration("30s")
		log.Info("stack not ready yet, requeuing...")
		r.recorder.Event(cnvrgApp, "Normal", "Creating", fmt.Sprintf("cnvrgapp %s not ready yet, done: %d%%", req.NamespacedName, percentageReady))
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
}

func (r *CnvrgAppReconciler) apply(app *mlopsv1.CnvrgApp) error {

	if err := registry.NewRegistryStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
		return err
	}

	if app.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress && app.Spec.Networking.Ingress.IstioGwEnabled {
		if err := networking.NewIstioGwState(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
	}

	if app.Spec.Dbs.Pg.Enabled {
		if err := dbs.NewPgStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
	}

	if app.Spec.Dbs.Minio.Enabled {
		if err := dbs.NewMinioStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
	}

	if app.Spec.Dbs.Es.Enabled {
		if err := dbs.NewElasticStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
		if app.Spec.Dbs.Es.Kibana.Enabled {
			if err := dbs.NewKibanaStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}
		if app.Spec.Dbs.Es.Elastalert.Enabled {
			if err := dbs.NewElastAlertStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}
	}

	if app.Spec.Dbs.Redis.Enabled {
		if err := dbs.NewRedisStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
	}

	if app.Spec.Dbs.Prom.Enabled {
		if err := dbs.NewPromStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
		if app.Spec.Dbs.Prom.Grafana.Enabled {
			if err := dbs.NewGrafanaStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}
	}

	if app.Spec.SSO.Enabled {

		if app.Spec.SSO.Pki.Enabled {
			if err := sso2.NewPkiStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}

		if app.Spec.SSO.Jwks.Enabled {
			if err := sso2.NewJwksStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}

		if app.Spec.SSO.Central.Enabled {
			if err := sso2.NewCentralStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}

		if app.Spec.SSO.Proxy.Enabled {
			if err := sso2.NewProxyStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
				return err
			}
		}

	}

	if err := controlplane.NewControlPlaneStateManager(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.Status, app *mlopsv1.CnvrgApp) {

	if app.Status.Status == mlopsv1.StatusRemoving {
		r.Log.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	if status.Status == mlopsv1.StatusError {
		msg := fmt.Sprintf("%s/%s error acoured during reconcile", app.GetNamespace(), app.GetName())
		r.recorder.Event(app, "Warning", "ReconcileError", msg)
	}
	app.Status.Status = status.Status
	app.Status.Message = status.Message
	if status.Progress >= 0 {
		app.Status.Progress = status.Progress
	}
	if status.StackReadiness != nil {
		app.Status.StackReadiness = status.StackReadiness
	}
	if err := r.Status().Update(context.Background(), app); err != nil {
		r.recorder.Event(app, "Warning", "StatusUpdateError", err.Error())
	}
}

func (r *CnvrgAppReconciler) syncCnvrgAppSpec(name types.NamespacedName) (bool, error) {

	r.Log.Info("synchronizing app spec")

	// Fetch current app spec
	app, err := r.getCnvrgAppSpec(name)
	if err != nil {
		return false, err
	}
	if app == nil {
		return false, nil // probably cnvrgapp was removed
	}
	r.Log.WithValues("name", name, "ns", app.Namespace)

	// Get default app spec
	desiredSpec := mlopsv1.DefaultCnvrgAppSpec()

	if err := CalculateAndApplyAppDefaults(app, &desiredSpec, r.Client); err != nil {
		r.Log.Error(err, "can't calculate defaults")
		return false, err
	}

	// Merge current app spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, app.Spec, mergo.WithOverride, mergo.WithTransformers(controllers.CnvrgSpecBoolTransformer{})); err != nil {
		r.Log.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {

		if diff, equal := messagediff.PrettyDiff(desiredSpec, app.Spec); !equal {
			r.Log.Info("diff between desiredSpec and actual")
			r.Log.Info(diff)
		}

	}

	equal := reflect.DeepEqual(desiredSpec, app.Spec)
	if !equal {
		r.Log.Info("states are not equals, syncing and requeuing")
		app.Spec = desiredSpec
		if err := r.Update(context.Background(), app); err != nil && errors.IsConflict(err) {
			r.Log.Info("conflict updating app object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	r.Log.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgAppReconciler) getCnvrgAppSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var app mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		r.Log.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}

	return &app, nil
}

func (r *CnvrgAppReconciler) cleanup(cnvrgApp *mlopsv1.CnvrgApp) error {

	r.Log.Info("running finalizer cleanup")

	// remove cnvrg-db-init
	if err := r.cleanupDbInitCm(cnvrgApp); err != nil {
		return err
	}

	// cleanup pvc
	if err := r.cleanupPVCs(); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) cleanupPVCs() error {
	if !viper.GetBool("cleanup-pvc") {
		r.Log.Info("cleanup-pvc is false, skipping pvc deletion!")
		return nil
	}
	r.Log.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		r.Log.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
			if pvc.ObjectMeta.Labels["app"] == "prometheus" || pvc.ObjectMeta.Labels["app"] == "elasticsearch" {
				if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
					r.Log.Info("pvc already deleted")
				} else if err != nil {
					r.Log.Error(err, "error deleting prometheus pvc")
					return err
				}
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanupDbInitCm(desiredSpec *mlopsv1.CnvrgApp) error {
	r.Log.Info("running cnvrg-db-init cleanup")
	ctx := context.Background()
	dbInitCm := &v1core.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-db-init", Namespace: desiredSpec.Namespace}}
	err := r.Delete(ctx, dbInitCm)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("no need to delete cnvrg-db-init, cm not found")
	} else {
		r.Log.Error(err, "error deleting cnvrg-db-init")
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if viper.GetBool("create-crds") {
		if err := controlplane.NewControlPlaneCrdsStateManager(r.Client, r.Scheme, r.Log).Apply(); err != nil {
			return err
		}
	}

	appPredicate := predicate.Funcs{

		CreateFunc: func(createEvent event.CreateEvent) bool {
			msg := fmt.Sprintf("cnvrgapp: %s/%s has been created", createEvent.Object.GetNamespace(), createEvent.Object.GetName())
			r.recorder.Event(createEvent.Object, "Normal", "Created", msg)
			return true
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			r.Log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			shouldReconcile := e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
			if shouldReconcile {
				msg := fmt.Sprintf("cnvrgapp: %s/%s has been updated", e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
				r.recorder.Event(e.ObjectNew, "Normal", "Updated", msg)
			}
			return shouldReconcile
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			msg := fmt.Sprintf("cnvrgapp: %s/%s has been deleted", deleteEvent.Object.GetNamespace(), deleteEvent.Object.GetName())
			r.recorder.Event(deleteEvent.Object, "Normal", "SuccessfulDelete", msg)
			r.Log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return !deleteEvent.DeleteStateUnknown
		},
	}

	appOwnsPredicate := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {
			gvk := e.ObjectNew.GetObjectKind().GroupVersionKind()
			app := mlopsv1.DefaultCnvrgAppSpec()
			healthCheckWorkloads := []string{
				"ingresscheck",
				"sidekiq",
				"searchkiq",
				"systemkiq",
				app.ControlPlane.WebApp.SvcName,
				app.Dbs.Pg.SvcName,
				app.Dbs.Minio.SvcName,
				app.Dbs.Redis.SvcName,
				app.Dbs.Es.SvcName,
			}
			if gvk == desired.Kinds[desired.DeploymentGVK] || gvk == desired.Kinds[desired.StatefulSetGVK] || gvk == desired.Kinds[desired.JobGVK] {
				if controllers.ContainsString(healthCheckWorkloads, e.ObjectNew.GetName()) {
					return true
				}
			}
			r.Log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			r.Log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return true
		},
	}
	r.recorder = mgr.GetEventRecorderFor("cnvrgapp")
	cnvrgAppController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgApp{}, builder.WithPredicates(appPredicate))

	for _, v := range desired.Kinds {

		if strings.Contains(v.Group, "istio.io") {
			continue
		}
		if strings.Contains(v.Group, "openshift.io") {
			continue
		}
		if strings.Contains(v.Group, "coreos.com") {
			continue
		}
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(v)
		cnvrgAppController.Owns(u, builder.WithPredicates(appOwnsPredicate))
	}

	r.Log.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgAppController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}

func (r *CnvrgAppReconciler) CheckJobReadiness(name types.NamespacedName) (bool, error) {
	ctx := context.Background()
	job := &v1batch.Job{}

	if err := r.Get(ctx, name, job); err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if job.Status.Succeeded >= 1 {
		return true, nil
	}

	return false, nil
}

func (r *CnvrgAppReconciler) CheckDeploymentReadiness(name types.NamespacedName) (bool, error) {
	ctx := context.Background()
	deployment := &v1apps.Deployment{}

	if err := r.Get(ctx, name, deployment); err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}

func (r *CnvrgAppReconciler) CheckStatefulSetReadiness(name types.NamespacedName) (bool, error) {

	ctx := context.Background()
	sts := &v1apps.StatefulSet{}

	if err := r.Get(ctx, name, sts); err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if sts.Status.Replicas == sts.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}
