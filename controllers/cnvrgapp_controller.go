package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/controlplane"
	"github.com/cnvrg-operator/pkg/dbs"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/cnvrg-operator/pkg/logging"
	"github.com/cnvrg-operator/pkg/monitoring"
	"github.com/cnvrg-operator/pkg/networking"
	"github.com/cnvrg-operator/pkg/registry"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"io/ioutil"
	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
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
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var cnvrgAppLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	cnvrgAppLog = r.Log.WithValues("name", req.NamespacedName)
	cnvrgAppLog.Info("starting cnvrgapp reconciliation")

	equal, err := r.syncCnvrgAppSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		return ctrl.Result{Requeue: true}, nil
	}

	cnvrgApp, err := r.getCnvrgAppSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgApp == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	// setup finalizer
	if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrginfraFinalizer) {
			cnvrgApp.ObjectMeta.Finalizers = append(cnvrgApp.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			if err := r.Update(context.Background(), cnvrgApp); err != nil {
				cnvrgInfraLog.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrginfraFinalizer) {
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
			cnvrgInfra.ObjectMeta.Finalizers = removeString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				cnvrgInfraLog.Info("error in removing finalizer, checking if cnvrgInfra object still exists")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	ready, percentageReady, err := r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}
	if percentageReady == 100 {
		percentageReady = 99
	}
	r.updateStatusMessage(mlopsv1.StatusReconciling, fmt.Sprintf("reconciling... (%d%%)", percentageReady), cnvrgApp)

	if err := r.applyManifests(cnvrgApp); err != nil {
		return ctrl.Result{}, err
	}

	// get control plan readiness
	ready, percentageReady, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", percentageReady)
	cnvrgAppLog.Info(statusMsg)

	if ready {
		r.updateStatusMessage(mlopsv1.StatusReady, statusMsg, cnvrgApp)
		return ctrl.Result{}, nil
	} else {
		requeueAfter, err := time.ParseDuration("30s")
		if err != nil {
			cnvrgAppLog.Error(err, "wrong duration for requeueAfter")
			return ctrl.Result{}, err
		}
		cnvrgAppLog.Info("stack not ready yet, requeuing...")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
}

func (r *CnvrgAppReconciler) esCredsSecret(app *mlopsv1.CnvrgApp) (user string, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: app.Spec.Dbs.Es.CredsRef, Namespace: app.Namespace}
	creds := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := r.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(app, &creds, r.Scheme); err != nil {
			cnvrgAppLog.Error(err, "error set controller reference", "name", namespacedName.Name)
			return "", "", err
		}

		pass = desired.RandomString()
		creds.Data = map[string][]byte{
			"CNVRG_ES_USER":          []byte(user), // envs for webapp/sidekiq
			"CNVRG_ES_PASS":          []byte(pass), // envs for webapp/sidekiq
			"ES_USERNAME":            []byte(user), // envs for elastalerts
			"ES_PASSWORD":            []byte(pass), // envs for elastalerts
			"ELASTICSEARCH_USERNAME": []byte(user), // envs for kibana
			"ELASTICSEARCH_PASSWORD": []byte(pass), // envs for kibana
		}
		if err := r.Create(context.Background(), &creds); err != nil {
			cnvrgAppLog.Error(err, "error creating es creds", "name", namespacedName.Name)
			return "", "", err
		}
		return user, pass, nil
	} else if err != nil {
		cnvrgAppLog.Error(err, "can't check if es creds secret exists", "name", namespacedName.Name)
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_USER"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_USER", namespacedName.Name)
		cnvrgAppLog.Error(err, "missing required field")
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_PASS"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_PASS", namespacedName.Name)
		cnvrgAppLog.Error(err, "missing required field")
		return "", "", err
	}

	return string(creds.Data["CNVRG_ES_USER"]), string(creds.Data["CNVRG_ES_PASS"]), nil
}

func (r *CnvrgAppReconciler) pgCredsSecret(app *mlopsv1.CnvrgApp) error {

	namespacedName := types.NamespacedName{Name: app.Spec.Dbs.Pg.CredsRef, Namespace: app.Namespace}
	creds := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := r.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(app, &creds, r.Scheme); err != nil {
			cnvrgAppLog.Error(err, "error set controller reference", "name", namespacedName.Name)
			return err
		}
		user := "cnvrg"
		pass := desired.RandomString()
		database := "cnvrg_production"
		creds.Data = map[string][]byte{
			"POSTGRESQL_USER":            []byte(user),
			"POSTGRESQL_PASSWORD":        []byte(pass),
			"POSTGRESQL_ADMIN_PASSWORD":  []byte(pass),
			"POSTGRESQL_DATABASE":        []byte(database),
			"POSTGRESQL_MAX_CONNECTIONS": []byte(strconv.Itoa(app.Spec.Dbs.Pg.MaxConnections)),
			"POSTGRESQL_SHARED_BUFFERS":  []byte(app.Spec.Dbs.Pg.SharedBuffers),
			// required vars for the app
			"POSTGRES_DB":       []byte(database),
			"POSTGRES_PASSWORD": []byte(pass),
			"POSTGRES_USER":     []byte(user),
			"POSTGRES_HOST":     []byte(app.Spec.Dbs.Pg.SvcName),
		}
		if err := r.Create(context.Background(), &creds); err != nil {
			cnvrgAppLog.Error(err, "error creating pg creds", "name", namespacedName.Name)
			return err
		}
		return nil
	} else if err != nil {
		cnvrgAppLog.Error(err, "can't check if pg creds secret exists", "name", namespacedName.Name)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(cnvrgApp *mlopsv1.CnvrgApp) (bool, int, error) {

	readyState := make(map[string]bool)

	// check webapp status
	if cnvrgApp.Spec.ControlPlane.WebApp.Enabled == true {
		name := types.NamespacedName{Name: cnvrgApp.Spec.ControlPlane.WebApp.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["webApp"] = ready
	}

	// check sidekiq status
	if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled == true {
		name := types.NamespacedName{Name: "sidekiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["sidekiq"] = ready
	}

	// check searchkiq status
	if cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled == true {
		name := types.NamespacedName{Name: "searchkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["searchkiq"] = ready
	}

	// check systemkiq status
	if cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled == true {
		name := types.NamespacedName{Name: "systemkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["searchkiq"] = ready
	}

	// check postgres status
	if cnvrgApp.Spec.Dbs.Pg.Enabled == true {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Pg.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["pg"] = ready
	}

	// check minio status
	if cnvrgApp.Spec.Dbs.Minio.Enabled == true {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Minio.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["minio"] = ready
	}

	// check redis status
	if cnvrgApp.Spec.Dbs.Redis.Enabled == true {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Redis.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["redis"] = ready
	}

	// check es status
	if cnvrgApp.Spec.Dbs.Es.Enabled == true {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Es.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, err
		}
		// if es is ready, trigger fluentbit reconfiguration
		if ready {
			cnvrgAppLog.Info("es is ready, triggering fluentbit reconfiguration")
			if err := r.triggerInfraReconciler(cnvrgApp, "add"); err != nil {
				return false, 0, err
			}
		}
		readyState["es"] = ready
	}

	// check kibana status
	if cnvrgApp.Spec.Logging.Enabled == true && cnvrgApp.Spec.Logging.Kibana.Enabled == true {
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
	cnvrgInfraLog.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgApp.Namespace,
		Data:      map[string]interface{}{"Registry": cnvrgApp.Spec.Registry},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgApp, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}
	// dbs
	if err := r.dbsState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// networking
	cnvrgAppLog.Info("applying networking")
	if err := desired.Apply(networking.CnvrgAppNetworkingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// logging
	if err := r.loggingState(cnvrgApp); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// controlplane
	cnvrgAppLog.Info("applying controlplane")
	if err := desired.Apply(controlplane.State(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
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
	cnvrgAppLog.Info("applying logging")
	kibanaConfigSecretData, err := r.getKibanaConfigSecretData(app)
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	if err := desired.Apply(logging.KibanaConfSecret(*kibanaConfigSecretData), app, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}

	if err := desired.Apply(logging.CnvrgAppLoggingState(app), app, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) dbsState(app *mlopsv1.CnvrgApp) error {
	// dbs
	cnvrgAppLog.Info("applying dbs")
	// creds for es
	if _, _, err := r.esCredsSecret(app); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	// creds for pg
	if err := r.pgCredsSecret(app); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	// creds for redis
	if err := desired.CreateRedisCredsSecret(
		app,
		app.Spec.Dbs.Redis.CredsRef,
		app.Namespace,
		fmt.Sprintf("%s:%d", app.Spec.Dbs.Redis.SvcName, app.Spec.Dbs.Redis.Port),
		r,
		r.Scheme,
		cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	if err := desired.Apply(dbs.AppDbsState(app), app, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) monitoringState(app *mlopsv1.CnvrgApp) error {
	// monitoring
	cnvrgAppLog.Info("applying monitoring")
	if err := desired.CreatePromCredsSecret(app, app.Spec.Monitoring.Prometheus.CredsRef, app.Namespace, r, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	if err := r.upstreamPrometheusConfig(app); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	// grafana dashboards
	cnvrgAppLog.Info("applying grafana dashboards ")
	if err := r.createGrafanaDashboards(app); err != nil {
		return err
	}
	// grafana datasource
	cnvrgInfraLog.Info("applying grafana datasource")

	user, pass, err := desired.GetPromCredsSecret(app.Spec.Monitoring.Prometheus.CredsRef, app.Namespace, r, cnvrgAppLog)
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}
	grafanaDatasourceData := desired.TemplateData{
		Namespace: app.Namespace,
		Data: map[string]interface{}{
			"Svc":  app.Spec.Monitoring.Prometheus.SvcName,
			"Port": app.Spec.Monitoring.Prometheus.Port,
			"User": user,
			"Pass": pass,
		},
	}

	if err := desired.Apply(monitoring.GrafanaDSState(grafanaDatasourceData), app, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}

	if err := desired.Apply(monitoring.AppMonitoringState(app), app, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), app)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) getCnvrgInfra() (*mlopsv1.CnvrgInfra, error) {

	cnvrgAppInfra := &mlopsv1.CnvrgInfraList{}

	if err := r.List(context.Background(), cnvrgAppInfra); err != nil {
		cnvrgAppLog.Error(err, "can't list CnvrgInfra objects")
		return nil, err
	}

	if len(cnvrgAppInfra.Items) == 0 {
		cnvrgAppLog.Info("no CnvrgInfra objects was deployed, skipping infra reconciler")
		return nil, fmt.Errorf("no CnvrgInfra objects was deployed, skipping infra reconciler")
	}

	return &cnvrgAppInfra.Items[0], nil
}

func (r *CnvrgAppReconciler) upstreamPrometheusConfig(app *mlopsv1.CnvrgApp) error {
	namespacedName := types.NamespacedName{Name: app.Spec.Monitoring.Prometheus.UpstreamRef, Namespace: app.Namespace}
	upstreamSecret := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := r.Get(context.Background(), namespacedName, &upstreamSecret); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(app, &upstreamSecret, r.Scheme); err != nil {
			cnvrgAppLog.Error(err, "error set controller reference", "name", namespacedName.Name)
			return err
		}

		infra, err := r.getCnvrgInfra()
		if err != nil {
			cnvrgAppLog.Error(err, "can't get cnvrgInfra object", "name", namespacedName.Name)
			return err
		}

		user, pass, err := desired.GetPromCredsSecret(infra.Spec.Monitoring.Prometheus.CredsRef, infra.Spec.InfraNamespace, r.Client, cnvrgInfraLog)
		if err != nil {
			cnvrgAppLog.Error(err, "can't get cnvrgInfra prometheus creds", "name", namespacedName.Name)
			return err
		}
		upstreamPrometheus := fmt.Sprintf("prometheus-operated.%s.svc:%d", infra.Spec.InfraNamespace, infra.Spec.Monitoring.Prometheus.Port)
		promUpstreamConfig := desired.PrometheusUpstreamConfig(user, pass, app.Namespace, upstreamPrometheus)
		cnvrgAppLog.V(1).Info(promUpstreamConfig)
		upstreamSecret.Data = map[string][]byte{"prometheus-additional.yaml": []byte(promUpstreamConfig)}
		if err := r.Create(context.Background(), &upstreamSecret); err != nil {
			cnvrgAppLog.Error(err, "error crating upstream prometheus static configs")
			return err
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) getKibanaConfigSecretData(app *mlopsv1.CnvrgApp) (*desired.TemplateData, error) {
	kibanaHost := "0.0.0.0"
	kibanaPort := strconv.Itoa(app.Spec.Logging.Kibana.Port)
	esUser, esPass, err := r.esCredsSecret(app)
	if err != nil {
		cnvrgAppLog.Error(err, "can't fetch es creds")
		return nil, err
	}
	if app.Spec.SSO.Enabled == true {
		kibanaHost = "127.0.0.1"
		kibanaPort = "3000"
	}

	return &desired.TemplateData{
		Namespace: app.Namespace,
		Data: map[string]interface{}{
			"Host":   kibanaHost,
			"Port":   kibanaPort,
			"EsHost": fmt.Sprintf("http://%s.%s.svc:%d", app.Spec.Dbs.Es.SvcName, app.Namespace, app.Spec.Dbs.Es.Port),
			"EsUser": esUser,
			"EsPass": esPass,
		},
	}, nil

}

func (r *CnvrgAppReconciler) createGrafanaDashboards(cnvrgApp *mlopsv1.CnvrgApp) error {

	if cnvrgApp.Spec.Monitoring.Enabled != true {
		cnvrgInfraLog.Info("monitoring disabled, skipping grafana deployment")
		return nil
	}
	basePath := "/pkg/monitoring/tmpl/grafana/dashboards-data/"
	dashboards := desired.GrafanaAppDashboards
	// if namespace tenancy is off, deploy all (infra + app) grafana dashboards
	if cnvrgApp.Spec.NamespaceTenancy == false {
		dashboards = desired.GrafanaInfraDashboards
	}
	for _, dashboard := range dashboards {
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
				Namespace: cnvrgApp.Namespace,
			},
			Data: map[string]string{filepath.Base(f.Name()): string(b)},
		}
		if err := ctrl.SetControllerReference(cnvrgApp, cm, r.Scheme); err != nil {
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

func (r *CnvrgAppReconciler) triggerInfraReconciler(cnvrgApp *mlopsv1.CnvrgApp, op string) error {

	infra, err := r.getCnvrgInfra()

	if err != nil {
		return err
	}

	name := types.NamespacedName{
		Name:      infra.Spec.InfraReconcilerCm,
		Namespace: infra.Spec.InfraNamespace,
	}

	cm := &v1core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
	}

	esUser, esPass, err := r.esCredsSecret(cnvrgApp)
	if err != nil {
		cnvrgAppLog.Error(err, "failed to fetch es creds ")
		return err
	}

	appInstance := mlopsv1.AppInstance{SpecName: cnvrgApp.Name, SpecNs: cnvrgApp.Namespace, EsUser: esUser, EsPass: esPass}
	appInstanceBytes, err := json.Marshal(appInstance)
	if err != nil {
		cnvrgAppLog.Error(err, "failed to marshal app instance ")
		return err
	}
	if err := r.Get(context.Background(), name, cm); err != nil && errors.IsNotFound(err) {
		cnvrgAppLog.Info("infra reconciler cm does not exists, skipping", name, name)
		return nil
	} else if err != nil {
		cnvrgAppLog.Error(err, "can't get cm", "name", name)
		return err
	}

	if op == "add" {
		if cm.Data == nil {
			cm.Data = map[string]string{cnvrgApp.Namespace: string(appInstanceBytes)}
		} else {
			cm.Data[cnvrgApp.Namespace] = string(appInstanceBytes)
		}
	}
	if op == "remove" {
		delete(cm.Data, cnvrgApp.Namespace)
	}
	if err := r.Update(context.Background(), cm); err != nil {
		cnvrgAppLog.Error(err, "can't update cm", "cm", name)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgApp *mlopsv1.CnvrgApp) {
	if cnvrgApp.Status.Status == mlopsv1.StatusRemoving {
		cnvrgAppLog.Info("skipping status update, current cnvrg spec under removing status...")
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
		cnvrgAppLog.Error(err, "can't update status")
	}

}

func (r *CnvrgAppReconciler) syncCnvrgAppSpec(name types.NamespacedName) (bool, error) {

	cnvrgAppLog.Info("synchronizing cnvrgApp spec")

	// Fetch current cnvrgApp spec
	cnvrgApp, err := r.getCnvrgAppSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgApp == nil {
		return false, nil // probably cnvrgapp was removed
	}
	cnvrgAppLog = r.Log.WithValues("name", name, "ns", cnvrgApp.Namespace)

	// Get default cnvrgApp spec
	desiredSpec := mlopsv1.DefaultCnvrgAppSpec()

	// Merge current cnvrgApp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgApp.Spec, mergo.WithOverride); err != nil {
		cnvrgAppLog.Error(err, "can't merge")
		return false, err
	}

	equal := reflect.DeepEqual(desiredSpec, cnvrgApp.Spec)
	if !equal {
		cnvrgAppLog.Info("states are not equals, syncing and requeuing")
		cnvrgApp.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgApp); err != nil && errors.IsConflict(err) {
			cnvrgAppLog.Info("conflict updating cnvrgApp object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	cnvrgAppLog.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgAppReconciler) getCnvrgAppSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &cnvrgApp); err != nil {
		if errors.IsNotFound(err) {
			cnvrgAppLog.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		cnvrgAppLog.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgApp, nil
}

func (r *CnvrgAppReconciler) cleanup(cnvrgApp *mlopsv1.CnvrgApp) error {

	cnvrgAppLog.Info("running finalizer cleanup")

	// remove cnvrg-db-init
	if err := r.cleanupDbInitCm(cnvrgApp); err != nil {
		return err
	}

	// update infra reconciler cm
	if err := r.triggerInfraReconciler(cnvrgApp, "remove"); err != nil {
		return err
	}

	// cleanup pvc
	if err := r.cleanupPVCs(); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) cleanupPVCs() error {
	cnvrgAppLog.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		cnvrgAppLog.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
			if pvc.ObjectMeta.Labels["app"] == "prometheus" || pvc.ObjectMeta.Labels["app"] == "elasticsearch" {
				if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
					cnvrgInfraLog.Info("pvc already deleted")
				} else if err != nil {
					cnvrgInfraLog.Error(err, "error deleting prometheus pvc")
					return err
				}
			}
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) cleanupDbInitCm(desiredSpec *mlopsv1.CnvrgApp) error {
	cnvrgAppLog.Info("running cnvrg-db-init cleanup")
	ctx := context.Background()
	dbInitCm := &v1core.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-db-init", Namespace: desiredSpec.Namespace}}
	err := r.Delete(ctx, dbInitCm)
	if err != nil && errors.IsNotFound(err) {
		cnvrgAppLog.Info("no need to delete cnvrg-db-init, cm not found")
	} else {
		cnvrgAppLog.Error(err, "error deleting cnvrg-db-init")
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cnvrgAppLog = r.Log.WithValues("initializing", "crds")

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(e.ObjectOld) {
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgApp)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgApp)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrgapp spec wasn't changed, assuming status update, won't reconcile
				cnvrgAppLog.V(1).Info("update received", "shouldReconcileOnSpecChange", shouldReconcileOnSpecChange)

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
	cnvrgAppLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
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
