package controllers

import (
	"context"
	"github.com/cnvrg-operator/pkg/cnvrginfra/istio"
	"github.com/cnvrg-operator/pkg/cnvrginfra/registry"
	"github.com/cnvrg-operator/pkg/cnvrginfra/storage"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mlopsv1 "github.com/cnvrg-operator/api/v1"
)

// CnvrgInfraReconciler reconciles a CnvrgInfra object
type CnvrgInfraReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var cnvrgInfraLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras/status,verbs=get;update;patch

func (r *CnvrgInfraReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	cnvrgInfraLog = r.Log.WithValues("name", req.NamespacedName)
	cnvrgInfraLog.Info("starting cnvrginfra reconciliation")

	desiredSpec, err := r.defineDesiredSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if desiredSpec == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	// Setup finalizer
	if desiredSpec.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			desiredSpec.ObjectMeta.Finalizers = append(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), desiredSpec); err != nil {
				cnvrgInfraLog.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			r.updateStatusMessage(mlopsv1.STATUS_REMOVING, "removing cnvrg spec", desiredSpec)
			if err := r.cleanup(desiredSpec); err != nil {
				return ctrl.Result{}, err
			}
			desiredSpec.ObjectMeta.Finalizers = removeString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), desiredSpec); err != nil {
				cnvrgInfraLog.Info("error in removing finalizer, checking if cnvrgApp object still exists")
				// if update was failed, make sure that cnvrgApp still exists
				spec, e := r.getCnvrgInfraSpec(req.NamespacedName)
				if spec == nil && e == nil {
					return ctrl.Result{}, nil // probably spec was deleted, stop reconcile
				}
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// set reconciling status
	r.updateStatusMessage(mlopsv1.STATUS_RECONCILING, "reconciling", desiredSpec)

	// infra base config
	if err := desired.Apply(registry.State(desiredSpec), desiredSpec, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec)
		return ctrl.Result{}, err
	}

	// Istio
	if err := desired.Apply(istio.State(desiredSpec), desiredSpec, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec)
		return ctrl.Result{}, err
	}

	// Istio
	if err := desired.Apply(storage.State(desiredSpec), desiredSpec, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec)
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.STATUS_HEALTHY, "successfully reconciled", desiredSpec)
	return ctrl.Result{}, nil
}

func (r *CnvrgInfraReconciler) defineDesiredSpec(name types.NamespacedName) (*mlopsv1.CnvrgInfra, error) {
	cnvrgApp, err := r.getCnvrgInfraSpec(name)
	if err != nil {
		return nil, err
	}
	// probably cnvrgapp was removed
	if cnvrgApp == nil {
		return nil, nil
	}
	desiredSpec := mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}
	if err := mergo.Merge(&desiredSpec, cnvrgApp, mergo.WithOverride); err != nil {
		cnvrgInfraLog.Error(err, "can't merge")
		return nil, err
	}
	cnvrgInfraLog = r.Log.WithValues("name", name, "ns", desiredSpec.Namespace)
	return &desiredSpec, nil
}

func (r *CnvrgInfraReconciler) getCnvrgInfraSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgInfra, error) {
	ctx := context.Background()
	var cnvrgInfra mlopsv1.CnvrgInfra
	if err := r.Get(ctx, namespacedName, &cnvrgInfra); err != nil {
		if errors.IsNotFound(err) {
			cnvrgInfraLog.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		cnvrgInfraLog.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgInfra, nil
}

func (r *CnvrgInfraReconciler) cleanup(desiredSpec *mlopsv1.CnvrgInfra) error {
	cnvrgInfraLog.Info("running finalizer cleanup")

	// remove istio
	if err := r.cleanupIstio(desiredSpec); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) cleanupIstio(desiredSpec *mlopsv1.CnvrgInfra) error {
	cnvrgInfraLog.Info("running istio cleanup")
	ctx := context.Background()
	istioManifests := istio.State(desiredSpec)
	for _, m := range istioManifests {
		// Make sure IstioOperator was deployed
		if m.GVR == desired.Kinds[desired.IstioGVR] {
			if err := m.GenerateDeployable(desiredSpec); err != nil {
				cnvrgInfraLog.Error(err, "can't make manifest deployable")
				return err
			}
			if err := r.Delete(ctx, m.Obj); err != nil {
				if errors.IsNotFound(err) {
					cnvrgInfraLog.Info("istio instance not found - probably removed previously")
					return nil
				}
				return err
			}
			istioExists := true
			cnvrgInfraLog.Info("wait for istio instance removal")
			for istioExists {
				err := r.Get(ctx, types.NamespacedName{Name: m.Obj.GetName(), Namespace: m.Obj.GetNamespace()}, m.Obj)
				if err != nil && errors.IsNotFound(err) {
					cnvrgInfraLog.Info("istio instance was successfully removed")
					istioExists = false
				}
				if istioExists {
					cnvrgInfraLog.Info("istio instance still present, will sleep of 1 sec, and check again...")
				}
			}
		}
	}
	return nil
}

func (r *CnvrgInfraReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgInfra *mlopsv1.CnvrgInfra) {
	if cnvrgInfra.Status.Status == mlopsv1.STATUS_REMOVING {
		cnvrgInfraLog.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	ctx := context.Background()
	cnvrgInfra.Status.Status = status
	cnvrgInfra.Status.Message = message
	if err := r.Status().Update(ctx, cnvrgInfra); err != nil {
		cnvrgInfraLog.Error(err, "can't update status")
	}
	// This check is to make sure that the status is indeed updated
	// short reconciliations loop might cause status to be applied but not yet saved into BD
	// and leads to error: "the object has been modified; please apply your changes to the latest version and try again"
	// to avoid this error, fetch the object and compare the status
	statusCheckAttempts := 3
	for {
		cnvrgInfra, err := r.getCnvrgInfraSpec(types.NamespacedName{Namespace: cnvrgInfra.Namespace, Name: cnvrgInfra.Name})
		if err != nil {
			cnvrgInfraLog.Error(err, "can't validate status update")
		}
		cnvrgInfraLog.V(1).Info("expected status", "status", status, "message", message)
		cnvrgInfraLog.V(1).Info("current status", "status", cnvrgInfra.Status.Status, "message", cnvrgInfra.Status.Message)
		if cnvrgInfra.Status.Status == status && cnvrgInfra.Status.Message == message {
			break
		}
		if statusCheckAttempts == 0 {
			cnvrgInfraLog.Info("can't verify status update, status checks attempts exceeded")
			break
		}
		statusCheckAttempts--
		cnvrgInfraLog.V(1).Info("validating status update", "attempts", statusCheckAttempts)
		time.Sleep(1 * time.Second)
	}
}

func (r *CnvrgInfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cnvrgInfraLog = r.Log.WithValues("initializing", "crds")

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Warn("deploy-depended-crds is to false, I hope CRDs was deployed ahead, if not I will fail...")
	}

	if viper.GetBool("own-istio-resources") {
		err := desired.Apply(istio.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
		if err != nil {
			cnvrgInfraLog.Error(err, "can't apply networking CRDs")
			os.Exit(1)
		}
	}

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(e.ObjectOld) {
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgInfra)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgInfra)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrginfra spec wasn't changed, assuming status update, won't reconcile
				cnvrgInfraLog.V(1).Info("update received", "shouldReconcileOnSpecChange", shouldReconcileOnSpecChange)

				return !shouldReconcileOnSpecChange

			}
			return true
		},
	}

	cnvrgInfraController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgInfra{}).
		WithEventFilter(p)

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
		cnvrgInfraController.Owns(u)
	}

	return cnvrgInfraController.
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
