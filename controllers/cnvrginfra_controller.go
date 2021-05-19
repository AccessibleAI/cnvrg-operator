package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/controlplane"
	"github.com/cnvrg-operator/pkg/dbs"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/cnvrg-operator/pkg/gpu"
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
	"time"
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
		requeueAfter, _ := time.ParseDuration("3s")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
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
			cnvrgInfra, err := r.getCnvrgInfraSpec(req.NamespacedName)
			if err != nil {
				return ctrl.Result{}, err
			}
			if cnvrgInfra == nil {
				return ctrl.Result{}, nil
			}
			cnvrgInfra.ObjectMeta.Finalizers = removeString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				cnvrgInfraLog.Info("error in removing finalizer, checking if cnvrgInfra object still exists")
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

func (r *CnvrgInfraReconciler) getCnvrgAppInstances(infra *mlopsv1.CnvrgInfra) ([]mlopsv1.AppInstance, error) {

	cmName := types.NamespacedName{Namespace: infra.Spec.InfraNamespace, Name: mlopsv1.InfraReconcilerCm}

	cnvrgAppCm := &v1.ConfigMap{}
	if err := r.Get(context.Background(), cmName, cnvrgAppCm); err != nil && errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var apps []mlopsv1.AppInstance
	for _, appJson := range cnvrgAppCm.Data {
		var app mlopsv1.AppInstance
		if err := json.Unmarshal([]byte(appJson), &app); err != nil {
			cnvrgInfraLog.Error(err, "error decoding AppInstance")
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *CnvrgInfraReconciler) applyManifests(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	var reconcileResult error

	// registry
	cnvrgInfraLog.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data: map[string]interface{}{
			"Registry":    cnvrgInfra.Spec.Registry,
			"Annotations": cnvrgInfra.Annotations,
			"Labels":      cnvrgInfra.Labels,
		},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// storage
	cnvrgInfraLog.Info("applying storage")
	if err := desired.Apply(storage.State(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// redis
	if *cnvrgInfra.Spec.Dbs.Redis.Enabled || *cnvrgInfra.Spec.SSO.Enabled {
		cnvrgInfraLog.Info("applying redis")
		if err := desired.CreateRedisCredsSecret(cnvrgInfra,
			cnvrgInfra.Spec.Dbs.Redis.CredsRef,
			cnvrgInfra.Spec.InfraNamespace,
			fmt.Sprintf("%s:%d", cnvrgInfra.Spec.Dbs.Redis.SvcName, cnvrgInfra.Spec.Dbs.Redis.Port),
			r,
			r.Scheme,
			cnvrgInfraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
		if err := desired.Apply(dbs.InfraDbsState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	// logging
	cnvrgInfraLog.Info("applying logging")
	cnvrgApps, err := r.getCnvrgAppInstances(cnvrgInfra)
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	fluentbitData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data: map[string]interface{}{
			"AppInstance": cnvrgApps,
			"Annotations": cnvrgInfra.Annotations,
			"Labels":      cnvrgInfra.Labels,
		},
	}
	if err := desired.Apply(logging.FluentbitConfigurationState(fluentbitData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	if err := desired.Apply(logging.InfraLoggingState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// istio
	cnvrgInfraLog.Info("applying istio")
	if err := desired.Apply(networking.IstioInstanceState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// monitoring
	cnvrgInfraLog.Info("applying monitoring")
	if err := r.monitoringState(cnvrgInfra); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// mpi infra
	cnvrgInfraLog.Info("applying mpi infra")
	if err := desired.Apply(controlplane.MpiInfraState(), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// nvidia device plugin
	if *cnvrgInfra.Spec.Gpu.NvidiaDp.Enabled {
		cnvrgInfraLog.Info("nvidia device plugin")
		nvidiaDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"NvidiaDp": cnvrgInfra.Spec.Gpu.NvidiaDp,
				"Registry": cnvrgInfra.Spec.Registry,
			},
		}
		if err := desired.Apply(gpu.NvidiaDpState(nvidiaDpData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	return reconcileResult
}

func (r *CnvrgInfraReconciler) monitoringState(infra *mlopsv1.CnvrgInfra) error {

	if *infra.Spec.Monitoring.Prometheus.Enabled {
		err := desired.CreatePromCredsSecret(infra,
			infra.Spec.Monitoring.Prometheus.CredsRef,
			infra.Spec.InfraNamespace,
			fmt.Sprintf("http://%s.%s.svc:%d", infra.Spec.Monitoring.Prometheus.SvcName, infra.Spec.InfraNamespace, infra.Spec.Monitoring.Prometheus.Port),
			r,
			r.Scheme,
			cnvrgInfraLog)
		if err != nil {

			return err
		}
	}

	if *infra.Spec.Monitoring.Grafana.Enabled {
		// grafana dashboards
		cnvrgInfraLog.Info("applying grafana dashboards")
		if err := r.createGrafanaDashboards(infra); err != nil {
			return err
		}

		// grafana datasource
		cnvrgInfraLog.Info("applying grafana datasource")
		url, basicAuthUser, basicAuthPass, err := desired.GetPromCredsSecret(infra.Spec.Monitoring.Prometheus.CredsRef, infra.Spec.InfraNamespace, r, cnvrgInfraLog)
		if err != nil {
			return err
		}
		grafanaDatasourceData := desired.TemplateData{
			Namespace: infra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"Url":         url,
				"User":        basicAuthUser,
				"Pass":        basicAuthPass,
				"Annotations": infra.Annotations,
				"Labels":      infra.Labels,
			},
		}
		if err := desired.Apply(monitoring.GrafanaDSState(grafanaDatasourceData), infra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
			return err
		}
	}
	// monitoring
	cnvrgInfraLog.Info("applying monitoring")
	if err := desired.Apply(monitoring.InfraMonitoringState(infra), infra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) createGrafanaDashboards(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	if !*cnvrgInfra.Spec.Monitoring.Grafana.Enabled {
		cnvrgInfraLog.Info("grafana disabled, skipping grafana deployment")
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
			cnvrgAppLog.Error(err, "conflict updating cnvrgInfra object, requeue for reconciliations...")
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

	// todo: completely remove istio cleanup when you are emotional ready
	// remove istio
	//if err := r.cleanupIstio(cnvrgInfra); err != nil {
	//	return err
	//}

	// cleanup pvc
	if err := r.cleanupPVCs(cnvrgInfra); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) cleanupPVCs(infra *mlopsv1.CnvrgInfra) error {
	if !viper.GetBool("cleanup-pvc") {
		cnvrgAppLog.Info("cleanup-pvc is false, skipping pvc deletion!")
		return nil
	}
	cnvrgAppLog.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		cnvrgAppLog.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if pvc.Namespace == infra.Spec.InfraNamespace {
			if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
				if pvc.ObjectMeta.Labels["app"] == "prometheus" {
					if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
						cnvrgInfraLog.Info("prometheus pvc already deleted")
					} else if err != nil {
						cnvrgInfraLog.Error(err, "error deleting prometheus pvc")
						return err
					}
				}
			}
		}
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
			if m.TemplateData == nil {
				m.TemplateData = cnvrgInfra
			}
			if err := m.GenerateDeployable(); err != nil {
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
				time.Sleep(1 * time.Second)
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
		cnvrgInfraLog.Error(err, "can't update status")
	}
}

func (r *CnvrgInfraReconciler) createInfraReconcilerTriggerCm(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: mlopsv1.InfraReconcilerCm, Namespace: cnvrgInfra.Spec.InfraNamespace},
	}
	if err := ctrl.SetControllerReference(cnvrgInfra, cm, r.Scheme); err != nil {
		cnvrgInfraLog.Error(err, "failed to set ControllerReference", "cm", mlopsv1.InfraReconcilerCm)
		return err
	}
	if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
		cnvrgInfraLog.Info("already exists", "cm", mlopsv1.InfraReconcilerCm)
	} else if err != nil {
		cnvrgInfraLog.Error(err, "error creating", "cm", mlopsv1.InfraReconcilerCm)
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cnvrgInfraLog = r.Log.WithValues("initializing", "crds")

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Info("deploy-depended-crds is false, I hope CRDs was deployed ahead and match expected versions, if not I will fail...")
	} else {

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
	}

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			// run reconcile only changing cnvrginfra/object marked for deletion
			if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(e.ObjectOld) {
				cnvrgInfraLog.V(1).Info("received UpdateEvent", "eventSourcesObjectName", e.MetaNew.GetName())
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgInfra)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgInfra)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrginfra spec wasn't changed, assuming status update, won't reconcile

				if !shouldReconcileOnSpecChange && viper.GetBool("verbose") {
					cnvrgInfraLog.V(1).Info("printing the diff between oldObject.Spec and newObject.Spec")
					diff, _ := messagediff.PrettyDiff(oldObject.Spec, newObject.Spec)
					cnvrgInfraLog.V(1).Info(diff)
				}

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
