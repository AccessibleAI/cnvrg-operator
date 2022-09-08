package dbs

import (
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PromStateManager struct {
	*desired.AssetsStateManager
	app            *mlopsv1.CnvrgApp
	promCreds      *desired.AssetsGroup
	promWebConfigs *desired.AssetsGroup
}

func NewPromStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "prom")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	return &PromStateManager{
		app:                app,
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/prom", f),
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
	if len(m.app.Spec.Dbs.Prom.ExtraPodsScrapeConfigs) > 0 {

		for _, podScrapeConfig := range m.app.Spec.Dbs.Prom.ExtraPodsScrapeConfigs {
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
