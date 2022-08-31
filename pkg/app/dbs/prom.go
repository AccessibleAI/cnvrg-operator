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
	app       *mlopsv1.CnvrgApp
	promCreds *desired.AssetsGroup
	promCm    *desired.AssetsGroup
}

func NewPromStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "prom")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	return &PromStateManager{
		app:                app,
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/prom", f),
	}
}

func (m *PromStateManager) Load() error {
	credsAsset := "creds.tpl"
	m.promCreds = desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: &credsAsset})
	if err := m.promCreds.LoadAssets(); err != nil {
		return err
	}
	configMapAsset := "cm.tpl"
	m.promCm = desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), &desired.LoadFilter{AssetName: &configMapAsset})
	if err := m.promCm.LoadAssets(); err != nil {
		return err
	}
	return nil
}

func (m *PromStateManager) Render() error {
	data, err := m.promCredsData()
	if err != nil {
		return err
	}

	if err := m.promCreds.Render(data); err != nil {
		return err
	}
	m.AssetsStateManager.AddToState(m.promCreds)

	if err := m.promCm.Render(data); err != nil {
		return err
	}
	m.AssetsStateManager.AddToState(m.promCm)

	return nil
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

	if err := m.Load(); err != nil {
		return err
	}

	if err := m.Render(); err != nil {
		return err
	}

	if err := m.AssetsStateManager.Apply(); err != nil {
		return nil
	}

	return nil
}
