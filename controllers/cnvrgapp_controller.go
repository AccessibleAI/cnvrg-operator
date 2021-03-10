package controllers

import (
	"context"
	"fmt"
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
	"time"
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

	r.Log.Info("starting reconciliation")
	cnvrgApp, err := r.getCnvrgSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	// probably cnvrgapp was removed
	if cnvrgApp == nil {
		return ctrl.Result{}, nil
	}

	desiredSpec, err := r.defineDesiredSpec(&cnvrgApp.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Setup finalizer
	if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			cnvrgApp.ObjectMeta.Finalizers = append(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), cnvrgApp); err != nil {
				r.Log.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	} else {
		if containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			if err := r.cleanup(desiredSpec); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// set reconciling status
	r.updateStatusMessage(mlopsv1.STATUS_RECONCILING, "reconciling", cnvrgApp, req.NamespacedName)

	// PostgreSQL
	if err := r.apply(pg.State(desiredSpec), desiredSpec, cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), cnvrgApp, req.NamespacedName)
		return ctrl.Result{}, err
	}

	// Networking
	if err := r.apply(networking.State(desiredSpec), desiredSpec, cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), cnvrgApp, req.NamespacedName)
		return ctrl.Result{}, err
	}

	// ControlPlan
	if err := r.apply(controlplan.State(desiredSpec), desiredSpec, cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), cnvrgApp, req.NamespacedName)
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.STATUS_HEALTHY, "successfully reconciled", cnvrgApp, req.NamespacedName)
	return ctrl.Result{}, nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgApp *mlopsv1.CnvrgApp, name types.NamespacedName) {
	ctx := context.Background()
	cnvrgApp.Status.Status = status
	cnvrgApp.Status.Message = message
	if err := r.Status().Update(ctx, cnvrgApp); err != nil {
		r.Log.Error(err, "can't update status")
	}
	// This check is to make sure that the status is indeed updated
	// short reconciliations loop might cause status to be applied but not yet saved into BD
	// and leads to error: "the object has been modified; please apply your changes to the latest version and try again"
	// to avoid this error, fetch the object and compare the status
	statusCheckAttempts := 10
	for {
		cnvrgApp, err := r.getCnvrgSpec(name)
		if err != nil {
			r.Log.Error(err, "can't validate status update")
		}
		zap.S().Debugf("current   status   [%v] [%v]", status, message)
		zap.S().Debugf("expeceted status   [%v] [%v]", cnvrgApp.Status.Status, cnvrgApp.Status.Message)
		if cnvrgApp.Status.Status == status && cnvrgApp.Status.Message == message {
			break
		}
		if statusCheckAttempts == 0 {
			r.Log.Error(fmt.Errorf("status update failed"), "can't update status")
		}
		statusCheckAttempts--
		zap.S().Debugf("validating status update, left attempt: %v", statusCheckAttempts)
		time.Sleep(1 * time.Second)
	}

}

func (r *CnvrgAppReconciler) defineDesiredSpec(cnvrgAppSpec *mlopsv1.CnvrgAppSpec) (*mlopsv1.CnvrgAppSpec, error) {

	desiredSpec := mlopsv1.DefaultSpec
	if err := mergo.Merge(&desiredSpec, cnvrgAppSpec, mergo.WithOverride); err != nil {
		r.Log.Error(err, "can't merge desiredSpec")
		return nil, err
	}
	return &desiredSpec, nil
}

func (r *CnvrgAppReconciler) getCnvrgSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &cnvrgApp); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		r.Log.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgApp, nil
}

func (r *CnvrgAppReconciler) apply(desiredManifests []*desired.State, desiredSpec *mlopsv1.CnvrgAppSpec, cnvrgApp *mlopsv1.CnvrgApp) error {
	ctx := context.Background()
	for _, manifest := range desiredManifests {
		if err := manifest.GenerateDeployable(desiredSpec); err != nil {
			r.Log.Error(err, "error generating deployable", "name", manifest.Name)
			return err
		}
		if manifest.Own {
			if err := ctrl.SetControllerReference(cnvrgApp, manifest.Obj, r.Scheme); err != nil {
				r.Log.Error(err, "error setting controller reference", "name", manifest.Name)
				return err
			}
		}
		if viper.GetBool("dry-run") {
			r.Log.Info("dry run enabled, skipping applying...")
			continue
		}
		err := r.Get(ctx, types.NamespacedName{Name: manifest.Name, Namespace: desiredSpec.CnvrgNs}, manifest.Obj)
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

func (r *CnvrgAppReconciler) cleanup(desiredSpec *mlopsv1.CnvrgAppSpec) error {
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
				err := r.Get(ctx, types.NamespacedName{Name: m.Name, Namespace: desiredSpec.CnvrgNs}, m.Obj)
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
		if err := r.apply(networking.Crds(), &mlopsv1.DefaultSpec, &mlopsv1.CnvrgApp{Spec: mlopsv1.DefaultSpec}); err != nil {
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
