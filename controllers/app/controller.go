package app

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/registry"
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
	"k8s.io/client-go/util/retry"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
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

var log logr.Logger

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
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgApp == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
		discoverOcpDefaultRouteHost(r.Client)
	}

	// setup finalizer
	if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			cnvrgApp.ObjectMeta.Finalizers = append(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(ctx, cnvrgApp); err != nil {
				log.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusRemoving, Message: "removing cnvrg spec"}, cnvrgApp)
			if err := r.cleanup(cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				if err := r.Update(ctx, cnvrgApp); err != nil {
					cnvrgApp, err := r.getCnvrgAppSpec(req.NamespacedName)
					if err != nil {
						log.Error(err, "error getting cnvrgapp for finalizer cleanup")
						return err
					}
					cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
					return r.Update(ctx, cnvrgApp)
				}
				return err
			})
			if err != nil {
				log.Info("error in removing finalizer")
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
	if err := r.applyManifests(cnvrgApp); err != nil {
		return ctrl.Result{}, err
	}

	// get control plan readiness
	ready, percentageReady, stackReadiness, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}
	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", percentageReady)
	log.Info(statusMsg)

	if ready { // ura, done
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

func (r *CnvrgAppReconciler) applyManifests(cnvrgApp *mlopsv1.CnvrgApp) error {
	// registry
	log.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgApp.Namespace,
		Data: map[string]interface{}{
			"Registry":    cnvrgApp.Spec.Registry,
			"Annotations": cnvrgApp.Spec.Annotations,
			"Labels":      cnvrgApp.Spec.Labels,
		},
	}
	// registry
	if err := desired.Apply(registry.State(registryData), cnvrgApp, r.Client, r.Scheme, log); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// dbs
	if err := r.dbsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// networking
	log.Info("applying networking")
	if err := desired.Apply(networking.AppNetworkingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, log); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// controlplane
	log.Info("applying controlplane")
	if err := desired.Apply(controlplane.State(cnvrgApp), cnvrgApp, r.Client, r.Scheme, log); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.Status, app *mlopsv1.CnvrgApp) {

	if app.Status.Status == mlopsv1.StatusRemoving {
		log.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	if status.Status == mlopsv1.StatusError {
		msg := fmt.Sprintf("%s/%s error acoured during reconcile", app.GetNamespace(), app.GetName())
		r.recorder.Event(app, "Warning", "ReconcileError", msg)
	}
	ctx := context.Background()
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		name := types.NamespacedName{Namespace: app.Namespace, Name: app.Name}
		app, err := r.getCnvrgAppSpec(name)
		if err != nil {
			return err
		}
		app.Status.Status = status.Status
		app.Status.Message = status.Message
		if status.Progress >= 0 {
			app.Status.Progress = status.Progress
		}
		if status.StackReadiness != nil {
			app.Status.StackReadiness = status.StackReadiness
		}
		err = r.Status().Update(ctx, app)
		return err
	})
	if err != nil {
		log.Error(err, "can't update status")
	}

}

func (r *CnvrgAppReconciler) syncCnvrgAppSpec(name types.NamespacedName) (bool, error) {

	log.Info("synchronizing cnvrgApp spec")

	// Fetch current cnvrgApp spec
	cnvrgApp, err := r.getCnvrgAppSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgApp == nil {
		return false, nil // probably cnvrgapp was removed
	}
	log = r.Log.WithValues("name", name, "ns", cnvrgApp.Namespace)

	// Get default cnvrgApp spec
	desiredSpec := mlopsv1.DefaultCnvrgAppSpec()

	if err := CalculateAndApplyAppDefaults(cnvrgApp, &desiredSpec, r.Client); err != nil {
		log.Error(err, "can't calculate defaults")
		return false, err
	}

	// Merge current cnvrgApp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgApp.Spec, mergo.WithOverride, mergo.WithTransformers(cnvrgSpecBoolTransformer{})); err != nil {
		log.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {

		if diff, equal := messagediff.PrettyDiff(desiredSpec, cnvrgApp.Spec); !equal {
			log.Info("diff between desiredSpec and actual")
			log.Info(diff)
		}

		if diff, equal := messagediff.PrettyDiff(cnvrgApp.Spec, desiredSpec); !equal {
			log.Info("diff between actual and desired")
			log.Info(diff)
		}

	}

	equal := reflect.DeepEqual(desiredSpec, cnvrgApp.Spec)
	if !equal {
		log.Info("states are not equals, syncing and requeuing")
		cnvrgApp.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgApp); err != nil && errors.IsConflict(err) {
			log.Info("conflict updating cnvrgApp object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	log.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgAppReconciler) getCnvrgAppSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var app mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			log.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		log.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}

	return &app, nil
}

func (r *CnvrgAppReconciler) cleanup(cnvrgApp *mlopsv1.CnvrgApp) error {

	log.Info("running finalizer cleanup")

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
		log.Info("cleanup-pvc is false, skipping pvc deletion!")
		return nil
	}
	log.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		log.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
			if pvc.ObjectMeta.Labels["app"] == "prometheus" || pvc.ObjectMeta.Labels["app"] == "elasticsearch" {
				if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
					log.Info("pvc already deleted")
				} else if err != nil {
					log.Error(err, "error deleting prometheus pvc")
					return err
				}
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanupDbInitCm(desiredSpec *mlopsv1.CnvrgApp) error {
	log.Info("running cnvrg-db-init cleanup")
	ctx := context.Background()
	dbInitCm := &v1core.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-db-init", Namespace: desiredSpec.Namespace}}
	err := r.Delete(ctx, dbInitCm)
	if err != nil && errors.IsNotFound(err) {
		log.Info("no need to delete cnvrg-db-init, cm not found")
	} else {
		log.Error(err, "error deleting cnvrg-db-init")
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log = r.Log.WithValues("initializing", "crds")

	appPredicate := predicate.Funcs{

		CreateFunc: func(createEvent event.CreateEvent) bool {
			msg := fmt.Sprintf("cnvrgapp: %s/%s has been created", createEvent.Object.GetNamespace(), createEvent.Object.GetName())
			r.recorder.Event(createEvent.Object, "Normal", "Created", msg)
			return true
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
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
			log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
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
				if containsString(healthCheckWorkloads, e.ObjectNew.GetName()) {
					return true
				}
			}
			log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return true
		},
	}
	r.recorder = mgr.GetEventRecorderFor("cnvrgapp")
	cnvrgAppController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgApp{}, builder.WithPredicates(appPredicate))

	for _, v := range desired.Kinds {

		if strings.Contains(v.Group, "istio.io") && viper.GetBool("own-istio-resources") == false {
			continue
		}
		if strings.Contains(v.Group, "openshift.io") && viper.GetBool("own-openshift-resources") == false {
			continue
		}
		if strings.Contains(v.Group, "coreos.com") && viper.GetBool("own-prometheus-resources") == false {
			continue
		}
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(v)
		cnvrgAppController.Owns(u, builder.WithPredicates(appOwnsPredicate))
	}

	log.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
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
