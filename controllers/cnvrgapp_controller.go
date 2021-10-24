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
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strconv"
	"strings"
	"time"
)

const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"

type CnvrgAppReconciler struct {
	client.Client
	recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

var appLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

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
			if err := r.Update(ctx, cnvrgApp); err != nil {
				appLog.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusRemoving, Message: "removing cnvrg spec"}, cnvrgApp)
			if err := r.cleanup(cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				if err := r.Update(ctx, cnvrgApp); err != nil {
					cnvrgApp, err := r.getCnvrgAppSpec(req.NamespacedName)
					if err != nil {
						appLog.Error(err, "error getting cnvrgapp for finalizer cleanup")
						return err
					}
					cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
					return r.Update(ctx, cnvrgApp)
				}
				return err
			})
			if err != nil {
				appLog.Info("error in removing finalizer")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// check if enabled control plane workloads are all in ready status
	ready, percentageReady, stackReadiness, err := r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	// even if all control plane workloads are ready, let operator finish the full reconcile loop
	if percentageReady == 100 {
		percentageReady = 99
	}
	s := mlopsv1.Status{
		Status:         mlopsv1.StatusReconciling,
		Message:        fmt.Sprintf("reconciling... (%d%%)", percentageReady),
		Progress:       percentageReady,
		StackReadiness: stackReadiness}
	r.updateStatusMessage(s, cnvrgApp)

	// apply spec manifests
	if err := r.applyManifests(cnvrgApp); err != nil {
		return ctrl.Result{}, err
	}

	// get control plan readiness
	ready, percentageReady, stackReadiness, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}
	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", percentageReady)
	appLog.Info(statusMsg)

	if ready { // ura, done
		s := mlopsv1.Status{
			Status:         mlopsv1.StatusReady,
			Message:        statusMsg,
			Progress:       percentageReady,
			StackReadiness: stackReadiness}
		r.updateStatusMessage(s, cnvrgApp)
		appLog.Info("stack is ready!")
		r.recorder.Event(cnvrgApp, "Normal", "Created", fmt.Sprintf("cnvrgapp %s successfully deployed", req.NamespacedName))
		return ctrl.Result{}, nil
	} else { // reconcile again
		requeueAfter, _ := time.ParseDuration("30s")
		appLog.Info("stack not ready yet, requeuing...")
		r.recorder.Event(cnvrgApp, "Normal", "Creating", fmt.Sprintf("cnvrgapp %s not ready yet, done: %d%%", req.NamespacedName, percentageReady))
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

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(cnvrgApp *mlopsv1.CnvrgApp) (bool, int, map[string]bool, error) {

	readyState := make(map[string]bool)

	// check webapp status
	if cnvrgApp.Spec.ControlPlane.WebApp.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.ControlPlane.WebApp.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["webApp"] = ready
	}

	// check sidekiq status
	if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled {
		name := types.NamespacedName{Name: "sidekiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["sidekiq"] = ready
	}

	// check searchkiq status
	if cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled {
		name := types.NamespacedName{Name: "searchkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["searchkiq"] = ready
	}

	// check systemkiq status
	if cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled {
		name := types.NamespacedName{Name: "systemkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["systemkiq"] = ready
	}

	// check postgres status
	if cnvrgApp.Spec.Dbs.Pg.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Pg.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["pg"] = ready
	}

	// check minio status
	if cnvrgApp.Spec.Dbs.Minio.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Minio.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["minio"] = ready
	}

	// check redis status
	if cnvrgApp.Spec.Dbs.Redis.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Redis.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["redis"] = ready
	}

	// check es status
	if cnvrgApp.Spec.Dbs.Es.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Es.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		// if es is ready, trigger fluentbit reconfiguration
		if ready {
			appLog.Info("es is ready, triggering fluentbit reconfiguration")
			if err := r.addFluentbitConfiguration(cnvrgApp); err != nil {
				return false, 0, nil, err
			}
		}
		readyState["es"] = ready
	}

	// check kibana status
	if cnvrgApp.Spec.Logging.Kibana.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Logging.Kibana.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["kibana"] = ready
	}

	// check prometheus status
	if cnvrgApp.Spec.Monitoring.Prometheus.Enabled {
		name := types.NamespacedName{Name: "prometheus-cnvrg-ccp-prometheus", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["prometheus"] = ready
	}

	// check prometheus status
	if cnvrgApp.Spec.Monitoring.Grafana.Enabled {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Monitoring.Grafana.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["grafana"] = ready
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

	return readyCount == len(readyState), percentageReady, readyState, nil
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
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// dbs
	if err := r.dbsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// backups
	if err := r.backupsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// networking
	appLog.Info("applying networking")
	if err := desired.Apply(networking.CnvrgAppNetworkingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// logging
	if err := r.loggingState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// controlplane
	appLog.Info("applying controlplane")
	if err := desired.Apply(controlplane.State(cnvrgApp), cnvrgApp, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	// monitoring
	if err := r.monitoringState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, cnvrgApp)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) loggingState(app *mlopsv1.CnvrgApp) error {

	if app.Spec.Logging.Kibana.Enabled {
		appLog.Info("applying logging")
		kibanaConfigSecretData, err := r.getKibanaConfigSecretData(app)
		if err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
		if err := desired.Apply(logging.KibanaConfSecret(*kibanaConfigSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
		if err := desired.Apply(logging.CnvrgAppLoggingState(app), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}

	return nil
}

func (r *CnvrgAppReconciler) dbsState(app *mlopsv1.CnvrgApp) error {

	// dbs
	appLog.Info("applying dbs")

	if app.Spec.Dbs.Es.Enabled {
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
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}

	if app.Spec.Dbs.Pg.Enabled {
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
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}

	if app.Spec.Dbs.Redis.Enabled {
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
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}

	if err := desired.Apply(dbs.AppDbsState(app), app, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) backupsState(app *mlopsv1.CnvrgApp) error {

	if app.Spec.Dbs.Pg.Enabled { // pg backups
		pgPvc := v1core.PersistentVolumeClaim{}
		pgPvcName := types.NamespacedName{Namespace: app.Namespace, Name: app.Spec.Dbs.Pg.PvcName}
		if err := r.Get(context.Background(), pgPvcName, &pgPvc); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
		if err := r.ApplyCapsuleAnnotations(app.Spec.Dbs.Pg.Backup, &pgPvc, "postgresql"); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) monitoringState(app *mlopsv1.CnvrgApp) error {

	// generate monitoring secrets (prometheus, prometheus upstream, and grafana data sources
	if err := r.generateMonitoringSecrets(app); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
		return err
	}
	// apply app monitoring state
	if err := desired.Apply(monitoring.AppMonitoringState(app), app, r.Client, r.Scheme, appLog); err != nil {
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) generateMonitoringSecrets(app *mlopsv1.CnvrgApp) error {

	if app.Spec.Monitoring.Prometheus.Enabled {
		appLog.Info("applying monitoring")
		user := "cnvrg"
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
				"User":        user,
				"Pass":        pass,
				"PassHash":    fmt.Sprintf("%s:%s", user, passHash),
				"PromUrl":     fmt.Sprintf("http://%s.%s.svc:%d", app.Spec.Monitoring.Prometheus.SvcName, app.Namespace, app.Spec.Monitoring.Prometheus.Port),
			},
		}

		appLog.Info("trying to generate prometheus creds (if still doesn't exists...)")
		if err := desired.Apply(monitoring.PromCreds(promSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}

		if err := r.createUpstreamPrometheusConfig(app); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}

	if app.Spec.Monitoring.Grafana.Enabled {
		// grafana dashboards
		appLog.Info("applying grafana dashboards ")
		if err := r.createGrafanaDashboards(app); err != nil {
			return err
		}
		// grafana datasource
		appLog.Info("applying grafana datasource")
		url, user, pass, err := desired.GetPromCredsSecret(app.Spec.Monitoring.Prometheus.CredsRef, app.Namespace, r.Client, appLog)
		if err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
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
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
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
		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
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
	if app.Spec.SSO.Enabled {
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

	if !cnvrgApp.Spec.Monitoring.Grafana.Enabled {
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

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.Status, app *mlopsv1.CnvrgApp) {

	if app.Status.Status == mlopsv1.StatusRemoving {
		appLog.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	if status.Status == mlopsv1.StatusError {
		msg := fmt.Sprintf("%s/%s error acoured during reconcile", app.GetNamespace(), app.GetName())
		r.recorder.Event(app, "Warning", "ReconcileError", msg)
	}
	ctx := context.Background()
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		name := types.NamespacedName{Namespace: app.Namespace, Name: app.Name}
		app, err := r.getCnvrgAppSpec(name)
		if err != nil {
			return err
		}
		app.Status.Status = status.Status
		app.Status.Message = status.Message
		if status.Progress >= 0 {
			app.Status.Progress = status.Progress
		}
		if status.StackReadiness != nil {
			app.Status.StackReadiness = status.StackReadiness
		}
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

	infra, err := r.getCnvrgInfra()
	if err != nil {
		appLog.Error(err, "can't get cnvrg infra")
		return false, err
	}

	calculateAndApplyAppDefaults(cnvrgApp, &desiredSpec, infra)

	// Merge current cnvrgApp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgApp.Spec, mergo.WithOverride, mergo.WithTransformers(cnvrgSpecBoolTransformer{})); err != nil {
		appLog.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {

		if diff, equal := messagediff.PrettyDiff(desiredSpec, cnvrgApp.Spec); !equal {
			appLog.Info("diff between desiredSpec and actual")
			appLog.Info(diff)
		}

		if diff, equal := messagediff.PrettyDiff(cnvrgApp.Spec, desiredSpec); !equal {
			appLog.Info("diff between actual and desired")
			appLog.Info(diff)
		}

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

	appPredicate := predicate.Funcs{

		CreateFunc: func(createEvent event.CreateEvent) bool {
			msg := fmt.Sprintf("cnvrgapp: %s/%s has been created", createEvent.Object.GetNamespace(), createEvent.Object.GetName())
			r.recorder.Event(createEvent.Object, "Normal", "Created", msg)
			return true
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			infraLog.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			shouldReconcile := e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
			if shouldReconcile {
				msg := fmt.Sprintf("cnvrgapp: %s/%s has been updated", e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
				r.recorder.Event(e.ObjectNew, "Normal", "Updated", msg)
			}
			return shouldReconcile
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			msg := fmt.Sprintf("cnvrgapp: %s/%s has been deleted", deleteEvent.Object.GetNamespace(), deleteEvent.Object.GetName())
			r.recorder.Event(deleteEvent.Object, "Normal", "SuccessfulDelete", msg)
			infraLog.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return !deleteEvent.DeleteStateUnknown
		},
	}

	appOwnsPredicate := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {
			infraLog.V(1).Info("received update event", "objectName", e.ObjectNew.GetName())
			return false
		},

		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			infraLog.V(1).Info("received delete event", "objectName", deleteEvent.Object.GetName())
			return true
		},
	}
	r.recorder = mgr.GetEventRecorderFor("cnvrgapp")
	cnvrgAppController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgApp{}, builder.WithPredicates(appPredicate))

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
		cnvrgAppController.Owns(u, builder.WithPredicates(appOwnsPredicate))
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

func (r *CnvrgAppReconciler) ApplyCapsuleAnnotations(b mlopsv1.Backup, pvc *v1core.PersistentVolumeClaim, serviceType string) error {
	if pvc.Annotations == nil {
		pvc.Annotations = map[string]string{}
	}
	pvc.Annotations["capsule.mlops.cnvrg.io/backup"] = "false"
	if b.Enabled {
		pvc.Annotations["capsule.mlops.cnvrg.io/backup"] = "true"
	}
	pvc.Annotations["capsule.mlops.cnvrg.io/serviceType"] = serviceType
	pvc.Annotations["capsule.mlops.cnvrg.io/bucketRef"] = b.BucketRef
	pvc.Annotations["capsule.mlops.cnvrg.io/credsRef"] = b.CredsRef
	pvc.Annotations["capsule.mlops.cnvrg.io/rotation"] = strconv.Itoa(b.Rotation)
	pvc.Annotations["capsule.mlops.cnvrg.io/period"] = b.Period

	if err := r.Update(context.Background(), pvc); err != nil {
		return err
	}
	return nil
}
