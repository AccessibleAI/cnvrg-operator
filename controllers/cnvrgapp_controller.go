package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/ingresscheck"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
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
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	"strconv"
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

var appLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	appLog = r.Log.WithValues("name", req.NamespacedName)
	appLog.Info("starting cnvrgapp reconciliation")

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
				appLog.Error(err, "failed to add finalizer")
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
						appLog.Error(err, "error getting cnvrgapp for finalizer cleanup")
						return err
					}
					cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
					return r.Update(ctx, cnvrgApp)
				}
				return err
			})
			if err != nil {
				appLog.Info("error in removing finalizer")
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
	appLog.Info(statusMsg)

	if ready { // ura, done
		s := mlopsv1.Status{
			Status:         mlopsv1.StatusReady,
			Message:        statusMsg,
			Progress:       percentageReady,
			StackReadiness: stackReadiness}
		r.updateStatusMessage(s, cnvrgApp)
		appLog.Info("stack is ready!")
		r.recorder.Event(cnvrgApp, "Normal", "Created", fmt.Sprintf("cnvrgapp %s successfully deployed", req.NamespacedName))
		return ctrl.Result{}, nil
	} else { // reconcile again
		requeueAfter, _ := time.ParseDuration("30s")
		appLog.Info("stack not ready yet, requeuing...")
		r.recorder.Event(cnvrgApp, "Normal", "Creating", fmt.Sprintf("cnvrgapp %s not ready yet, done: %d%%", req.NamespacedName, percentageReady))
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
}

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(cnvrgApp *mlopsv1.CnvrgApp) (bool, int, map[string]bool, error) {

	readyState := make(map[string]bool)

	// check ingresscheck status
	if cnvrgApp.Spec.IngressCheck.Enabled {
		name := types.NamespacedName{Name: "ingresscheck", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckJobReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["ingressCheck"] = ready
	}

	// check webapp status
	if cnvrgApp.Spec.ControlPlane.WebApp.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.ControlPlane.WebApp.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["webApp"] = ready
	}

	// check sidekiq status
	if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled {
		name := types.NamespacedName{Name: "sidekiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["sidekiq"] = ready
	}

	// check searchkiq status
	if cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled {
		name := types.NamespacedName{Name: "searchkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["searchkiq"] = ready
	}

	// check systemkiq status
	if cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled {
		name := types.NamespacedName{Name: "systemkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["systemkiq"] = ready
	}

	// check postgres status
	if cnvrgApp.Spec.Dbs.Pg.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Pg.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["pg"] = ready
	}

	// check minio status
	if cnvrgApp.Spec.Dbs.Minio.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Minio.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["minio"] = ready
	}

	// check redis status
	if cnvrgApp.Spec.Dbs.Redis.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Redis.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["redis"] = ready
	}

	// check es status
	if cnvrgApp.Spec.Dbs.Es.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Es.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		// if es is ready, trigger fluentbit reconfiguration
		if ready {
			appLog.Info("es is ready, triggering fluentbit reconfiguration")
			if err := r.addFluentbitConfiguration(cnvrgApp); err != nil {
				return false, 0, nil, err
			}
		}
		readyState["es"] = ready
	}

	// check kibana status
	if cnvrgApp.Spec.Logging.Kibana.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Logging.Kibana.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["kibana"] = ready
	}

	// check prometheus status
	if cnvrgApp.Spec.Monitoring.Prometheus.Enabled {
		name := types.NamespacedName{Name: "prometheus-cnvrg-ccp-prometheus", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["prometheus"] = ready
	}

	// check prometheus status
	if cnvrgApp.Spec.Monitoring.Grafana.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Monitoring.Grafana.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["grafana"] = ready
	}

	percentageReady := 0

	readyCount := 0

	for _, ready := range readyState {
		if ready {
			readyCount++
		}
	}

	if len(readyState) > 0 {
		percentageReady = readyCount * 100 / len(readyState)
	}

	return readyCount == len(readyState), percentageReady, readyState, nil
}

func (r *CnvrgAppReconciler) applyManifests(cnvrgApp *mlopsv1.CnvrgApp) error {
	// registry
	appLog.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgApp.Namespace,
		Data: map[string]interface{}{
			"Registry":    cnvrgApp.Spec.Registry,
			"Annotations": cnvrgApp.Spec.Annotations,
			"Labels":      cnvrgApp.Spec.Labels,
		},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// dbs
	if err := r.dbsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// backups
	if err := r.backupsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// networking
	appLog.Info("applying networking")
	if err := desired.Apply(networking.CnvrgAppNetworkingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// logging
	if err := r.loggingState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// controlplane
	appLog.Info("applying controlplane")
	if err := desired.Apply(controlplane.State(cnvrgApp), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// monitoring
	if err := r.monitoringState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// ingress check
	if err := r.ingressCheckState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) ingressCheckState(app *mlopsv1.CnvrgApp) error {
	// apply ingress check state
	if err := desired.Apply(ingresscheck.IngressCheckState(app), app, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) getCnvrgInfra() (*mlopsv1.CnvrgInfra, error) {

	cnvrgAppInfra := &mlopsv1.CnvrgInfraList{}

	if err := r.List(context.Background(), cnvrgAppInfra); err != nil {
		appLog.Error(err, "can't list CnvrgInfra objects")
		return nil, err
	}

	if len(cnvrgAppInfra.Items) == 0 {
		appLog.Info("no CnvrgInfra objects was deployed, skipping infra reconciler")
		return nil, errors.NewNotFound(schema.GroupResource{Group: "mlops.cnvrg.io", Resource: "CnvrgInfra"}, "cnvrg-infra")
	}

	return &cnvrgAppInfra.Items[0], nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.Status, app *mlopsv1.CnvrgApp) {

	if app.Status.Status == mlopsv1.StatusRemoving {
		appLog.Info("skipping status update, current cnvrg spec under removing status...")
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
		appLog.Error(err, "can't update status")
	}

}

func (r *CnvrgAppReconciler) syncCnvrgAppSpec(name types.NamespacedName) (bool, error) {

	appLog.Info("synchronizing cnvrgApp spec")

	// Fetch current cnvrgApp spec
	cnvrgApp, err := r.getCnvrgAppSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgApp == nil {
		return false, nil // probably cnvrgapp was removed
	}
	appLog = r.Log.WithValues("name", name, "ns", cnvrgApp.Namespace)

	// Get default cnvrgApp spec
	desiredSpec := mlopsv1.DefaultCnvrgAppSpec()

	infra, err := r.getCnvrgInfra()
	if err != nil && !errors.IsNotFound(err) {
		appLog.Error(err, "can't get cnvrg infra")
		//return false, err
	}

	if err := CalculateAndApplyAppDefaults(cnvrgApp, &desiredSpec, infra, r.Client); err != nil {
		appLog.Error(err, "can't calculate defaults")
		return false, err
	}

	// Merge current cnvrgApp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgApp.Spec, mergo.WithOverride, mergo.WithTransformers(cnvrgSpecBoolTransformer{})); err != nil {
		appLog.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {

		if diff, equal := messagediff.PrettyDiff(desiredSpec, cnvrgApp.Spec); !equal {
			appLog.Info("diff between desiredSpec and actual")
			appLog.Info(diff)
		}

		if diff, equal := messagediff.PrettyDiff(cnvrgApp.Spec, desiredSpec); !equal {
			appLog.Info("diff between actual and desired")
			appLog.Info(diff)
		}

	}

	equal := reflect.DeepEqual(desiredSpec, cnvrgApp.Spec)
	if !equal {
		appLog.Info("states are not equals, syncing and requeuing")
		cnvrgApp.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgApp); err != nil && errors.IsConflict(err) {
			appLog.Info("conflict updating cnvrgApp object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	appLog.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgAppReconciler) getCnvrgAppSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var app mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			appLog.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		appLog.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}

	return &app, nil
}

func (r *CnvrgAppReconciler) cleanup(cnvrgApp *mlopsv1.CnvrgApp) error {

	appLog.Info("running finalizer cleanup")

	// remove cnvrg-db-init
	if err := r.cleanupDbInitCm(cnvrgApp); err != nil {
		return err
	}

	// update infra reconciler cm
	if err := r.removeFluentbitConfiguration(cnvrgApp); err != nil {
		if err.Error() == "no CnvrgInfra objects was deployed, skipping infra reconciler" {
			appLog.Info("cnvrgInfra object not found, no need to trigger infra reconciler")
		} else {
			return err
		}
	}

	// cleanup pvc
	if err := r.cleanupPVCs(); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) cleanupPVCs() error {
	if !viper.GetBool("cleanup-pvc") {
		appLog.Info("cleanup-pvc is false, skipping pvc deletion!")
		return nil
	}
	appLog.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		appLog.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
			if pvc.ObjectMeta.Labels["app"] == "prometheus" || pvc.ObjectMeta.Labels["app"] == "elasticsearch" {
				if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
					appLog.Info("pvc already deleted")
				} else if err != nil {
					appLog.Error(err, "error deleting prometheus pvc")
					return err
				}
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanupDbInitCm(desiredSpec *mlopsv1.CnvrgApp) error {
	appLog.Info("running cnvrg-db-init cleanup")
	ctx := context.Background()
	dbInitCm := &v1core.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-db-init", Namespace: desiredSpec.Namespace}}
	err := r.Delete(ctx, dbInitCm)
	if err != nil && errors.IsNotFound(err) {
		appLog.Info("no need to delete cnvrg-db-init, cm not found")
	} else {
		appLog.Error(err, "error deleting cnvrg-db-init")
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	appLog = r.Log.WithValues("initializing", "crds")

	appPredicate := predicate.Funcs{

		CreateFunc: func(createEvent event.CreateEvent) bool {
			msg := fmt.Sprintf("cnvrgapp: %s/%s has been created", createEvent.Object.GetNamespace(), createEvent.Object.GetName())
			r.recorder.Event(createEvent.Object, "Normal", "Created", msg)
			return true
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			infraLog.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
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
			infraLog.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
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
			infraLog.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			infraLog.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
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

	appLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
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

func (r *CnvrgAppReconciler) ApplyCapsuleAnnotations(b mlopsv1.Backup, pvc *v1core.PersistentVolumeClaim, serviceType string) error {
	if pvc.Annotations == nil {
		pvc.Annotations = map[string]string{}
	}
	pvc.Annotations["capsule.mlops.cnvrg.io/backup"] = "false"
	if b.Enabled {
		pvc.Annotations["capsule.mlops.cnvrg.io/backup"] = "true"
	}
	pvc.Annotations["capsule.mlops.cnvrg.io/serviceType"] = serviceType
	pvc.Annotations["capsule.mlops.cnvrg.io/bucketRef"] = b.BucketRef
	pvc.Annotations["capsule.mlops.cnvrg.io/credsRef"] = b.CredsRef
	pvc.Annotations["capsule.mlops.cnvrg.io/rotation"] = strconv.Itoa(b.Rotation)
	pvc.Annotations["capsule.mlops.cnvrg.io/period"] = b.Period

	if err := r.Update(context.Background(), pvc); err != nil {
		return err
	}
	return nil
}
