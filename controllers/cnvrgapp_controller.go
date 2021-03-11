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
	"time"
)

const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"

type CnvrgAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var log logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	log = r.Log.WithValues("name", req.NamespacedName)
	log.Info("starting reconciliation")

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
				log.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			if err := r.cleanup(desiredSpec); err != nil {
				return ctrl.Result{}, err
			}
			desiredSpec.ObjectMeta.Finalizers = removeString(desiredSpec.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), desiredSpec); err != nil {
				log.Info("error in removing finalizer, checking if cnvrgApp object still exists")
				// if update was failed, make sure that cnvrgApp still exists
				spec, e := r.getCnvrgSpec(req.NamespacedName)
				if spec == nil && e == nil {
					return ctrl.Result{}, nil // probably spec was deleted, stop reconcile
				}
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// set reconciling status
	r.updateStatusMessage(mlopsv1.STATUS_RECONCILING, "reconciling", desiredSpec, req.NamespacedName)

	// PostgreSQL
	if err := r.apply(pg.State(desiredSpec), desiredSpec); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec, req.NamespacedName)
		return ctrl.Result{}, err
	}

	// Networking
	if err := r.apply(networking.State(desiredSpec), desiredSpec); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec, req.NamespacedName)
		return ctrl.Result{}, err
	}

	// ControlPlan
	if err := r.apply(controlplan.State(desiredSpec), desiredSpec); err != nil {
		r.updateStatusMessage(mlopsv1.STATUS_ERROR, err.Error(), desiredSpec, req.NamespacedName)
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.STATUS_HEALTHY, "successfully reconciled", desiredSpec, req.NamespacedName)
	return ctrl.Result{}, nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgApp *mlopsv1.CnvrgApp, name types.NamespacedName) {
	ctx := context.Background()
	cnvrgApp.Status.Status = status
	cnvrgApp.Status.Message = message
	if err := r.Status().Update(ctx, cnvrgApp); err != nil {
		log.Error(err, "can't update status")
	}
	// This check is to make sure that the status is indeed updated
	// short reconciliations loop might cause status to be applied but not yet saved into BD
	// and leads to error: "the object has been modified; please apply your changes to the latest version and try again"
	// to avoid this error, fetch the object and compare the status
	statusCheckAttempts := 10
	for {
		cnvrgApp, err := r.getCnvrgSpec(name)
		if err != nil {
			log.Error(err, "can't validate status update")
		}
		log.V(1).Info("expected status", "status", status, "message", message)
		log.V(1).Info("current status", "status", cnvrgApp.Status.Status, "message", cnvrgApp.Status.Message)
		if cnvrgApp.Status.Status == status && cnvrgApp.Status.Message == message {
			break
		}
		if statusCheckAttempts == 0 {
			log.Info("can't verify status update, status checks attempts exceeded")
			break
		}
		statusCheckAttempts--
		log.V(1).Info("validating status update", "attempts", statusCheckAttempts)
		time.Sleep(1 * time.Second)
	}
}

func (r *CnvrgAppReconciler) defineDesiredSpec(name types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	cnvrgApp, err := r.getCnvrgSpec(name)
	if err != nil {
		return nil, err
	}
	// probably cnvrgapp was removed
	if cnvrgApp == nil {
		return nil, nil
	}
	desiredSpec := mlopsv1.CnvrgApp{Spec: mlopsv1.DefaultSpec}
	if err := mergo.Merge(&desiredSpec, cnvrgApp, mergo.WithOverride); err != nil {
		log.Error(err, "can't merge")
		return nil, err
	}
	log = r.Log.WithValues("name", name, "ns", desiredSpec.Spec.CnvrgNs)
	return &desiredSpec, nil
}

func (r *CnvrgAppReconciler) getCnvrgSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &cnvrgApp); err != nil {
		if errors.IsNotFound(err) {
			log.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		log.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgApp, nil
}

func (r *CnvrgAppReconciler) apply(desiredManifests []*desired.State, desiredSpec *mlopsv1.CnvrgApp) error {
	ctx := context.Background()
	for _, manifest := range desiredManifests {
		if err := manifest.GenerateDeployable(desiredSpec); err != nil {
			log.Error(err, "error generating deployable", "name", manifest.Name)
			return err
		}
		if manifest.Own {
			if err := ctrl.SetControllerReference(desiredSpec, manifest.Obj, r.Scheme); err != nil {
				log.Error(err, "error setting controller reference", "name", manifest.Name)
				return err
			}
		}
		if viper.GetBool("dry-run") {
			log.Info("dry run enabled, skipping applying...")
			continue
		}
		fetchInto := &unstructured.Unstructured{}
		fetchInto.SetGroupVersionKind(manifest.GVR)
		err := r.Get(ctx, types.NamespacedName{Name: manifest.Name, Namespace: desiredSpec.Spec.CnvrgNs}, fetchInto)
		if err != nil && errors.IsNotFound(err) {
			log.Info("creating", "name", manifest.Name, "kind", manifest.GVR.Kind)
			if err := r.Create(ctx, manifest.Obj); err != nil {
				log.Error(err, "error creating object", "name", manifest.Name)
				return err
			}
		} else {
			manifest.Obj.SetResourceVersion(fetchInto.GetResourceVersion())
			err := r.Update(ctx, manifest.Obj)
			if err != nil {
				log.Info("error updating object", "manifest", manifest.TemplatePath)
				return err
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanup(desiredSpec *mlopsv1.CnvrgApp) error {
	log.Info("running finalizer cleanup")
	ctx := context.Background()
	// remove istio
	istioManifests := networking.State(desiredSpec)
	for _, m := range istioManifests {
		// Make sure IstioOperator was deployed
		if m.GVR == desired.Kinds[desired.IstioGVR] {
			if err := m.GenerateDeployable(desiredSpec); err != nil {
				log.Error(err, "can't make manifest deployable")
				return err
			}
			if err := r.Delete(ctx, m.Obj); err != nil {
				if errors.IsNotFound(err) {
					log.Info("istio instance not found - probably removed previously")
					return nil
				}
				return err
			}

			istioExists := true
			log.Info("wait for istio instance removal")
			for istioExists {
				err := r.Get(ctx, types.NamespacedName{Name: m.Name, Namespace: desiredSpec.Spec.CnvrgNs}, m.Obj)
				if err != nil && errors.IsNotFound(err) {
					log.Info("istio instance was successfully removed")
					istioExists = false
				}
				if istioExists {
					log.Info("istio instance still present, will sleep of 1 sec, and check again...")
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
			log.Error(err, "can't apply networking CRDs")
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
		For(&mlopsv1.CnvrgApp{}).
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
