package controllers

import (
	"context"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/controlplan"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/cnvrg-operator/pkg/networking"
	"github.com/cnvrg-operator/pkg/pg"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"

type CnvrgAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	r.Log.Info("starting reconciliation")
	desiredSpec, err := r.defineDesiredSpec(req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if desiredSpec == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	// set reconciling status
	r.updateStatusMessage(mlopsv1.STATUS_RECONCILING, "reconciling", desiredSpec)

	// Setup finalizer
	if desiredSpec.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			desiredSpec.ObjectMeta.Finalizers = append(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(ctx, desiredSpec); err != nil {
				r.Log.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			if err := r.cleanup(desiredSpec); err != nil {
				return ctrl.Result{}, err
			}
			desiredSpec.ObjectMeta.Finalizers = removeString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer)
			cnvrgApp, err := r.getCnvrgSpec(req)
			if err != nil {
				return ctrl.Result{}, err
			}
			if cnvrgApp == nil {
				return ctrl.Result{}, nil
			}
			if err := r.Update(ctx, desiredSpec); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// PostgreSQL
	if err := r.apply(pg.State(desiredSpec), desiredSpec); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec)
		return ctrl.Result{}, err
	}

	// Networking
	if err := r.apply(networking.State(desiredSpec), desiredSpec); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec)
		return ctrl.Result{}, err
	}

	// ControlPlan
	if err := r.apply(controlplan.State(desiredSpec), desiredSpec); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec)
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.STATUS_HEALTHY, "successfully reconciled", desiredSpec)
	return ctrl.Result{}, nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, desiredSpec *mlopsv1.CnvrgApp) {
	ctx := context.Background()
	desiredSpec.Status.Status = status
	desiredSpec.Status.Message = message
	if err := r.Status().Update(ctx, desiredSpec); err != nil {
		r.Log.Error(err, "can't update status")
	}
}

func (r *CnvrgAppReconciler) defineDesiredSpec(req ctrl.Request) (*mlopsv1.CnvrgApp, error) {
	cnvrgApp, err := r.getCnvrgSpec(req)
	if err != nil {
		return nil, err
	}
	// probably cnvrgapp was removed
	if cnvrgApp == nil {
		return nil, nil
	}
	desiredSpec := mlopsv1.CnvrgApp{Spec: mlopsv1.DefaultSpec}
	if err := mergo.Merge(&desiredSpec, cnvrgApp, mergo.WithOverride); err != nil {
		r.Log.Error(err, "can't merge")
		return nil, err
	}
	return &desiredSpec, nil
}

func (r *CnvrgAppReconciler) getCnvrgSpec(req ctrl.Request) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, req.NamespacedName, &cnvrgApp); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		r.Log.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgApp, nil
}

func (r *CnvrgAppReconciler) apply(desiredManifests []*desired.State, desiredSpec *mlopsv1.CnvrgApp) error {
	ctx := context.Background()
	for _, manifest := range desiredManifests {
		if err := manifest.GenerateDeployable(desiredSpec); err != nil {
			r.Log.Error(err, "error generating deployable", "name", manifest.Name)
			return err
		}
		if manifest.Own {
			if err := ctrl.SetControllerReference(desiredSpec, manifest.Obj, r.Scheme); err != nil {
				r.Log.Error(err, "error setting controller reference", "name", manifest.Name)
				return err
			}
		}
		if viper.GetBool("dry-run") {
			r.Log.Info("dry run enabled, skipping applying...")
			continue
		}
		err := r.Get(ctx, types.NamespacedName{Name: manifest.Name, Namespace: desiredSpec.Spec.CnvrgNs}, manifest.Obj)
		if err != nil && errors.IsNotFound(err) {
			r.Log.Info("creating", "name", manifest.Name, "kind", manifest.GVR.Kind)
			if err := r.Create(ctx, manifest.Obj); err != nil {
				r.Log.Error(err, "error creating object", "name", manifest.Name)
				return err
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanup(desiredSpec *mlopsv1.CnvrgApp) error {
	ctx := context.Background()
	// remove istio
	istioManifests := networking.State(desiredSpec)
	for _, m := range istioManifests {
		// Make sure IstioOperator was deployed
		if m.GVR == desired.Kinds[desired.IstioGVR] {
			if err := m.GenerateDeployable(desiredSpec); err != nil {
				r.Log.Error(err, "can't make manifest deployable")
				return err
			}
			if err := r.Delete(ctx, m.Obj); err != nil {
				if errors.IsNotFound(err) {
					r.Log.Info("istio instance not found - probably removed previously")
					return nil
				}
				return err
			}
			r.Log.Info("has to remove istio first")
			istioExists := true
			for istioExists {
				err := r.Get(ctx, types.NamespacedName{Name: m.Name, Namespace: desiredSpec.Spec.CnvrgNs}, m.Obj)
				if err != nil && errors.IsNotFound(err) {
					r.Log.Info("istio instance was successfully removed")
					istioExists = false
				}
				if istioExists {
					r.Log.Info("istio instance still present, will sleep of 1 sec, and check again...")
				}
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Warn("deploy-depended-crds is to false, I hope CRDs was deployed ahead, if not I will fail...")
	}
	if viper.GetBool("own-istio-resources") {
		if err := r.apply(networking.Crds(), &mlopsv1.CnvrgApp{Spec: mlopsv1.DefaultSpec}); err != nil {
			r.Log.Error(err, "can't apply networking CRDs")
			os.Exit(1)
		}
	}

	cnvrgAppController := ctrl.NewControllerManagedBy(mgr).For(&mlopsv1.CnvrgApp{})

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
		cnvrgAppController.Owns(u)
	}

	pred := predicate.GenerationChangedPredicate{}
	return cnvrgAppController.
		WithEventFilter(pred).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
