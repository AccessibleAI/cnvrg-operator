package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/capsule"
	"github.com/AccessibleAI/cnvrg-operator/pkg/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/dbs"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/gpu"
	"github.com/AccessibleAI/cnvrg-operator/pkg/logging"
	"github.com/AccessibleAI/cnvrg-operator/pkg/monitoring"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/priorityclass"
	"github.com/AccessibleAI/cnvrg-operator/pkg/registry"
	"github.com/AccessibleAI/cnvrg-operator/pkg/reloader"
	"github.com/AccessibleAI/cnvrg-operator/pkg/storage"
	"github.com/Dimss/crypt/apr1_crypt"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/opentracing/opentracing-go/log"
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
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"os"
	"path/filepath"
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

type CnvrgInfraReconciler struct {
	client.Client
	recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

var infraLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras/status,verbs=get;update;patch

func (r *CnvrgInfraReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	infraLog = r.Log.WithValues("name", req.NamespacedName)
	infraLog.Info("starting cnvrginfra reconciliation")

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

	r.updateStatusMessage(mlopsv1.StatusReconciling, "reconciling", cnvrgInfra)

	// apply manifests
	if err := r.applyManifests(cnvrgInfra); err != nil {
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.StatusHealthy, "successfully reconciled", cnvrgInfra)
	infraLog.Info("successfully reconciled")
	return ctrl.Result{}, nil
}

func (r *CnvrgInfraReconciler) applyManifests(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	var reconcileResult error

	// apply priority classes
	if err := desired.Apply(priorityclass.State(), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// registry
	infraLog.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data: map[string]interface{}{
			"Registry":    cnvrgInfra.Spec.Registry,
			"Annotations": cnvrgInfra.Spec.Annotations,
			"Labels":      cnvrgInfra.Spec.Labels,
		},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// infra reconciler trigger configmap
	if err := r.createInfraReconcilerTriggerCm(cnvrgInfra); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// config reloader
	infraLog.Info("applying config reloader")
	if err := desired.Apply(reloader.State(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// storage
	infraLog.Info("applying storage")
	if err := desired.Apply(storage.State(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// redis
	if cnvrgInfra.Spec.Dbs.Redis.Enabled || cnvrgInfra.Spec.SSO.Enabled {
		redisSecretData := desired.TemplateData{
			Data: map[string]interface{}{
				"Namespace":   cnvrgInfra.Spec.InfraNamespace,
				"Annotations": cnvrgInfra.Spec.Annotations,
				"Labels":      cnvrgInfra.Spec.Labels,
				"CredsRef":    cnvrgInfra.Spec.Dbs.Redis.CredsRef,
				"SvcName":     cnvrgInfra.Spec.Dbs.Redis.SvcName,
			},
		}
		infraLog.Info("trying to generate redis creds (if still doesn't exists...)")
		if err := desired.Apply(dbs.RedisCreds(redisSecretData), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			return err
		}
		infraLog.Info("applying redis")
		if err := desired.Apply(dbs.InfraDbsState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	// logging
	infraLog.Info("applying logging")
	cnvrgApps, err := r.getCnvrgAppInstances()
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	fluentbitData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data: map[string]interface{}{
			"AppInstance":           cnvrgApps,
			"Annotations":           cnvrgInfra.Spec.Annotations,
			"Labels":                cnvrgInfra.Spec.Labels,
			"ClusterInternalDomain": cnvrgInfra.Spec.ClusterInternalDomain,
			"CriType":               cnvrgInfra.Spec.Cri,
		},
	}
	if err := desired.Apply(logging.FluentbitConfigurationState(fluentbitData), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	if err := desired.Apply(logging.InfraLoggingState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// istio
	infraLog.Info("applying infra networking")
	if err := desired.Apply(networking.InfraNetworkingState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// monitoring
	infraLog.Info("applying monitoring")
	if err := r.monitoringState(cnvrgInfra); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// mpi infra
	infraLog.Info("applying mpi infra")
	if err := desired.Apply(controlplane.MpiInfraState(), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// nvidia device plugin
	if cnvrgInfra.Spec.Gpu.NvidiaDp.Enabled {
		infraLog.Info("nvidia device plugin")
		nvidiaDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"NvidiaDp":    cnvrgInfra.Spec.Gpu.NvidiaDp,
				"Registry":    cnvrgInfra.Spec.Registry,
				"ImageHub":    cnvrgInfra.Spec.ImageHub,
				"Annotations": cnvrgInfra.Spec.Annotations,
				"Labels":      cnvrgInfra.Spec.Labels,
			},
		}
		if err := desired.Apply(gpu.NvidiaDpState(nvidiaDpData), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	// habana device plugin
	if cnvrgInfra.Spec.Gpu.HabanaDp.Enabled {
		infraLog.Info("habana device plugin")
		habanaDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"HabanaDp":    cnvrgInfra.Spec.Gpu.HabanaDp,
				"Registry":    cnvrgInfra.Spec.Registry,
				"ImageHub":    cnvrgInfra.Spec.ImageHub,
				"Annotations": cnvrgInfra.Spec.Annotations,
				"Labels":      cnvrgInfra.Spec.Labels,
			},
		}
		if err := desired.Apply(gpu.HabanaDpState(habanaDpData), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	// metagpu device plugin
	if cnvrgInfra.Spec.Gpu.MetaGpuDp.Enabled {
		infraLog.Info("metagpu device plugin")
		metagpuDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"Annotations": cnvrgInfra.Spec.Annotations,
				"Labels":      cnvrgInfra.Spec.Labels,
				"MetaGpuDp":   cnvrgInfra.Spec.Gpu.MetaGpuDp,
				"ImageHub":    cnvrgInfra.Spec.ImageHub,
			},
		}
		// apply metagpu infra state
		if err := desired.Apply(gpu.MetagpudpState(metagpuDpData), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}

	}
	// turn on/off metagpu presence cm in each app ns based on the infra state
	if err := r.setMetagpuPresence(cnvrgInfra); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// capsule backup service
	infraLog.Info("applying capsule")
	if err := desired.Apply(capsule.State(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, infraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	return reconcileResult
}

func (r *CnvrgInfraReconciler) getCnvrgAppInstances() ([]*mlopsv1.AppInstance, error) {

	var appInstances []*mlopsv1.AppInstance
	apps := &mlopsv1.CnvrgAppList{}
	var opts []client.ListOption
	if err := r.List(context.Background(), apps, opts...); err != nil {
		return nil, err
	}

	for _, app := range apps.Items {
		esUser, esPass, err := getAppEsCredsSecret(r.Client, &app)
		if err != nil {
			log.Error(err)
			continue
		}
		appInstances = append(appInstances, &mlopsv1.AppInstance{
			SpecName: app.Name, SpecNs: app.Namespace, EsUser: esUser, EsPass: esPass,
		})
	}

	return appInstances, nil
}

func (r *CnvrgInfraReconciler) monitoringState(infra *mlopsv1.CnvrgInfra) error {

	if err := r.generateMonitoringSecrets(infra); err != nil {
		return err
	}
	if err := desired.Apply(monitoring.InfraMonitoringState(infra), infra, r.Client, r.Scheme, infraLog); err != nil {
		return err
	}
	return nil
}

func (r *CnvrgInfraReconciler) generateMonitoringSecrets(infra *mlopsv1.CnvrgInfra) error {

	if infra.Spec.Monitoring.Prometheus.Enabled {
		user := "cnvrg"
		pass := desired.RandomString()
		passHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
		if err != nil {
			infraLog.Error(err, "error generating prometheus hash")
			return err
		}
		promSecretData := desired.TemplateData{
			Data: map[string]interface{}{
				"Namespace":   infra.Spec.InfraNamespace,
				"Annotations": infra.Spec.Annotations,
				"Labels":      infra.Spec.Labels,
				"CredsRef":    infra.Spec.Monitoring.Prometheus.CredsRef,
				"User":        user,
				"Pass":        pass,
				"PassHash":    fmt.Sprintf("%s:%s", user, passHash),
				"PromUrl":     fmt.Sprintf("http://%s.%s.svc:%d", infra.Spec.Monitoring.Prometheus.SvcName, infra.Spec.InfraNamespace, infra.Spec.Monitoring.Prometheus.Port),
			},
		}
		infraLog.Info("trying to generate prometheus creds (if still doesn't exists...)")
		if err := desired.Apply(monitoring.PromCreds(promSecretData), infra, r.Client, r.Scheme, infraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), infra)
			return err
		}
	}

	if infra.Spec.Monitoring.Grafana.Enabled {
		// grafana dashboards
		infraLog.Info("applying grafana dashboards")
		if err := r.createGrafanaDashboards(infra); err != nil {
			return err
		}

		// grafana datasource
		infraLog.Info("applying grafana datasource")
		url, basicAuthUser, basicAuthPass, err := desired.GetPromCredsSecret(infra.Spec.Monitoring.Prometheus.CredsRef, infra.Spec.InfraNamespace, r.Client, infraLog)
		if err != nil {
			return err
		}
		grafanaDatasourceData := desired.TemplateData{
			Namespace: infra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"Url":         url,
				"User":        basicAuthUser,
				"Pass":        basicAuthPass,
				"Annotations": infra.Spec.Annotations,
				"Labels":      infra.Spec.Labels,
			},
		}
		if err := desired.Apply(monitoring.GrafanaDSState(grafanaDatasourceData), infra, r.Client, r.Scheme, infraLog); err != nil {
			return err
		}
	}

	return nil
}

func (r *CnvrgInfraReconciler) createGrafanaDashboards(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	if !cnvrgInfra.Spec.Monitoring.Grafana.Enabled {
		infraLog.Info("grafana disabled, skipping grafana deployment")
		return nil
	}

	basePath := "/pkg/monitoring/tmpl/grafana/dashboards-data/"
	for _, dashboard := range desired.GrafanaInfraDashboards {
		f, err := pkger.Open(basePath + dashboard)
		if err != nil {
			infraLog.Error(err, "error reading path", "path", dashboard)
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			infraLog.Error(err, "error reading", "file", dashboard)
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
			infraLog.Error(err, "error setting controller reference", "file", f.Name())
			return err
		}
		if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
			infraLog.V(1).Info("grafana dashboard already exists", "file", dashboard)
			continue
		} else if err != nil {
			infraLog.Error(err, "error reading", "file", dashboard)
			return err
		}
	}

	return nil

}

func (r *CnvrgInfraReconciler) syncCnvrgInfraSpec(name types.NamespacedName) (bool, error) {

	infraLog.Info("synchronizing cnvrgInfra spec")

	// Fetch current cnvrgInfra spec
	cnvrgInfra, err := r.getCnvrgInfraSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgInfra == nil {
		return true, nil // all (probably) good, cnvrginfra was removed
	}
	infraLog = r.Log.WithValues("name", name, "ns", cnvrgInfra.Spec.InfraNamespace)

	// Get default cnvrgInfra spec
	desiredSpec := mlopsv1.DefaultCnvrgInfraSpec()

	if err := calculateAndApplyInfraDefaults(cnvrgInfra, &desiredSpec, r.Client); err != nil {
		infraLog.Error(err, "can't calculate defaults")
		return false, err
	}

	// Merge current cnvrgInfra spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgInfra.Spec, mergo.WithOverride, mergo.WithTransformers(cnvrgSpecBoolTransformer{})); err != nil {
		infraLog.Error(err, "can't merge")
		return false, err
	}

	infraLog.V(1).Info("printing the diff between desiredSpec and actual")
	diff, _ := messagediff.PrettyDiff(desiredSpec, cnvrgInfra.Spec)
	infraLog.V(1).Info(diff)

	// Compare desiredSpec and current cnvrgInfra spec,
	// if they are not equal, update the cnvrgInfra spec with desiredSpec,
	// and return true for triggering new reconciliation
	equal := reflect.DeepEqual(desiredSpec, cnvrgInfra.Spec)
	if !equal {
		infraLog.Info("states are not equals, syncing and requeuing")
		cnvrgInfra.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgInfra); err != nil && errors.IsConflict(err) {
			infraLog.Error(err, "conflict updating cnvrgInfra object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	infraLog.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgInfraReconciler) getCnvrgInfraSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgInfra, error) {
	ctx := context.Background()
	var cnvrgInfra mlopsv1.CnvrgInfra
	if err := r.Get(ctx, namespacedName, &cnvrgInfra); err != nil {
		if errors.IsNotFound(err) {
			infraLog.Info("unable to fetch CnvrgInfra, probably cr was deleted")
			return nil, nil
		}
		infraLog.Error(err, "unable to fetch CnvrgInfra")
		return nil, err
	}
	return &cnvrgInfra, nil
}

func (r *CnvrgInfraReconciler) cleanupIstio(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	infraLog.Info("running istio cleanup")
	ctx := context.Background()
	istioManifests := networking.InfraNetworkingState(cnvrgInfra)
	for _, m := range istioManifests {
		// Make sure IstioOperator was deployed
		if m.GVK == desired.Kinds[desired.IstioGVK] {
			if m.TemplateData == nil {
				m.TemplateData = cnvrgInfra
			}
			if err := m.GenerateDeployable(); err != nil {
				infraLog.Error(err, "can't make manifest deployable")
				return err
			}
			if err := r.Delete(ctx, m.Obj); err != nil {
				if errors.IsNotFound(err) {
					infraLog.Info("istio instance not found - probably removed previously")
					return nil
				}
				return err
			}
			istioExists := true
			infraLog.Info("wait for istio instance removal")
			for istioExists {
				err := r.Get(ctx, types.NamespacedName{Name: m.Obj.GetName(), Namespace: m.Obj.GetNamespace()}, m.Obj)
				if err != nil && errors.IsNotFound(err) {
					infraLog.Info("istio instance was successfully removed")
					istioExists = false
				}
				if istioExists {
					infraLog.Info("istio instance still present, will sleep of 1 sec, and check again...")
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	return nil
}

func (r *CnvrgInfraReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgInfra *mlopsv1.CnvrgInfra) {
	if cnvrgInfra.Status.Status == mlopsv1.StatusRemoving {
		infraLog.Info("skipping status update, current cnvrg spec under removing status...")
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
		infraLog.Error(err, "can't update status")
	}
}

func (r *CnvrgInfraReconciler) setMetagpuPresence(infra *mlopsv1.CnvrgInfra) error {
	apps, err := r.getCnvrgAppInstances()
	if err != nil {
		return err
	}
	for _, app := range apps {
		mgDpPresence := desired.TemplateData{
			Namespace: app.SpecNs,
			Data: map[string]interface{}{
				"Annotations": infra.Spec.Annotations,
				"Labels":      infra.Spec.Labels,
				"Enabled":     infra.Spec.Gpu.MetaGpuDp.Enabled,
			},
		}
		if err := desired.Apply(gpu.MetagpudpPresenceState(mgDpPresence), infra, r.Client, r.Scheme, infraLog); err != nil {
			return err
		}
	}
	return nil
}

func (r *CnvrgInfraReconciler) createInfraReconcilerTriggerCm(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: mlopsv1.InfraReconcilerCm, Namespace: cnvrgInfra.Spec.InfraNamespace},
	}
	if err := ctrl.SetControllerReference(cnvrgInfra, cm, r.Scheme); err != nil {
		infraLog.Error(err, "failed to set ControllerReference", "cm", mlopsv1.InfraReconcilerCm)
		return err
	}
	if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
		infraLog.Info("already exists", "cm", mlopsv1.InfraReconcilerCm)
	} else if err != nil {
		infraLog.Error(err, "error creating", "cm", mlopsv1.InfraReconcilerCm)
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	infraLog = r.Log.WithValues("initializing", "crds")

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Info("deploy-depended-crds is false, I hope CRDs was deployed ahead and match expected versions, if not I will fail...")
	} else {

		if viper.GetBool("own-istio-resources") {
			err := desired.Apply(networking.IstioCrds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r.Client, r.Scheme, r.Log)
			if err != nil {
				infraLog.Error(err, "can't apply istio CRDs")
				os.Exit(1)
			}
		}

		if viper.GetBool("own-prometheus-resources") {
			err := desired.Apply(monitoring.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r.Client, r.Scheme, r.Log)
			if err != nil {
				infraLog.Error(err, "can't apply prometheus CRDs")
				os.Exit(1)
			}
		}

		err := desired.Apply(controlplane.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r.Client, r.Scheme, r.Log)
		if err != nil {
			infraLog.Error(err, "can't apply control plane crds")
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
			infraLog.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
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
			infraLog.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return !deleteEvent.DeleteStateUnknown
		},
	}

	infraOwnsPredicate := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {
			infraLog.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			if e.ObjectNew.GetName() == mlopsv1.InfraReconcilerCm {
				// Infra reconciler ConfigMap should always trigger reconcile loop
				return true
			}

			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			infraLog.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
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

	infraLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgInfraController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}
