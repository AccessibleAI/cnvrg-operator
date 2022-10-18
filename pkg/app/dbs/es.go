package dbs

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/Dimss/crypt/apr1_crypt"
	"github.com/go-logr/logr"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

type ElasticStateManager struct {
	*desired.AssetsStateManager
	app      *mlopsv1.CnvrgApp
	esSecret *desired.AssetsGroup
}

type KibanaStateManager struct {
	*desired.AssetsStateManager
	app        *mlopsv1.CnvrgApp
	kibanaConf *desired.AssetsGroup
}

type ElastAlertStateManager struct {
	*desired.AssetsStateManager
	app    *mlopsv1.CnvrgApp
	eaConf *desired.AssetsGroup
}

func NewElasticStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "elastic")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	asm := desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/es/elastic", f)
	return &ElasticStateManager{AssetsStateManager: asm, app: app}
}

func NewKibanaStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "kibana")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	asm := desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/es/kibana", f)
	return &KibanaStateManager{AssetsStateManager: asm, app: app}
}

func NewElastAlertStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "elastAlert")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	asm := desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/es/elastalert", f)
	return &ElastAlertStateManager{AssetsStateManager: asm, app: app}
}

func (m *ElasticStateManager) Load() error {

	assetName := []string{"secret.tpl"}
	f := &desired.LoadFilter{AssetName: assetName}
	m.esSecret = desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
	if err := m.esSecret.LoadAssets(); err != nil {
		return err
	}

	return nil
}

func (m *ElasticStateManager) RenderDeployment() error {

	esSecretData := map[string]interface{}{
		"Namespace":   m.app.Namespace,
		"CredsRef":    m.app.Spec.Dbs.Es.CredsRef,
		"EsUrl":       fmt.Sprintf("%s.%s.svc:%d", m.app.Spec.Dbs.Es.SvcName, m.app.Namespace, m.app.Spec.Dbs.Es.Port),
		"Annotations": m.app.Spec.Annotations,
		"Labels":      m.app.Spec.Labels,
	}

	if err := m.esSecret.Render(esSecretData); err != nil {
		return err
	}
	m.AssetsStateManager.AddToState(m.esSecret)

	return nil

}

func (m *ElasticStateManager) Apply() error {
	if err := m.Load(); err != nil {
		return err
	}

	if err := m.RenderDeployment(); err != nil {
		return err
	}

	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}
	return nil
}

func (m *KibanaStateManager) Load() error {
	assetName := []string{"secret.tpl"}
	f := &desired.LoadFilter{AssetName: assetName}
	m.kibanaConf = desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
	if err := m.kibanaConf.LoadAssets(); err != nil {
		return err
	}
	return nil

}

func (m *KibanaStateManager) RenderDeployment() error {
	if data, err := m.getKibanaConfigSecretData(); err == nil {
		if err := m.kibanaConf.Render(data); err != nil {
			return err
		}
		m.AssetsStateManager.AddToState(m.kibanaConf)
	} else {
		return err
	}
	return nil
}

func (m *KibanaStateManager) getKibanaConfigSecretData() (map[string]interface{}, error) {
	kibanaHost := "0.0.0.0"
	kibanaPort := strconv.Itoa(m.app.Spec.Dbs.Es.Kibana.Port)
	esUser, esPass, err := m.getEsCredsSecret()
	if err != nil {
		m.Log().Error(err, "can't fetch es creds")
		return nil, err
	}
	if m.app.Spec.SSO.Enabled {
		kibanaHost = "127.0.0.1"
		kibanaPort = "3000"
	}
	return map[string]interface{}{
		"Namespace":   m.app.Namespace,
		"Host":        kibanaHost,
		"Port":        kibanaPort,
		"EsHost":      fmt.Sprintf("http://%s.%s.svc:%d", m.app.Spec.Dbs.Es.SvcName, m.app.Namespace, m.app.Spec.Dbs.Es.Port),
		"EsUser":      esUser,
		"EsPass":      esPass,
		"Annotations": m.app.Spec.Annotations,
		"Labels":      m.app.Spec.Labels,
	}, nil
}

func (m *KibanaStateManager) getEsCredsSecret() (user string, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: m.app.Spec.Dbs.Es.CredsRef, Namespace: m.app.Namespace}
	creds := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := m.AssetsStateManager.C.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		return "", "", err
	} else if err != nil {
		m.Log().Error(err, "can't check if es creds secret exists", "name", namespacedName.Name)
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_USER"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_USER", namespacedName.Name)
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_ES_PASS"]; !ok {
		err := fmt.Errorf("es creds secret %s missing require field CNVRG_ES_PASS", namespacedName.Name)
		return "", "", err
	}

	return string(creds.Data["CNVRG_ES_USER"]), string(creds.Data["CNVRG_ES_PASS"]), nil
}

func (m *KibanaStateManager) Apply() error {
	if err := m.Load(); err != nil {
		return nil
	}

	if err := m.RenderDeployment(); err != nil {
		return nil
	}

	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}
	return nil
}

func (m *ElastAlertStateManager) Load() error {
	assetName := []string{"credsec.tpl"}
	f := &desired.LoadFilter{AssetName: assetName}
	m.eaConf = desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
	if err := m.eaConf.LoadAssets(); err != nil {
		return err
	}
	return nil
}

func (m *ElastAlertStateManager) RenderDeployment() error {
	eaConfig, err := m.getElastAlertConfigs()
	if err != nil {
		return err
	}
	if err := m.eaConf.Render(eaConfig); err != nil {
		return err
	}
	m.AssetsStateManager.AddToState(m.eaConf)
	return nil
}

func (m *ElastAlertStateManager) getElastAlertConfigs() (map[string]interface{}, error) {
	user := "cnvrg"
	pass := desired.RandomString()
	passwordHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
	if err != nil {
		return nil, err
	}

	httpSchema := "http://"
	if m.app.Spec.Networking.HTTPS.Enabled {
		httpSchema = "https://"
	}

	return map[string]interface{}{
		"Namespace":     m.app.Namespace,
		"Annotations":   m.app.Spec.Annotations,
		"Labels":        m.app.Spec.Labels,
		"CredsRef":      m.app.Spec.Dbs.Es.Elastalert.CredsRef,
		"User":          user,
		"Pass":          pass,
		"Htpasswd":      fmt.Sprintf("%s:%s", user, passwordHash),
		"ElastAlertUrl": fmt.Sprintf("%s%s.%s", httpSchema, m.app.Spec.Dbs.Es.Elastalert.SvcName, m.app.Spec.ClusterDomain),
	}, nil

}

func (m *ElastAlertStateManager) Apply() error {

	if err := m.Load(); err != nil {
		return err
	}

	if err := m.RenderDeployment(); err != nil {
		return err
	}

	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}
