package app

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
)

func (r *CnvrgAppReconciler) getEsCredsSecret(app *mlopsv1.CnvrgApp) (user string, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: app.Spec.Dbs.Es.CredsRef, Namespace: app.Namespace}
	creds := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := r.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		r.Log.Error(err, "es-creds secret not found!")
		return "", "", err
	} else if err != nil {
		r.Log.Error(err, "can't check if es creds secret exists", "name", namespacedName.Name)
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_USER"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_USER", namespacedName.Name)
		r.Log.Error(err, "missing required field")
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_PASS"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_PASS", namespacedName.Name)
		r.Log.Error(err, "missing required field")
		return "", "", err
	}

	return string(creds.Data["CNVRG_ES_USER"]), string(creds.Data["CNVRG_ES_PASS"]), nil
}

func (r *CnvrgAppReconciler) dbsState(app *mlopsv1.CnvrgApp) error {

	//// dbs
	//r.Log.Info("applying dbs")
	//
	//if app.Spec.Dbs.Es.Enabled {
	//	esSecretData := desired.TemplateData{
	//		Data: map[string]interface{}{
	//			"Namespace":   app.Namespace,
	//			"CredsRef":    app.Spec.Dbs.Es.CredsRef,
	//			"EsUrl":       fmt.Sprintf("%s.%s.svc:%d", app.Spec.Dbs.Es.SvcName, app.Namespace, app.Spec.Dbs.Es.Port),
	//			"Annotations": app.Spec.Annotations,
	//			"Labels":      app.Spec.Labels,
	//		},
	//	}
	//	r.Log.Info("trying to generate es creds (if still doesn't exists...)")
	//	if err := desired.Apply(dbs.EsCreds(esSecretData), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//
	//	kibanaConfigSecretData, err := r.getKibanaConfigSecretData(app)
	//	if err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//	if err := desired.Apply(dbs.KibanaConfSecret(*kibanaConfigSecretData), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//	if err := desired.Apply(dbs.KibanaState(app), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//
	//}
	//
	//if app.Spec.Dbs.Pg.Enabled {
	//	pgSecretData := desired.TemplateData{
	//		Data: map[string]interface{}{
	//			"Namespace":          app.Namespace,
	//			"CredsRef":           app.Spec.Dbs.Pg.CredsRef,
	//			"Annotations":        app.Spec.Annotations,
	//			"Labels":             app.Spec.Labels,
	//			"MaxConnections":     app.Spec.Dbs.Pg.MaxConnections,
	//			"SharedBuffers":      app.Spec.Dbs.Pg.SharedBuffers,
	//			"EffectiveCacheSize": app.Spec.Dbs.Pg.EffectiveCacheSize,
	//			"SvcName":            app.Spec.Dbs.Pg.SvcName,
	//		},
	//	}
	//	r.Log.Info("trying to generate pg creds (if still doesn't exists...)")
	//	if err := desired.Apply(dbs.PgCreds(pgSecretData), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//}
	//
	//if app.Spec.Dbs.Redis.Enabled {
	//	redisSecretData := desired.TemplateData{
	//		Data: map[string]interface{}{
	//			"Namespace":   app.Namespace,
	//			"Annotations": app.Spec.Annotations,
	//			"Labels":      app.Spec.Labels,
	//			"CredsRef":    app.Spec.Dbs.Redis.CredsRef,
	//			"SvcName":     app.Spec.Dbs.Redis.SvcName,
	//		},
	//	}
	//	r.Log.Info("trying to generate redis creds (if still doesn't exists...)")
	//	if err := desired.Apply(dbs.RedisCreds(redisSecretData), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//}
	//
	//if app.Spec.Dbs.Prom.Enabled {
	//	if err := r.createPromDBCreds(app); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//
	//	if _, _, passHash, err := desired.GetPromCredsSecret("prom-creds", app.Namespace, r.Client, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	} else {
	//		promData := desired.TemplateData{
	//			Data: map[string]interface{}{
	//				"Namespace":     app.Namespace,
	//				"Annotations":   app.Spec.Annotations,
	//				"Labels":        app.Spec.Labels,
	//				"PassHash":      passHash,
	//				"ClusterDomain": app.Spec.ClusterDomain,
	//				"HttpsEnabled":  app.Spec.Networking.HTTPS.Enabled,
	//				"RegistryName":  app.Spec.Registry.Name,
	//			},
	//		}
	//		if err := desired.Apply(dbs.ApplyAppPrometheus(app, promData), app, r.Client, r.Scheme, r.Log); err != nil {
	//			r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//			return err
	//		}
	//	}
	//}
	//
	//if app.Spec.Dbs.Elastalert.Enabled {
	//	// create elastalert creds ref
	//	data, err := generateElastalertCreds(app)
	//	if err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//	if err := desired.Apply(dbs.ElastCreds(data), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//
	//	if err := desired.Apply(dbs.ElastAlert(), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//}
	//
	//if err := desired.Apply(dbs.AppDbsState(app), app, r.Client, r.Scheme, r.Log); err != nil {
	//	r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//	return err
	//}
	return nil
}

func (r *CnvrgAppReconciler) createPromDBCreds(app *mlopsv1.CnvrgApp) error {
	//if app.Spec.Dbs.Prom.Enabled {
	//	user := "cnvrg"
	//	pass := desired.RandomString()
	//	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	//	if err != nil {
	//		r.Log.Error(err, "error generating prometheus hash")
	//		return err
	//	}
	//	promSecretData := desired.TemplateData{
	//		Data: map[string]interface{}{
	//			"Namespace":   app.Namespace,
	//			"Annotations": app.Spec.Annotations,
	//			"Labels":      app.Spec.Labels,
	//			"CredsRef":    "prom-creds",
	//			"User":        user,
	//			"Pass":        pass,
	//			"PassHash":    string(passHash),
	//			"PromUrl":     fmt.Sprintf("http://%s.%s.svc:%d", "prom", app.Namespace, 9090),
	//		},
	//	}
	//	if err := desired.Apply(dbs.PromDBCreds(promSecretData), app, r.Client, r.Scheme, r.Log); err != nil {
	//		r.updateStatusMessage(mlopsv1.Status{Status: mlopsv1.StatusError, Message: err.Error(), Progress: -1}, app)
	//		return err
	//	}
	//}
	return nil
}

func generateElastalertCreds(app *mlopsv1.CnvrgApp) (*desired.TemplateData, error) {
	//user := "cnvrg"
	//pass := desired.RandomString()
	//passwordHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
	//if err != nil {
	//	//r.Log.Error(err, "error generating elastalert hash")
	//	return nil, err
	//}
	//
	//httpSchema := "http://"
	//if app.Spec.Networking.HTTPS.Enabled {
	//	httpSchema = "https://"
	//}
	//
	//data := &desired.TemplateData{
	//	Data: map[string]interface{}{
	//		"Namespace":     app.Namespace,
	//		"Annotations":   app.Spec.Annotations,
	//		"Labels":        app.Spec.Labels,
	//		"CredsRef":      app.Spec.Dbs.Elastalert.CredsRef,
	//		"User":          user,
	//		"Pass":          pass,
	//		"Htpasswd":      fmt.Sprintf("%s:%s", user, passwordHash),
	//		"ElastAlertUrl": fmt.Sprintf("%s%s.%s", httpSchema, app.Spec.Dbs.Elastalert.SvcName, app.Spec.ClusterDomain),
	//	},
	//}
	//
	//return data, nil
	return nil, nil
}

func (r *CnvrgAppReconciler) getKibanaConfigSecretData(app *mlopsv1.CnvrgApp) (*desired.TemplateData, error) {
	kibanaHost := "0.0.0.0"
	kibanaPort := strconv.Itoa(app.Spec.Dbs.Es.Kibana.Port)
	esUser, esPass, err := r.getEsCredsSecret(app)
	if err != nil {
		r.Log.Error(err, "can't fetch es creds")
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
