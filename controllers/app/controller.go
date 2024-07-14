package app

import (
	"context"
	"dario.cat/mergo"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/dbs"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/registry"
	sso2 "github.com/AccessibleAI/cnvrg-operator/pkg/app/sso"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	"gopkg.in/d4l3k/messagediff.v1"
	v1apps "k8s.io/api/apps/v1"
	v1batch "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"time"
)

const systemStatusHealthCheckLabelName = "cnvrg-system-status-check"
const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"
const RolloutAnnotation = "kubectl.kubernetes.io/restartedAt"

type CnvrgAppReconciler struct {
	client.Client
	recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

// +kubebuilder:rbac:groups=mlops.cnvrg.io,namespace=cnvrg,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,namespace=cnvrg,resources=cnvrgapps/status,verbs=get;update;patch

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
	stackReadiness, err := r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	// even if all control plane workloads are ready, let operator finish the full reconcile loop
	if stackReadiness.percentageReady == 100 {
		stackReadiness.percentageReady = 99
	}

	//check are feature flags were updated
	var featureFlagsChanged bool
	newFeatureFlagsHash := hashStringsMap(cnvrgApp.Spec.ControlPlane.BaseConfig.FeatureFlags)

	lastFeatureFlagsHash := cnvrgApp.Status.LastFeatureFlagsHash

	// feature flags are changed
	if newFeatureFlagsHash != lastFeatureFlagsHash && lastFeatureFlagsHash != "" {
		featureFlagsChanged = true
	}

	s := mlopsv1.Status{
		Status:               mlopsv1.StatusReconciling,
		Message:              fmt.Sprintf("reconciling... (%d%%)", stackReadiness.percentageReady),
		Progress:             stackReadiness.percentageReady,
		StackReadiness:       stackReadiness.readyState,
		LastFeatureFlagsHash: newFeatureFlagsHash,
	}
	r.updateStatusMessage(s, cnvrgApp)

	// apply spec manifests
	if err := r.apply(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return ctrl.Result{}, err
	}

	// get control plan readiness
	stackReadiness, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	var rollingOnFeatureFlagUpdate bool

	//feature flags were changed, therefore need rollout scheduler
	if featureFlagsChanged {
		rollingOnFeatureFlagUpdate = true

		if cnvrgApp.Spec.ControlPlane.WebApp.Enabled {
			err := r.RollDeployment(types.NamespacedName{Name: cnvrgApp.Spec.ControlPlane.WebApp.SvcName, Namespace: cnvrgApp.Namespace})
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled {
			err = r.RollDeployment(types.NamespacedName{Name: sidekiq, Namespace: cnvrgApp.Namespace})
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		if cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled {
			err = r.RollDeployment(types.NamespacedName{Name: searchkiq, Namespace: cnvrgApp.Namespace})
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		if cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled {
			err = r.RollDeployment(types.NamespacedName{Name: systemkiq, Namespace: cnvrgApp.Namespace})
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		if cnvrgApp.Spec.ControlPlane.CnvrgScheduler.Enabled {
			err = r.RollDeployment(types.NamespacedName{Name: scheduler, Namespace: cnvrgApp.Namespace})
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", stackReadiness.percentageReady)
	log.Info(statusMsg)

	if stackReadiness.isReady && !rollingOnFeatureFlagUpdate { // u r done and no need to roll due to feature flags change, done
		s := mlopsv1.Status{
			Status:               mlopsv1.StatusReady,
			Message:              statusMsg,
			Progress:             stackReadiness.percentageReady,
			StackReadiness:       stackReadiness.readyState,
			LastFeatureFlagsHash: newFeatureFlagsHash,
		}
		r.updateStatusMessage(s, cnvrgApp)
		log.Info("stack is ready!")
		r.recorder.Event(cnvrgApp, "Normal", "Created", fmt.Sprintf("cnvrgapp %s successfully deployed", req.NamespacedName))
		return ctrl.Result{}, nil
	} else { // reconcile again
		requeueAfter, _ := time.ParseDuration("30s")
		var logMessage, eventMessage string
		if rollingOnFeatureFlagUpdate {
			logMessage = "rolling apps due to feature flags change..."
			eventMessage = "rolling: feature flags changed"

		} else {
			logMessage = "stack not ready yet, requeuing..."
			eventMessage = fmt.Sprintf("cnvrgapp %s not ready yet, done: %d%%", req.NamespacedName, stackReadiness.percentageReady)

		}
		log.Info(logMessage)
		r.recorder.Event(cnvrgApp, "Normal", "Creating", eventMessage)
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

	if app.Spec.Networking.Proxy.Enabled {
		if err := networking.NewHttpProxyState(app, r.Client, r.Scheme, r.Log).Apply(); err != nil {
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
	app.Status.LastFeatureFlagsHash = status.LastFeatureFlagsHash
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

	appOwnsPredicate := r.appOwnsPredicateFuncs()

	r.recorder = mgr.GetEventRecorderFor("cnvrgapp")
	a := &v1apps.Deployment{}
	a.GroupVersionKind()
	cnvrgAppController := ctrl.
		NewControllerManagedBy(mgr).
		Owns(&v1apps.Deployment{}, builder.WithPredicates(appOwnsPredicate)).
		Owns(&v1apps.StatefulSet{}, builder.WithPredicates(appOwnsPredicate)).
		For(&mlopsv1.CnvrgApp{}, builder.WithPredicates(appPredicate))

	r.Log.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))

	return cnvrgAppController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}

func (r *CnvrgAppReconciler) appOwnsPredicateFuncs() predicate.Funcs {
	return predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			if controllers.ContainsString(labelsMapToList(e.ObjectNew.GetLabels()), systemStatusHealthCheckLabelName) {
				return true
			}

			r.Log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			r.Log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return true
		},
	}
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

func (r *CnvrgAppReconciler) RollDeployment(name types.NamespacedName) error {
	ctx := context.Background()
	deployment := &v1apps.Deployment{}

	if err := r.Get(ctx, name, deployment); err != nil {
		return fmt.Errorf("failed to get deployment for rollout %s/%s : %v", name.Namespace, name.Name, err)
	}

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Annotations[RolloutAnnotation] = time.Now().Format(time.RFC3339)

	if err := r.Client.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to get update deployment for rollout %s/%s : %v", name.Namespace, name.Name, err)
	}

	return nil
}
