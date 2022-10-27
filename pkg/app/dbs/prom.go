package dbs

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"golang.org/x/crypto/bcrypt"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type PromStateManager struct {
	*desired.AssetsStateManager
	app            *mlopsv1.CnvrgApp
	promCreds      *desired.AssetsGroup
	promWebConfigs *desired.AssetsGroup
}

type GrafanaStateManager struct {
	*desired.AssetsStateManager
	app        *mlopsv1.CnvrgApp
	dashboards []*grafanaDashboard
}

type grafanaDashboard struct {
	dashboardName string
	fileName      string
	content       string
}

func newGrafanaDashboards() (dashboards []*grafanaDashboard, err error) {
	dirPath := fsRoot + "/prom/grafana/dashboards"
	dirEntries, err := fs.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}
		f, err := fs.ReadFile(dirPath + "/" + e.Name())
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, &grafanaDashboard{
			dashboardName: strings.TrimSuffix(filepath.Base(e.Name()), filepath.Ext(e.Name())),
			fileName:      filepath.Base(e.Name()),
			content:       string(f),
		})
	}
	return
}

func NewPromStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "prom")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	return &PromStateManager{
		app:                app,
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/prom/prometheus", f),
	}
}

func NewGrafanaStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "grafana")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	dashboards, err := newGrafanaDashboards()
	if err != nil {
		l.Error(err, "failed to load grafana dashboards")
	}
	return &GrafanaStateManager{
		app:                app,
		dashboards:         dashboards,
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/prom/grafana", f),
	}
}

func (m *PromStateManager) promCredsAssetsToState() error {
	assets := []string{"creds.tpl", "webconfigs.tpl"}

	credsAssets := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: assets})

	if err := credsAssets.LoadAssets(); err != nil {
		return err
	}

	data, err := m.promCredsData()
	if err != nil {
		return err
	}

	if err := credsAssets.Render(data); err != nil {
		return err
	}

	m.AddToState(credsAssets)

	return nil
}

func (m *PromStateManager) defaultRoleBinding() error {
	assets := []string{"rolebinding.tpl"}

	roleBinding := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: assets})
	if err := roleBinding.LoadAssets(); err != nil {
		return nil
	}
	data := map[string]interface{}{
		"Namespace": m.app.Namespace,
		"Spec": map[string]interface{}{
			"Annotations": m.app.Spec.Annotations,
			"Labels":      m.app.Spec.Labels,
		},
		"CnvrgNamespace": m.app.Namespace,
	}

	if err := roleBinding.Render(data); err != nil {
		return err
	}

	m.AddToState(roleBinding)

	return nil
}

func (m *PromStateManager) extraPodsScrapeConfigs() {
	if len(m.app.Spec.Dbs.Prom.ExtraScrapeConfigs) > 0 {

		for _, podScrapeConfig := range m.app.Spec.Dbs.Prom.ExtraScrapeConfigs {
			if podScrapeConfig.Namespace != m.app.Namespace {
				assets := []string{"role.tpl", "rolebinding.tpl"}
				extraRbac := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: assets})
				if err := extraRbac.LoadAssets(); err != nil {
					m.Log().Error(err, "unable to setup RBAC for extra pods scrape configs")
					continue
				}
				rbacData := map[string]interface{}{
					"Namespace": podScrapeConfig.Namespace,
					"Spec": map[string]interface{}{
						"Annotations": map[string]string{},
						"Labels":      map[string]string{},
					},
					"CnvrgNamespace": m.app.Namespace,
				}
				if err := extraRbac.Render(rbacData); err != nil {
					m.Log().Error(err, "failed to render RBAC for extra pods scrape configs")
					continue
				}
				m.AddToState(extraRbac)
			}
		}
	}
}

func (m *PromStateManager) promCredsData() (map[string]interface{}, error) {
	user := "cnvrg"
	pass := desired.RandomString()
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"Namespace":   m.app.Namespace,
		"Annotations": m.app.Spec.Annotations,
		"Labels":      m.app.Spec.Labels,
		"CredsRef":    m.app.Spec.Dbs.Prom.CredsRef,
		"User":        user,
		"Pass":        pass,
		"PassHash":    string(passHash),
		"PromUrl":     fmt.Sprintf("http://%s.%s.svc:%d", "prom", m.app.Namespace, 9090),
	}, nil
}

func (m *PromStateManager) Apply() error {

	// create default role binding for prometheus
	if err := m.defaultRoleBinding(); err != nil {
		return err
	}
	// parse & render prom creds
	if err := m.promCredsAssetsToState(); err != nil {
		return err
	}
	// try to apply extra pods scrape configs
	m.extraPodsScrapeConfigs()
	// apply the final state
	if err := m.AssetsStateManager.Apply(); err != nil {
		return nil
	}
	return nil
}

func (m *GrafanaStateManager) promCreds() (url, user, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: m.app.Spec.Dbs.Prom.CredsRef, Namespace: m.app.Namespace}
	creds := v1core.Secret{ObjectMeta: v1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := m.C.Get(context.Background(), namespacedName, &creds); err != nil {
		return "", "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_USER"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_USER", namespacedName.Name)
		return "", "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_PASS"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_PASS", namespacedName.Name)
		return "", "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_URL"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_URL", namespacedName.Name)
		return "", "", "", err
	}

	return string(creds.Data["CNVRG_PROMETHEUS_URL"]), string(creds.Data["CNVRG_PROMETHEUS_USER"]), string(creds.Data["CNVRG_PROMETHEUS_PASS"]), nil
}

func (m *GrafanaStateManager) dataSources() error {
	assets := []string{"datasource.tpl"}
	ds := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: assets})
	if err := ds.LoadAssets(); err != nil {
		return err
	}

	url, user, pass, err := m.promCreds()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Namespace":   m.app.Namespace,
		"Annotations": m.app.Spec.Annotations,
		"Labels":      m.app.Spec.Labels,
		"Url":         url,
		"User":        user,
		"Pass":        pass,
	}

	if err := ds.Render(data); err != nil {
		return err
	}

	m.AddToState(ds)

	return nil

}

func (m *GrafanaStateManager) createDashboards() error {

	for _, d := range m.dashboards {
		cm := &v1core.ConfigMap{
			ObjectMeta: v1.ObjectMeta{
				Name:      d.dashboardName,
				Namespace: m.app.Namespace,
			},
			Data: map[string]string{d.fileName: d.content},
		}
		if err := m.C.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
			continue
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (m *GrafanaStateManager) deployment() error {

	dep := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: []string{"dep.tpl"}})
	if err := dep.LoadAssets(); err != nil {
		return err
	}

	data := map[string]interface{}{
		"Namespace":  m.app.Namespace,
		"Spec":       m.app.Spec,
		"Dashboards": m.dashboardsToList(),
	}
	if err := dep.Render(data); err != nil {
		return err
	}

	m.AddToState(dep)

	return nil
}

func (m *GrafanaStateManager) dashboardsToList() (dashboards []string) {
	for _, d := range m.dashboards {
		dashboards = append(dashboards, d.dashboardName)
	}
	return
}

func (m *GrafanaStateManager) oauthProxy() error {

	cm := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: []string{"oauth.tpl"}})

	if err := cm.LoadAssets(); err != nil {
		return err
	}

	if err := cm.Render(m.app); err != nil {
		return err
	}

	m.AddToState(cm)

	return nil

}

func (m *GrafanaStateManager) Apply() error {

	// create grafana dashboards as config maps
	if err := m.createDashboards(); err != nil {
		return err
	}
	// create prometheus data source configmap for grafana
	if err := m.dataSources(); err != nil {
		return err
	}
	// create grafana deployment
	if err := m.deployment(); err != nil {
		return err
	}
	// create oauth2-proxy configs if sso enabled
	if m.app.Spec.SSO.Enabled {
		if err := m.oauthProxy(); err != nil {
			return err
		}
	}
	// apply all the default assets
	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}
