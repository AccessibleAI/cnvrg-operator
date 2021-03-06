package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/dbs"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/logging"
	"github.com/AccessibleAI/cnvrg-operator/pkg/monitoring"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/registry"
	"github.com/Dimss/crypt/apr1_crypt"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"gopkg.in/d4l3k/messagediff.v1"
	"io/ioutil"
	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sort"
	"strconv"
	"strings"
	"time"
)

const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"

type CnvrgAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var appLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	appLog = r.Log.WithValues("name", req.NamespacedName)
	appLog.Info("starting cnvrgapp reconciliation")

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

	// setup finalizer
	if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			cnvrgApp.ObjectMeta.Finalizers = append(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), cnvrgApp); err != nil {
				appLog.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			r.updateStatusMessage(mlopsv1.StatusRemoving, "removing cnvrg spec", cnvrgApp)
			if err := r.cleanup(cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgInfra, err := r.getCnvrgAppSpec(req.NamespacedName)
			if err != nil {
				return ctrl.Result{}, err
			}
			if cnvrgInfra == nil {
				return ctrl.Result{}, nil
			}
			cnvrgInfra.ObjectMeta.Finalizers = removeString(cnvrgInfra.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				appLog.Info("error in removing finalizer, checking if cnvrgInfra object still exists")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// check if enabled control plane workloads are all in ready status
	ready, percentageReady, err := r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	// even if all control plane workloads are ready, let operator finish the full reconcile loop
	if percentageReady == 100 {
		percentageReady = 99
	}
	r.updateStatusMessage(mlopsv1.StatusReconciling, fmt.Sprintf("reconciling... (%d%%)", percentageReady), cnvrgApp)

	// apply spec manifests
	if err := r.applyManifests(cnvrgApp); err != nil {
		return ctrl.Result{}, err
	}

	// get control plan readiness
	ready, percentageReady, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}
	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", percentageReady)
	appLog.Info(statusMsg)

	if ready { // ura, done
		r.updateStatusMessage(mlopsv1.StatusReady, statusMsg, cnvrgApp)
		appLog.Info("stack is ready!")
		return ctrl.Result{}, nil
	} else { // reconcile again
		requeueAfter, err := time.ParseDuration("30s")
		if err != nil {
			appLog.Error(err, "wrong duration for requeueAfter")
			return ctrl.Result{}, err
		}
		appLog.Info("stack not ready yet, requeuing...")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
}

func (r *CnvrgAppReconciler) getEsCredsSecret(app *mlopsv1.CnvrgApp) (user string, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: app.Spec.Dbs.Es.CredsRef, Namespace: app.Namespace}
	creds := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := r.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		appLog.Error(err, "es-creds secret not found!")
		return "", "", err
	} else if err != nil {
		appLog.Error(err, "can't check if es creds secret exists", "name", namespacedName.Name)
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_USER"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_USER", namespacedName.Name)
		appLog.Error(err, "missing required field")
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_PASS"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_PASS", namespacedName.Name)
		appLog.Error(err, "missing required field")
		return "", "", err
	}

	return string(creds.Data["CNVRG_ES_USER"]), string(creds.Data["CNVRG_ES_PASS"]), nil
}

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(cnvrgApp *mlopsv1.CnvrgApp) (bool, int, error) {

	readyState := make(map[string]bool)

	// check webapp status
	if *cnvrgApp.Spec.ControlPlane.WebApp.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.ControlPlane.WebApp.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["webApp"] = ready
	}

	// check sidekiq status
	if *cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled {
		name := types.NamespacedName{Name: "sidekiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["sidekiq"] = ready
	}

	// check searchkiq status
	if *cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled {
		name := types.NamespacedName{Name: "searchkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["searchkiq"] = ready
	}

	// check systemkiq status
	if *cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled {
		name := types.NamespacedName{Name: "systemkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["searchkiq"] = ready
	}

	// check postgres status
	if *cnvrgApp.Spec.Dbs.Pg.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Pg.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["pg"] = ready
	}

	// check minio status
	if *cnvrgApp.Spec.Dbs.Minio.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Minio.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["minio"] = ready
	}

	// check redis status
	if *cnvrgApp.Spec.Dbs.Redis.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Redis.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["redis"] = ready
	}

	// check es status
	if *cnvrgApp.Spec.Dbs.Es.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Es.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, err
		}
		// if es is ready, trigger fluentbit reconfiguration
		if ready {
			appLog.Info("es is ready, triggering fluentbit reconfiguration")
			if err := r.addFluentbitConfiguration(cnvrgApp); err != nil {
				return false, 0, err
			}
		}
		readyState["es"] = ready
	}

	// check kibana status
	if *cnvrgApp.Spec.Logging.Kibana.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Logging.Kibana.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["kibana"] = ready
	}

	percentageReady := 0

	readyCount := 0

	for _, ready := range readyState {
		if ready {
			readyCount++
		}
	}

	if len(readyState) > 0 {
		percentageReady = readyCount * 100 / len(readyState)
	}

	return readyCount == len(readyState), percentageReady, nil
}

func (r *CnvrgAppReconciler) applyManifests(cnvrgApp *mlopsv1.CnvrgApp) error {

	// registry
	appLog.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgApp.Namespace,
		Data: map[string]interface{}{
			"Registry":    cnvrgApp.Spec.Registry,
			"Annotations": cnvrgApp.Spec.Annotations,
			"Labels":      cnvrgApp.Spec.Labels,
		},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// dbs
	if err := r.dbsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// networking
	appLog.Info("applying networking")
	if err := desired.Apply(networking.CnvrgAppNetworkingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// logging
	if err := r.loggingState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// controlplane
	appLog.Info("applying controlplane")
	if err := desired.Apply(controlplane.State(cnvrgApp), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// monitoring
	if err := r.monitoringState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) loggingState(app *mlopsv1.CnvrgApp) error {

	if *app.Spec.Logging.Kibana.Enabled {
		appLog.Info("applying logging")
		kibanaConfigSecretData, err := r.getKibanaConfigSecretData(app)
		if err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
		if err := desired.Apply(logging.KibanaConfSecret(*kibanaConfigSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
		if err := desired.Apply(logging.CnvrgAppLoggingState(app), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
	}

	return nil
}

func (r *CnvrgAppReconciler) dbsState(app *mlopsv1.CnvrgApp) error {

	// dbs
	appLog.Info("applying dbs")

	if *app.Spec.Dbs.Es.Enabled {
		esSecretData := desired.TemplateData{
			Data: map[string]interface{}{
				"Namespace":   app.Namespace,
				"CredsRef":    app.Spec.Dbs.Es.CredsRef,
				"EsUrl":       fmt.Sprintf("%s.%s.svc:%d", app.Spec.Dbs.Es.SvcName, app.Namespace, app.Spec.Dbs.Es.Port),
				"Annotations": app.Spec.Annotations,
				"Labels":      app.Spec.Labels,
			},
		}
		appLog.Info("trying to generate es creds (if still doesn't exists...)")
		if err := desired.Apply(dbs.EsCreds(esSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
	}

	if *app.Spec.Dbs.Pg.Enabled {
		pgSecretData := desired.TemplateData{
			Data: map[string]interface{}{
				"Namespace":          app.Namespace,
				"CredsRef":           app.Spec.Dbs.Pg.CredsRef,
				"Annotations":        app.Spec.Annotations,
				"Labels":             app.Spec.Labels,
				"MaxConnections":     app.Spec.Dbs.Pg.MaxConnections,
				"SharedBuffers":      app.Spec.Dbs.Pg.SharedBuffers,
				"EffectiveCacheSize": app.Spec.Dbs.Pg.EffectiveCacheSize,
				"SvcName":            app.Spec.Dbs.Pg.SvcName,
			},
		}
		appLog.Info("trying to generate pg creds (if still doesn't exists...)")
		if err := desired.Apply(dbs.PgCreds(pgSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
	}

	if *app.Spec.Dbs.Redis.Enabled {
		redisSecretData := desired.TemplateData{
			Data: map[string]interface{}{
				"Namespace":   app.Namespace,
				"Annotations": app.Spec.Annotations,
				"Labels":      app.Spec.Labels,
				"CredsRef":    app.Spec.Dbs.Redis.CredsRef,
				"SvcName":     app.Spec.Dbs.Redis.SvcName,
			},
		}
		appLog.Info("trying to generate redis creds (if still doesn't exists...)")
		if err := desired.Apply(dbs.RedisCreds(redisSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
	}

	if err := desired.Apply(dbs.AppDbsState(app), app, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) monitoringState(app *mlopsv1.CnvrgApp) error {

	// generate monitoring secrets (prometheus, prometheus upstream, and grafana data sources
	if err := r.generateMonitoringSecrets(app); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	// apply app monitoring state
	if err := desired.Apply(monitoring.AppMonitoringState(app), app, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) generateMonitoringSecrets(app *mlopsv1.CnvrgApp) error {
	if *app.Spec.Monitoring.Prometheus.Enabled {
		appLog.Info("applying monitoring")
		pass := desired.RandomString()
		passHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
		if err != nil {
			appLog.Error(err, "error generating prometheus hash")
			return err
		}
		promSecretData := desired.TemplateData{
			Data: map[string]interface{}{
				"Namespace":   app.Namespace,
				"Annotations": app.Spec.Annotations,
				"Labels":      app.Spec.Labels,
				"CredsRef":    app.Spec.Monitoring.Prometheus.CredsRef,
				"User":        "cnvrg",
				"Pass":        pass,
				"PassHash":    passHash,
				"PromUrl":     fmt.Sprintf("http://%s.%s.svc:%d", app.Spec.Monitoring.Prometheus.SvcName, app.Namespace, app.Spec.Monitoring.Prometheus.Port),
			},
		}

		appLog.Info("trying to generate prometheus creds (if still doesn't exists...)")
		if err := desired.Apply(monitoring.PromCreds(promSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}

		if err := r.createUpstreamPrometheusConfig(app); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
	}

	if *app.Spec.Monitoring.Grafana.Enabled {
		// grafana dashboards
		appLog.Info("applying grafana dashboards ")
		if err := r.createGrafanaDashboards(app); err != nil {
			return err
		}
		// grafana datasource
		appLog.Info("applying grafana datasource")
		url, user, pass, err := desired.GetPromCredsSecret(app.Spec.Monitoring.Prometheus.CredsRef, app.Namespace, r, appLog)
		if err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
		grafanaDatasourceData := desired.TemplateData{
			Namespace: app.Namespace,
			Data: map[string]interface{}{
				"Url":  url,
				"User": user,
				"Pass": pass,
			},
		}

		if err := desired.Apply(monitoring.GrafanaDSState(grafanaDatasourceData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
			return err
		}
	}

	return nil
}

func (r *CnvrgAppReconciler) getCnvrgInfra() (*mlopsv1.CnvrgInfra, error) {

	cnvrgAppInfra := &mlopsv1.CnvrgInfraList{}

	if err := r.List(context.Background(), cnvrgAppInfra); err != nil {
		appLog.Error(err, "can't list CnvrgInfra objects")
		return nil, err
	}

	if len(cnvrgAppInfra.Items) == 0 {
		appLog.Info("no CnvrgInfra objects was deployed, skipping infra reconciler")
		return nil, errors.NewNotFound(schema.GroupResource{Group: "mlops.cnvrg.io", Resource: "CnvrgInfra"}, "cnvrg-infra")
	}

	return &cnvrgAppInfra.Items[0], nil
}

func (r *CnvrgAppReconciler) createUpstreamPrometheusConfig(app *mlopsv1.CnvrgApp) error {
	infra, err := r.getCnvrgInfra()
	if err != nil {
		appLog.Error(err, "can't get cnvrgInfra object ")
		return err
	}

	_, user, pass, err := desired.GetPromCredsSecret(infra.Spec.Monitoring.Prometheus.CredsRef, infra.Spec.InfraNamespace, r.Client, appLog)
	if err != nil {
		appLog.Error(err, "can't get cnvrgInfra prometheus creds")
		return err
	}

	promUpstreamData := desired.TemplateData{
		Data: map[string]interface{}{
			"Namespace":   app.Namespace,
			"Annotations": app.Spec.Annotations,
			"Labels":      app.Spec.Labels,
			"CredsRef":    app.Spec.Monitoring.Prometheus.UpstreamRef,
			"User":        user,
			"Pass":        pass,
			"Upstream":    fmt.Sprintf("prometheus-operated.%s.svc:%d", infra.Spec.InfraNamespace, infra.Spec.Monitoring.Prometheus.Port),
		},
	}

	if err := desired.Apply(monitoring.PromUpstreamCreds(promUpstreamData), app, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) getKibanaConfigSecretData(app *mlopsv1.CnvrgApp) (*desired.TemplateData, error) {
	kibanaHost := "0.0.0.0"
	kibanaPort := strconv.Itoa(app.Spec.Logging.Kibana.Port)
	esUser, esPass, err := r.getEsCredsSecret(app)
	if err != nil {
		appLog.Error(err, "can't fetch es creds")
		return nil, err
	}
	if *app.Spec.SSO.Enabled {
		kibanaHost = "127.0.0.1"
		kibanaPort = "3000"
	}
	return &desired.TemplateData{
		Namespace: app.Namespace,
		Data: map[string]interface{}{
			"Host":        kibanaHost,
			"Port":        kibanaPort,
			"EsHost":      fmt.Sprintf("http://%s.%s.svc:%d", app.Spec.Dbs.Es.SvcName, app.Namespace, app.Spec.Dbs.Es.Port),
			"EsUser":      esUser,
			"EsPass":      esPass,
			"Annotations": app.Spec.Annotations,
			"Labels":      app.Spec.Labels,
		},
	}, nil

}

func (r *CnvrgAppReconciler) createGrafanaDashboards(cnvrgApp *mlopsv1.CnvrgApp) error {

	if !*cnvrgApp.Spec.Monitoring.Grafana.Enabled {
		appLog.Info("grafana disabled, skipping grafana deployment")
		return nil
	}
	basePath := "/pkg/monitoring/tmpl/grafana/dashboards-data/"
	for _, dashboard := range desired.GrafanaAppDashboards {
		f, err := pkger.Open(basePath + dashboard)
		if err != nil {
			appLog.Error(err, "error reading path", "path", dashboard)
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			appLog.Error(err, "error reading", "file", dashboard)
			return err
		}
		cm := &v1core.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name())),
				Namespace: cnvrgApp.Namespace,
			},
			Data: map[string]string{filepath.Base(f.Name()): string(b)},
		}
		if err := ctrl.SetControllerReference(cnvrgApp, cm, r.Scheme); err != nil {
			appLog.Error(err, "error setting controller reference", "file", f.Name())
			return err
		}
		if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
			appLog.V(1).Info("grafana dashboard already exists", "file", dashboard)
			continue
		} else if err != nil {
			appLog.Error(err, "error reading", "file", dashboard)
			return err
		}
	}

	return nil

}

func (r *CnvrgAppReconciler) addFluentbitConfiguration(cnvrgApp *mlopsv1.CnvrgApp) error {
	infra, err := r.getCnvrgInfra()
	if err != nil {
		return err
	}

	name := types.NamespacedName{Name: mlopsv1.InfraReconcilerCm, Namespace: infra.Spec.InfraNamespace}
	infraReconcilerCm := &v1core.ConfigMap{}

	esUser, esPass, err := r.getEsCredsSecret(cnvrgApp)
	if err != nil {
		appLog.Error(err, "failed to fetch es creds")
		return err
	}

	appInstance := mlopsv1.AppInstance{SpecName: cnvrgApp.Name, SpecNs: cnvrgApp.Namespace, EsUser: esUser, EsPass: esPass}
	appInstanceBytes, err := json.Marshal(appInstance)
	if err != nil {
		appLog.Error(err, "failed to marshal app instance ")
		return err
	}
	if err := r.Get(context.Background(), name, infraReconcilerCm); err != nil {
		appLog.Error(err, "can't get reconciler cm", "name", name)
		return err
	}
	if infraReconcilerCm.Data == nil {
		infraReconcilerCm.Data = map[string]string{cnvrgApp.Namespace: string(appInstanceBytes)}
	} else {
		infraReconcilerCm.Data[cnvrgApp.Namespace] = string(appInstanceBytes)
	}
	if err := r.Update(context.Background(), infraReconcilerCm); err != nil {
		appLog.Error(err, "can't update cm", "cm", name)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) removeFluentbitConfiguration(cnvrgApp *mlopsv1.CnvrgApp) error {
	infra, err := r.getCnvrgInfra()
	if err != nil && errors.IsNotFound(err) {
		appLog.Info("cnvrg infra not found, skipping fluentbit cleanup")
		return nil
	} else if err != nil {
		appLog.Info("error getting cnvrg infra, trying reconcile again...")
		return err
	}
	name := types.NamespacedName{Name: mlopsv1.InfraReconcilerCm, Namespace: infra.Spec.InfraNamespace}
	infraReconcilerCm := &v1core.ConfigMap{}
	if err := r.Get(context.Background(), name, infraReconcilerCm); err != nil && errors.IsNotFound(err) {
		appLog.Info("infra reconciler configmap not found, skipping fluentbit cleanup")
		return nil
	} else if err != nil {
		appLog.Error(err, "can't get reconciler cm", "name", name)
		return err
	}
	delete(infraReconcilerCm.Data, cnvrgApp.Namespace)
	if err := r.Update(context.Background(), infraReconcilerCm); err != nil {
		appLog.Error(err, "can't update cm", "cm", name)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgApp *mlopsv1.CnvrgApp) {
	if cnvrgApp.Status.Status == mlopsv1.StatusRemoving {
		appLog.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	ctx := context.Background()
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		name := types.NamespacedName{Namespace: cnvrgApp.Namespace, Name: cnvrgApp.Name}
		app, err := r.getCnvrgAppSpec(name)
		if err != nil {
			return err
		}
		app.Status.Status = status
		app.Status.Message = message
		err = r.Status().Update(ctx, app)
		return err
	})
	if err != nil {
		appLog.Error(err, "can't update status")
	}

}

func (r *CnvrgAppReconciler) syncCnvrgAppSpec(name types.NamespacedName) (bool, error) {

	appLog.Info("synchronizing cnvrgApp spec")

	// Fetch current cnvrgApp spec
	cnvrgApp, err := r.getCnvrgAppSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgApp == nil {
		return false, nil // probably cnvrgapp was removed
	}
	appLog = r.Log.WithValues("name", name, "ns", cnvrgApp.Namespace)

	// Get default cnvrgApp spec
	desiredSpec := mlopsv1.DefaultCnvrgAppSpec()

	calculateAndApplyAppDefaults(cnvrgApp, &desiredSpec)

	// Merge current cnvrgApp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgApp.Spec, mergo.WithOverride); err != nil {
		appLog.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {
		appLog.V(1).Info("printing the diff between desiredSpec and actual")
		diff, _ := messagediff.PrettyDiff(desiredSpec, cnvrgApp.Spec)
		appLog.V(1).Info(diff)
	}

	equal := reflect.DeepEqual(desiredSpec, cnvrgApp.Spec)
	if !equal {
		appLog.Info("states are not equals, syncing and requeuing")
		cnvrgApp.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgApp); err != nil && errors.IsConflict(err) {
			appLog.Info("conflict updating cnvrgApp object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	appLog.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgAppReconciler) getCnvrgAppSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var app mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			appLog.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		appLog.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}

	return &app, nil
}

func (r *CnvrgAppReconciler) cleanup(cnvrgApp *mlopsv1.CnvrgApp) error {

	appLog.Info("running finalizer cleanup")

	// remove cnvrg-db-init
	if err := r.cleanupDbInitCm(cnvrgApp); err != nil {
		return err
	}

	// update infra reconciler cm
	if err := r.removeFluentbitConfiguration(cnvrgApp); err != nil {
		if err.Error() == "no CnvrgInfra objects was deployed, skipping infra reconciler" {
			appLog.Info("cnvrgInfra object not found, no need to trigger infra reconciler")
		} else {
			return err
		}
	}

	// cleanup pvc
	if err := r.cleanupPVCs(); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) cleanupPVCs() error {
	if !viper.GetBool("cleanup-pvc") {
		appLog.Info("cleanup-pvc is false, skipping pvc deletion!")
		return nil
	}
	appLog.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		appLog.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
			if pvc.ObjectMeta.Labels["app"] == "prometheus" || pvc.ObjectMeta.Labels["app"] == "elasticsearch" {
				if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
					appLog.Info("pvc already deleted")
				} else if err != nil {
					appLog.Error(err, "error deleting prometheus pvc")
					return err
				}
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanupDbInitCm(desiredSpec *mlopsv1.CnvrgApp) error {
	appLog.Info("running cnvrg-db-init cleanup")
	ctx := context.Background()
	dbInitCm := &v1core.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-db-init", Namespace: desiredSpec.Namespace}}
	err := r.Delete(ctx, dbInitCm)
	if err != nil && errors.IsNotFound(err) {
		appLog.Info("no need to delete cnvrg-db-init, cm not found")
	} else {
		appLog.Error(err, "error deleting cnvrg-db-init")
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	appLog = r.Log.WithValues("initializing", "crds")

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {
			appLog.V(1).Info("received UpdateEvent", "eventSourcesObjectName", e.MetaNew.GetName())
			if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(e.ObjectOld) {
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgApp)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgApp)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrgapp spec wasn't changed, assuming status update, won't reconcile
				appLog.V(1).Info("update received", "shouldReconcileOnSpecChange", shouldReconcileOnSpecChange)

				return !shouldReconcileOnSpecChange
			}
			return true
		},
	}

	cnvrgAppController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgApp{}).
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
		cnvrgAppController.Owns(u)
	}
	appLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgAppController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
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

func calculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, desiredAppSpec *mlopsv1.CnvrgAppSpec) {
	// set default heap size for ES if not set by user
	if strings.Contains(app.Spec.Dbs.Es.Requests.Memory, "Gi") && app.Spec.Dbs.Es.JavaOpts == "" {
		requestMem := strings.TrimSuffix(app.Spec.Dbs.Es.Requests.Memory, "Gi")
		mem, err := strconv.Atoi(requestMem)
		if err == nil {
			heapMem := mem / 2
			if heapMem > 0 {
				desiredAppSpec.Dbs.Es.JavaOpts = fmt.Sprintf("-Xms%dg -Xmx%dg", heapMem, heapMem)
			}
		}
	}

	if app.Spec.Networking.Ingress.IstioGwName == "" {
		desiredAppSpec.Networking.Ingress.IstioGwName = fmt.Sprintf(mlopsv1.IstioGwName, app.Namespace)
	}

	// if proxy.enabled wasn't set in spec, make sure explicitly set it before initiating noProxy values
	if app.Spec.Networking.Proxy.Enabled == nil {
		app.Spec.Networking.Proxy.Enabled = desiredAppSpec.Networking.Proxy.Enabled
	}

	if *app.Spec.Networking.Proxy.Enabled {
		desiredAppSpec.Networking.Proxy.NoProxy = app.Spec.Networking.Proxy.NoProxy
		// make sure no_proxy includes all default values
		for _, defaultNoProxy := range networking.DefaultNoProxy() {
			if !containsString(desiredAppSpec.Networking.Proxy.NoProxy, defaultNoProxy) {
				desiredAppSpec.Networking.Proxy.NoProxy = append(desiredAppSpec.Networking.Proxy.NoProxy, defaultNoProxy)
			}
		}
		// sort slices before compare
		sort.Strings(desiredAppSpec.Networking.Proxy.NoProxy)
		sort.Strings(app.Spec.Networking.Proxy.NoProxy)
		// if slice are not equal, use desiredAppSpec no_proxy
		if !reflect.DeepEqual(desiredAppSpec.Networking.Proxy.NoProxy, app.Spec.Networking.Proxy.NoProxy) {
			app.Spec.Networking.Proxy.NoProxy = nil
		}
	}

}
