package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/logging"
	"github.com/Dimss/crypt/apr1_crypt"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
)

func (r *CnvrgAppReconciler) loggingState(app *mlopsv1.CnvrgApp) error {

	if app.Spec.Logging.Kibana.Enabled {
		appLog.Info("applying kibana")
		kibanaConfigSecretData, err := r.getKibanaConfigSecretData(app)
		if err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
		if err := desired.Apply(logging.KibanaConfSecret(*kibanaConfigSecretData), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
		if err := desired.Apply(logging.CnvrgAppKibanaState(app), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
	}

	if app.Spec.Logging.Elastalert.Enabled {
		appLog.Info("applying elastalert")

		// create elastalert creds ref
		data, err := generateElastalertCreds(app)
		if err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}
		if err := desired.Apply(logging.ElastCreds(data), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}

		if err := desired.Apply(logging.ElastAlert(), app, r.Client, r.Scheme, appLog); err != nil {
			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
			return err
		}

	}

	return nil
}

func generateElastalertCreds(app *mlopsv1.CnvrgApp) (*desired.TemplateData, error) {
	user := "cnvrg"
	pass := desired.RandomString()
	passwordHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
	if err != nil {
		appLog.Error(err, "error generating elastalert hash")
		return nil, err
	}

	httpSchema := "http://"
	if app.Spec.Networking.HTTPS.Enabled {
		httpSchema = "https://"
	}

	data := &desired.TemplateData{
		Data: map[string]interface{}{
			"Namespace":     app.Namespace,
			"Annotations":   app.Spec.Annotations,
			"Labels":        app.Spec.Labels,
			"CredsRef":      app.Spec.Logging.Elastalert.CredsRef,
			"User":          user,
			"Pass":          pass,
			"Htpasswd":      fmt.Sprintf("%s:%s", user, passwordHash),
			"ElastAlertUrl": fmt.Sprintf("%s%s.%s", httpSchema, app.Spec.Logging.Elastalert.SvcName, app.Spec.ClusterDomain),
		},
	}

	return data, nil
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

func (r *CnvrgAppReconciler) addFluentbitConfiguration(cnvrgApp *mlopsv1.CnvrgApp) error {
	infra, err := r.getCnvrgInfra()
	if errors.IsNotFound(err) {
		appLog.Info("cnvrginfra not found, skipping fluentbit configuration")
		return nil
	}
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
