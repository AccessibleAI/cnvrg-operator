package infra

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/infra/gpu"
	"github.com/AccessibleAI/cnvrg-operator/pkg/infra/istio"
	"github.com/AccessibleAI/cnvrg-operator/pkg/registry"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/d4l3k/messagediff.v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"os"
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

const CnvrginfraFinalizer = "cnvrginfra.mlops.cnvrg.io/finalizer"

type CnvrgInfraReconciler struct {
	client.Client
	recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

var log logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras/status,verbs=get;update;patch

func (r *CnvrgInfraReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log = r.Log.WithValues("name", req.NamespacedName)
	log.Info("starting cnvrginfra reconciliation")

	// sync specs between actual and defaults
	equal, err := r.syncCnvrgInfraSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		requeueAfter, _ := time.ParseDuration("3s")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	// specs are synced, proceed reconcile
	cnvrgInfra, err := r.getCnvrgInfraSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgInfra == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	// setup finalizer
	if cnvrgInfra.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer) {
			cnvrgInfra.ObjectMeta.Finalizers = append(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			if err := r.Update(ctx, cnvrgInfra); err != nil {
				log.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer) {
			r.updateStatusMessage(mlopsv1.StatusRemoving, "removing cnvrg spec", cnvrgInfra)
			if err := r.cleanup(cnvrgInfra); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgInfra.ObjectMeta.Finalizers = removeString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				if err := r.Update(ctx, cnvrgInfra); err != nil {
					cnvrgInfra, err := r.getCnvrgInfraSpec(req.NamespacedName)
					if err != nil {
						log.Error(err, "error getting cnvrginfra for finalizer cleanup")
						return err
					}
					cnvrgInfra.ObjectMeta.Finalizers = removeString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
					return r.Update(ctx, cnvrgInfra)
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

	r.updateStatusMessage(mlopsv1.StatusReconciling, "reconciling", cnvrgInfra)

	// apply manifests
	if err := r.applyManifests(cnvrgInfra); err != nil {
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.StatusHealthy, "successfully reconciled", cnvrgInfra)
	log.Info("successfully reconciled")
	return ctrl.Result{}, nil
}

func (r *CnvrgInfraReconciler) applyManifests(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	var reconcileResult error

	// registry
	log.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgInfra.Namespace,
		Data: map[string]interface{}{
			"Registry": cnvrgInfra.Spec.Registry,
		},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgInfra, r.Client, r.Scheme, log); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// istio
	log.Info("applying infra networking")
	if err := desired.Apply(istio.IstioState(), cnvrgInfra, r.Client, r.Scheme, log); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// nvidia device plugin
	if cnvrgInfra.Spec.Gpu.NvidiaDp.Enabled {
		log.Info("nvidia device plugin")
		nvidiaDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Namespace,
			Data: map[string]interface{}{
				"NvidiaDp": cnvrgInfra.Spec.Gpu.NvidiaDp,
				"Registry": cnvrgInfra.Spec.Registry,
				"ImageHub": cnvrgInfra.Spec.ImageHub,
			},
		}
		if err := desired.Apply(gpu.NvidiaDpState(nvidiaDpData), cnvrgInfra, r.Client, r.Scheme, log); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	// habana device plugin
	if cnvrgInfra.Spec.Gpu.HabanaDp.Enabled {
		log.Info("habana device plugin")
		habanaDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Namespace,
			Data: map[string]interface{}{
				"HabanaDp": cnvrgInfra.Spec.Gpu.HabanaDp,
				"Registry": cnvrgInfra.Spec.Registry,
				"ImageHub": cnvrgInfra.Spec.ImageHub,
			},
		}
		if err := desired.Apply(gpu.HabanaDpState(habanaDpData), cnvrgInfra, r.Client, r.Scheme, log); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	// metagpu device plugin
	if cnvrgInfra.Spec.Gpu.MetaGpuDp.Enabled {
		log.Info("metagpu device plugin")
		metagpuDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Namespace,
			Data: map[string]interface{}{
				"MetaGpuDp": cnvrgInfra.Spec.Gpu.MetaGpuDp,
				"ImageHub":  cnvrgInfra.Spec.ImageHub,
			},
		}
		// apply metagpu infra state
		if err := desired.Apply(gpu.MetagpudpState(metagpuDpData), cnvrgInfra, r.Client, r.Scheme, log); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}

	}

	return reconcileResult
}

func (r *CnvrgInfraReconciler) syncCnvrgInfraSpec(name types.NamespacedName) (bool, error) {

	log.Info("synchronizing cnvrgInfra spec")

	// Fetch current cnvrgInfra spec
	cnvrgInfra, err := r.getCnvrgInfraSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgInfra == nil {
		return true, nil // all (probably) good, cnvrginfra was removed
	}
	log = r.Log.WithValues("name", name, "ns", cnvrgInfra.Namespace)

	// Get default cnvrgInfra spec
	desiredSpec := mlopsv1.DefaultCnvrgInfraSpec()

	if err := CalculateAndApplyInfraDefaults(cnvrgInfra, &desiredSpec, r.Client); err != nil {
		log.Error(err, "can't calculate defaults")
		return false, err
	}

	// Merge current cnvrgInfra spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgInfra.Spec, mergo.WithOverride, mergo.WithTransformers(cnvrgSpecBoolTransformer{})); err != nil {
		log.Error(err, "can't merge")
		return false, err
	}

	log.V(1).Info("printing the diff between desiredSpec and actual")
	diff, _ := messagediff.PrettyDiff(desiredSpec, cnvrgInfra.Spec)
	log.V(1).Info(diff)

	// Compare desiredSpec and current cnvrgInfra spec,
	// if they are not equal, update the cnvrgInfra spec with desiredSpec,
	// and return true for triggering new reconciliation
	equal := reflect.DeepEqual(desiredSpec, cnvrgInfra.Spec)
	if !equal {
		log.Info("states are not equals, syncing and requeuing")
		cnvrgInfra.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgInfra); err != nil && errors.IsConflict(err) {
			log.Error(err, "conflict updating cnvrgInfra object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	log.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgInfraReconciler) getCnvrgInfraSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgInfra, error) {
	ctx := context.Background()
	var cnvrgInfra mlopsv1.CnvrgInfra
	if err := r.Get(ctx, namespacedName, &cnvrgInfra); err != nil {
		if errors.IsNotFound(err) {
			log.Info("unable to fetch CnvrgInfra, probably cr was deleted")
			return nil, nil
		}
		log.Error(err, "unable to fetch CnvrgInfra")
		return nil, err
	}
	return &cnvrgInfra, nil
}

func (r *CnvrgInfraReconciler) cleanup(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	log.Info("running finalizer cleanup")

	// cleanup pvc
	if err := r.cleanupPVCs(cnvrgInfra); err != nil {
		return err
	}

	log.Info("cleanup has been finished")
	return nil
}

func (r *CnvrgInfraReconciler) cleanupPVCs(infra *mlopsv1.CnvrgInfra) error {
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
		if pvc.Namespace == infra.Namespace {
			if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
				if pvc.ObjectMeta.Labels["app"] == "prometheus" {
					if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
						log.Info("prometheus pvc already deleted")
					} else if err != nil {
						log.Error(err, "error deleting prometheus pvc")
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *CnvrgInfraReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgInfra *mlopsv1.CnvrgInfra) {
	if cnvrgInfra.Status.Status == mlopsv1.StatusRemoving {
		log.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	ctx := context.Background()
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		name := types.NamespacedName{Namespace: "", Name: cnvrgInfra.Name}
		infra, err := r.getCnvrgInfraSpec(name)
		if err != nil {
			return err
		}
		infra.Status.Status = status
		infra.Status.Message = message
		err = r.Status().Update(ctx, infra)
		return err
	})
	if err != nil {
		log.Error(err, "can't update status")
	}
}

func (r *CnvrgInfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log = r.Log.WithValues("initializing", "crds")

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Info("deploy-depended-crds is false, I hope CRDs was deployed ahead and match expected versions, if not I will fail...")
	} else {

		if viper.GetBool("own-istio-resources") {
			err := desired.Apply(istio.IstioCrds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r.Client, r.Scheme, r.Log)
			if err != nil {
				log.Error(err, "can't apply istio CRDs")
				os.Exit(1)
			}
		}

		err := desired.Apply(controlplane.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r.Client, r.Scheme, r.Log)
		if err != nil {
			log.Error(err, "can't apply control plane crds")
			os.Exit(1)
		}
	}

	infraPredicate := predicate.Funcs{

		CreateFunc: func(createEvent event.CreateEvent) bool {
			msg := fmt.Sprintf("cnvrginfra: %s has been created", createEvent.Object.GetName())
			r.recorder.Event(createEvent.Object, "Normal", "Created", msg)
			return true
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			shouldReconcile := e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
			if shouldReconcile {
				msg := fmt.Sprintf("cnvrginfra: %s has been updated", e.ObjectNew.GetName())
				r.recorder.Event(e.ObjectNew, "Normal", "Updated", msg)
			}
			return shouldReconcile
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			msg := fmt.Sprintf("cnvrginfra: %s has been deleted", deleteEvent.Object.GetName())
			r.recorder.Event(deleteEvent.Object, "Normal", "SuccessfulDelete", msg)
			log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return !deleteEvent.DeleteStateUnknown
		},
	}

	infraOwnsPredicate := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {
			log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return true
		},
	}
	r.recorder = mgr.GetEventRecorderFor("cnvrginfra")
	cnvrgInfraController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgInfra{}, builder.WithPredicates(infraPredicate))

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
		cnvrgInfraController.Owns(u, builder.WithPredicates(infraOwnsPredicate))
	}

	log.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgInfraController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}
