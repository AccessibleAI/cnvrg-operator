package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/controlplane"
	"github.com/cnvrg-operator/pkg/dbs"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/cnvrg-operator/pkg/logging"
	"github.com/cnvrg-operator/pkg/monitoring"
	"github.com/cnvrg-operator/pkg/networking"
	"github.com/cnvrg-operator/pkg/registry"
	"github.com/cnvrg-operator/pkg/storage"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/d4l3k/messagediff.v1"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"os"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

const CnvrginfraFinalizer = "cnvrginfra.mlops.cnvrg.io/finalizer"

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

	equal, err := r.syncCnvrgInfraSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		return ctrl.Result{Requeue: true}, nil
	}

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
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				cnvrgInfraLog.Error(err, "failed to add finalizer")
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
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				cnvrgInfraLog.Info("error in removing finalizer, checking if cnvrgApp object still exists")
				// if update was failed, make sure that cnvrgInfra still exists
				spec, e := r.getCnvrgInfraSpec(req.NamespacedName)
				if spec == nil && e == nil {
					return ctrl.Result{}, nil // probably spec was deleted, stop reconcile
				}
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

	// infra reconciler trigger configmap
	if err := r.createInfraReconcilerTriggerCm(cnvrgInfra); err != nil {
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.StatusHealthy, "successfully reconciled", cnvrgInfra)
	cnvrgInfraLog.Info("successfully reconciled")
	return ctrl.Result{}, nil
}

func (r *CnvrgInfraReconciler) getCnvrgAppInstances(infra *mlopsv1.CnvrgInfra) ([]mlopsv1.CnvrgAppInstance, error) {
	var cnvrgAppInstances []mlopsv1.CnvrgAppInstance
	cmName := types.NamespacedName{Namespace: infra.Spec.InfraNamespace, Name: infra.Spec.InfraReconcilerCm}
	if cmName.Name == "" {
		cmName.Name = mlopsv1.DefaultCnvrgInfraSpec().InfraReconcilerCm
	}
	if cmName.Namespace == "" {
		cmName.Namespace = infra.Spec.InfraNamespace
	}
	cnvrgAppCm := &v1.ConfigMap{}
	if err := r.Get(context.Background(), cmName, cnvrgAppCm); err != nil && errors.IsNotFound(err) {
		return cnvrgAppInstances, nil
	} else if err != nil {
		return nil, err
	}
	for cnvrgAppNamespace, cnvrgAppName := range cnvrgAppCm.Data {
		cnvrgAppInstances = append(cnvrgAppInstances, mlopsv1.CnvrgAppInstance{
			Name:      cnvrgAppName,
			Namespace: cnvrgAppNamespace,
		})
	}

	return cnvrgAppInstances, nil
}

func (r *CnvrgInfraReconciler) applyManifests(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	var reconcileResult error

	// storage
	cnvrgInfraLog.Info("applying storage")
	if err := desired.Apply(storage.State(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// logging
	cnvrgInfraLog.Info("applying logging")
	if err := desired.Apply(logging.InfraLoggingState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// grafana dashboards
	cnvrgInfraLog.Info("applying grafana dashboards")
	if err := r.createGrafanaDashboards(cnvrgInfra); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// monitoring
	cnvrgInfraLog.Info("applying monitoring")
	if err := desired.Apply(monitoring.InfraMonitoringState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// registry
	cnvrgInfraLog.Info("applying registry")
	if err := desired.Apply(registry.State(), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// istio
	cnvrgInfraLog.Info("applying istio")
	if err := desired.Apply(networking.IstioInstanceState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// redis
	cnvrgInfraLog.Info("applying redis")
	if err := desired.Apply(dbs.InfraDbsState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// mpi infra
	cnvrgInfraLog.Info("applying mpi infra")
	if err := desired.Apply(controlplane.MpiInfraState(), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	return reconcileResult
}

func (r *CnvrgInfraReconciler) createGrafanaDashboards(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	if cnvrgInfra.Spec.Monitoring.Enabled != "true" {
		cnvrgInfraLog.Info("monitoring disabled, skipping grafana deployment")
		return nil
	}

	basePath := "/pkg/monitoring/tmpl/grafana/dashboards-data/"
	for _, dashboard := range desired.GrafanaInfraDashboards {
		f, err := pkger.Open(basePath + dashboard)
		if err != nil {
			cnvrgAppLog.Error(err, "error reading path", "path", dashboard)
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			cnvrgAppLog.Error(err, "error reading", "file", dashboard)
			return err
		}
		cm := &v1core.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name())),
				Namespace: cnvrgInfra.Spec.InfraNamespace,
			},
			Data: map[string]string{filepath.Base(f.Name()): string(b)},
		}
		if err := ctrl.SetControllerReference(cnvrgInfra, cm, r.Scheme); err != nil {
			cnvrgAppLog.Error(err, "error setting controller reference", "file", f.Name())
			return err
		}
		if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
			cnvrgAppLog.V(1).Info("grafana dashboard already exists", "file", dashboard)
			continue
		} else if err != nil {
			cnvrgAppLog.Error(err, "error reading", "file", dashboard)
			return err
		}
	}

	return nil

}

func (r *CnvrgInfraReconciler) syncCnvrgInfraSpec(name types.NamespacedName) (bool, error) {

	cnvrgInfraLog.Info("synchronizing cnvrgInfra spec")

	// Fetch current cnvrgInfra spec
	cnvrgInfra, err := r.getCnvrgInfraSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgInfra == nil {
		return false, nil // probably cnvrgapp was removed
	}
	cnvrgInfraLog = r.Log.WithValues("name", name, "ns", cnvrgInfra.Spec.InfraNamespace)

	// Get default cnvrgInfra spec
	desiredSpec := mlopsv1.DefaultCnvrgInfraSpec()
	cnvrgAppInstances, err := r.getCnvrgAppInstances(cnvrgInfra)
	if err != nil {
		return false, err
	}
	desiredSpec.CnvrgAppInstances = cnvrgAppInstances

	// Merge current cnvrgInfra spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgInfra.Spec, mergo.WithOverride); err != nil {
		cnvrgInfraLog.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {
		cnvrgInfraLog.V(1).Info("printing the diff between desiredSpec and actual")
		diff, _ := messagediff.PrettyDiff(desiredSpec, cnvrgInfra.Spec)
		cnvrgInfraLog.V(1).Info(diff)
	}

	// Compare desiredSpec and current cnvrgInfra spec,
	// if they are not equal, update the cnvrgInfra spec with desiredSpec,
	// and return true for triggering new reconciliation
	equal := reflect.DeepEqual(desiredSpec, cnvrgInfra.Spec)
	if !equal {
		cnvrgInfraLog.Info("states are not equals, syncing and requeuing")
		cnvrgInfra.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgInfra); err != nil && errors.IsConflict(err) {
			cnvrgAppLog.Info("conflict updating cnvrgInfra object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	// make sure cnvrgAppInstances are synced
	equal = reflect.DeepEqual(desiredSpec.CnvrgAppInstances, cnvrgAppInstances)
	if !equal {
		cnvrgInfraLog.Info("states are not equals (invalid cnvrgAppInstances), syncing and requeuing")
		// cnvrgApp instances must be calculated at runtime
		cnvrgInfra.Spec.CnvrgAppInstances = cnvrgAppInstances
		if err := r.Update(context.Background(), cnvrgInfra); err != nil && errors.IsConflict(err) {
			cnvrgAppLog.Info("conflict updating cnvrgInfra object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	cnvrgInfraLog.Info("states are equals, no need to sync")
	return equal, nil
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

func (r *CnvrgInfraReconciler) cleanup(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cnvrgInfraLog.Info("running finalizer cleanup")

	// remove istio
	if err := r.cleanupIstio(cnvrgInfra); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) cleanupIstio(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cnvrgInfraLog.Info("running istio cleanup")
	ctx := context.Background()
	istioManifests := networking.IstioInstanceState(cnvrgInfra)
	for _, m := range istioManifests {
		// Make sure IstioOperator was deployed
		if m.GVR == desired.Kinds[desired.IstioGVR] {
			if err := m.GenerateDeployable(cnvrgInfra); err != nil {
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
	if cnvrgInfra.Status.Status == mlopsv1.StatusRemoving {
		cnvrgInfraLog.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	ctx := context.Background()
	cnvrgInfra.Status.Status = status
	cnvrgInfra.Status.Message = message
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := r.Status().Update(ctx, cnvrgInfra)
		return err
	})
	if err != nil {
		cnvrgInfraLog.Error(err, "can't update status")
	}
	//// This check is to make sure that the status is indeed updated
	//// short reconciliations loop might cause status to be applied but not yet saved into BD
	//// and leads to error: "the object has been modified; please apply your changes to the latest version and try again"
	//// to avoid this error, fetch the object and compare the status
	//statusCheckAttempts := 3
	//for {
	//	cnvrgInfra, err := r.getCnvrgInfraSpec(types.NamespacedName{Namespace: cnvrgInfra.Spec.InfraNamespace, Name: cnvrgInfra.Name})
	//	if err != nil {
	//		cnvrgInfraLog.Error(err, "can't validate status update")
	//	}
	//	cnvrgInfraLog.V(1).Info("expected status", "status", status, "message", message)
	//	cnvrgInfraLog.V(1).Info("current status", "status", cnvrgInfra.Status.Status, "message", cnvrgInfra.Status.Message)
	//	if cnvrgInfra.Status.Status == status && cnvrgInfra.Status.Message == message {
	//		break
	//	}
	//	if statusCheckAttempts == 0 {
	//		cnvrgInfraLog.Info("can't verify status update, status checks attempts exceeded")
	//		break
	//	}
	//	statusCheckAttempts--
	//	cnvrgInfraLog.V(1).Info("validating status update", "attempts", statusCheckAttempts)
	//	time.Sleep(1 * time.Second)
	//}
}

func (r *CnvrgInfraReconciler) createInfraReconcilerTriggerCm(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: cnvrgInfra.Spec.InfraReconcilerCm, Namespace: cnvrgInfra.Spec.InfraNamespace},
	}
	if err := ctrl.SetControllerReference(cnvrgInfra, cm, r.Scheme); err != nil {
		cnvrgInfraLog.Error(err, "failed to set ControllerReference", "cm", cnvrgInfra.Spec.InfraReconcilerCm)
		return err
	}
	if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
		cnvrgInfraLog.Info("already exists", "cm", cnvrgInfra.Spec.InfraReconcilerCm)
	} else if err != nil {
		cnvrgInfraLog.Error(err, "error creating", "cm", cnvrgInfra.Spec.InfraReconcilerCm)
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cnvrgInfraLog = r.Log.WithValues("initializing", "crds")

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Warn("deploy-depended-crds is to false, I hope CRDs was deployed ahead, if not I will fail...")
	}

	if viper.GetBool("own-istio-resources") {
		err := desired.Apply(networking.IstioCrds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
		if err != nil {
			cnvrgInfraLog.Error(err, "can't apply istio CRDs")
			os.Exit(1)
		}
	}

	if viper.GetBool("own-prometheus-resources") {
		err := desired.Apply(monitoring.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
		if err != nil {
			cnvrgInfraLog.Error(err, "can't apply prometheus CRDs")
			os.Exit(1)
		}
	}

	err := desired.Apply(controlplane.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
	if err != nil {
		cnvrgInfraLog.Error(err, "can't apply MPI CRDs")
		os.Exit(1)
	}

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			// run reconcile only changing cnvrginfra/object marked for deletion
			if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(e.ObjectOld) {
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgInfra)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgInfra)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrginfra spec wasn't changed, assuming status update, won't reconcile
				cnvrgInfraLog.V(1).Info("cnvrginfra update received", "shouldReconcileOnSpecChange", shouldReconcileOnSpecChange)

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
	cnvrgInfraLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgInfraController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}
