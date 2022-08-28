package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/monitoring"
	"github.com/Dimss/crypt/apr1_crypt"
	"github.com/markbates/pkger"
	"io/ioutil"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

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
		appLog.Info("applying grafana dashboards")
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
