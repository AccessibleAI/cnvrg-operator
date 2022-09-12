package thirdparty

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	controllers "github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/thirdparty/istio"
	"github.com/AccessibleAI/cnvrg-operator/pkg/thirdparty/metagpu"
	"github.com/AccessibleAI/cnvrg-operator/pkg/thirdparty/nvidia"
	ctpregistry "github.com/AccessibleAI/cnvrg-operator/pkg/thirdparty/registry"

	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"gopkg.in/d4l3k/messagediff.v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

type CnvrgThirdPartyReconciler struct {
	client.Client
	recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

var log logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgthirdparties,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgthirdparties/status,verbs=get;update;patch

func (r *CnvrgThirdPartyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log = r.Log.WithValues("name", req.NamespacedName)
	log.Info("starting cnvrgThirdParty reconciliation")

	// sync specs between actual and defaults
	equal, err := r.syncCnvrgThirdPartySpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		requeueAfter, _ := time.ParseDuration("3s")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	// specs are synced, proceed reconcile
	cnvrgTp, err := r.getCnvrgThirdPartySpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgTp == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	r.updateStatusMessage(mlopsv1.StatusReconciling, "reconciling", cnvrgTp)

	// apply manifests
	if err := r.applyManifests(cnvrgTp); err != nil {
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.StatusHealthy, "successfully reconciled", cnvrgTp)
	log.Info("successfully reconciled")
	return ctrl.Result{}, nil
}

func (r *CnvrgThirdPartyReconciler) applyManifests(ctp *mlopsv1.CnvrgThirdParty) error {

	var reconcileResult error

	// registry
	if err := ctpregistry.NewRegistryStateManager(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
		reconcileResult = err
	}

	// istio
	if ctp.Spec.Istio.Enabled {
		if err := istio.NewIstioCrdsState(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			reconcileResult = err
		}

		if err := istio.NewIstioInstanceState(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			reconcileResult = err
		}
	}

	// nvidia
	if ctp.Spec.Nvidia.DevicePlugin.Enabled || ctp.Spec.Nvidia.MetricsExporter.Enabled {
		if err := nvidia.NewNvidiaRbacState(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			reconcileResult = err
		}
	}

	if ctp.Spec.Nvidia.DevicePlugin.Enabled {
		if err := nvidia.NewNvidiaStateManager(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			reconcileResult = err
		}
	}

	if ctp.Spec.Nvidia.MetricsExporter.Enabled {
		if err := nvidia.NewMetricsExporterState(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			reconcileResult = err
		}
	}

	// metagpu
	if ctp.Spec.Metagpu.Enabled {
		if err := metagpu.NewDevicePluginState(ctp, r.Client, r.Scheme, r.Log).Apply(); err != nil {
			reconcileResult = err
		}
	}

	return reconcileResult
}

func (r *CnvrgThirdPartyReconciler) syncCnvrgThirdPartySpec(name types.NamespacedName) (bool, error) {

	log.Info("synchronizing cnvrgThirdParty spec")

	// Fetch current cnvrgThirdParty spec
	ctp, err := r.getCnvrgThirdPartySpec(name)
	if err != nil {
		return false, err
	}
	if ctp == nil {
		return true, nil // all (probably) good, ctp was removed
	}
	log = r.Log.WithValues("name", name, "ns", ctp.Namespace)

	// Get default ctp spec
	desiredSpec := mlopsv1.DefaultCnvrgThirdPartySpec()

	if err := applyCtpDefaults(ctp, &desiredSpec, r.Client); err != nil {
		r.Log.Error(err, "failed to apply ctp defaults ")
		return false, err
	}

	// Merge current ctp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, ctp.Spec, mergo.WithOverride, mergo.WithTransformers(controllers.CnvrgSpecBoolTransformer{})); err != nil {
		log.Error(err, "can't merge")
		return false, err
	}

	log.V(1).Info("printing the diff between desiredSpec and actual")
	diff, _ := messagediff.PrettyDiff(desiredSpec, ctp.Spec)
	log.V(1).Info(diff)

	// Compare desiredSpec and current ctp spec,
	// if they are not equal, update the ctp spec with desiredSpec,
	// and return true for triggering new reconciliation
	equal := reflect.DeepEqual(desiredSpec, ctp.Spec)
	if !equal {
		log.Info("states are not equals, syncing and requeuing")
		ctp.Spec = desiredSpec
		if err := r.Update(context.Background(), ctp); err != nil && errors.IsConflict(err) {
			log.Error(err, "conflict updating ctp object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	log.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgThirdPartyReconciler) getCnvrgThirdPartySpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgThirdParty, error) {
	ctx := context.Background()
	var cnvrgThirdParty mlopsv1.CnvrgThirdParty
	if err := r.Get(ctx, namespacedName, &cnvrgThirdParty); err != nil {
		if errors.IsNotFound(err) {
			log.Info("unable to fetch CnvrgThirdParty, probably cr was deleted")
			return nil, nil
		}
		log.Error(err, "unable to fetch CnvrgThirdParty")
		return nil, err
	}
	return &cnvrgThirdParty, nil
}

func (r *CnvrgThirdPartyReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgTp *mlopsv1.CnvrgThirdParty) {
	if cnvrgTp.Status.Status == mlopsv1.StatusRemoving {
		log.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		cnvrgTp.Status.Status = status
		cnvrgTp.Status.Message = message
		err := r.Status().Update(context.Background(), cnvrgTp)
		return err
	})
	if err != nil {
		log.Error(err, "can't update status")
	}
}

func (r *CnvrgThirdPartyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log = r.Log.WithValues("initializing", "crds")

	tpPredicate := predicate.Funcs{

		CreateFunc: func(createEvent event.CreateEvent) bool {
			msg := fmt.Sprintf("cnvrgtp: %s has been created", createEvent.Object.GetName())
			r.recorder.Event(createEvent.Object, "Normal", "Created", msg)
			return true
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			shouldReconcile := e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
			if shouldReconcile {
				msg := fmt.Sprintf("cnvrgtp: %s has been updated", e.ObjectNew.GetName())
				r.recorder.Event(e.ObjectNew, "Normal", "Updated", msg)
			}
			return shouldReconcile
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			msg := fmt.Sprintf("cnvrgtp: %s has been deleted", deleteEvent.Object.GetName())
			r.recorder.Event(deleteEvent.Object, "Normal", "SuccessfulDelete", msg)
			log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return !deleteEvent.DeleteStateUnknown
		},
	}

	cnvrgTpOwnsPredicate := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {
			log.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			log.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return true
		},
	}
	r.recorder = mgr.GetEventRecorderFor("cnvrgtp")
	cnvrgTpController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgThirdParty{}, builder.WithPredicates(tpPredicate))

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
		cnvrgTpController.Owns(u, builder.WithPredicates(cnvrgTpOwnsPredicate))
	}

	log.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgTpController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}
