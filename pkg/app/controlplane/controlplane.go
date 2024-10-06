package controlplane

import (
	"context"
	"embed"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const fsRoot = "tmpl"

//go:embed  tmpl/*
var fs embed.FS

type CpCrdsStateManager struct {
	*desired.AssetsStateManager
}

type CpStateManager struct {
	*desired.AssetsStateManager
	app *mlopsv1.CnvrgApp
}

func NewControlPlaneStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "controlPlane")
	asm := desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot, nil)
	return &CpStateManager{AssetsStateManager: asm, app: app}
}

func (m *CpStateManager) LoadKiqs(kiqName string, hpa bool) error {
	kiqAsset := []string{fmt.Sprintf("%s.tpl", kiqName)}

	kiq := desired.NewAssetsGroup(fs, m.RootPath()+"/sidekiqs", m.Log(), &desired.LoadFilter{AssetName: kiqAsset})
	if err := kiq.LoadAssets(); err != nil {
		return err
	}

	m.AddToAssets(kiq)

	if hpa {
		hpaAsset := []string{fmt.Sprintf("%s-hpa.tpl", kiqName)}
		kiqHpa := desired.NewAssetsGroup(fs, m.RootPath()+"/sidekiqs", m.Log(), &desired.LoadFilter{AssetName: hpaAsset})
		if err := kiqHpa.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(kiqHpa)
	}
	return nil
}

func (m *CpStateManager) renderSmtpConfigs() error {
	assets := []string{fsRoot + "/conf/cm/secret-smtp.tpl"}
	f := &desired.LoadFilter{AssetName: assets}
	cfg := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
	if err := cfg.LoadAssets(); err != nil {
		return err
	}

	configData, err := m.smtpCfgData()
	if err != nil {
		return err
	}

	if err = cfg.Render(configData); err != nil {
		return err
	}

	m.AddToState(cfg)

	return nil
}

func (m *CpStateManager) smtpCfgData() (map[string]interface{}, error) {
	var userName, password string

	if m.app.Spec.ControlPlane.SMTP.CredentialsSecretRef != "" {
		secret := &corev1.Secret{}
		if err := m.C.Get(context.Background(), types.NamespacedName{Name: m.app.Spec.ControlPlane.SMTP.CredentialsSecretRef, Namespace: m.app.Namespace}, secret); err != nil {
			return nil, err
		}
		userName = string(secret.Data["username"])
		password = string(secret.Data["password"])
	} else {
		userName = m.app.Spec.ControlPlane.SMTP.Username
		password = m.app.Spec.ControlPlane.SMTP.Password
	}

	d := map[string]interface{}{
		"Namespace":         m.app.Namespace,
		"Annotations":       m.app.Spec.Annotations,
		"Server":            m.app.Spec.ControlPlane.SMTP.Server,
		"Port":              m.app.Spec.ControlPlane.SMTP.Port,
		"Username":          userName,
		"Password":          password,
		"Domain":            m.app.Spec.ControlPlane.SMTP.Domain,
		"OpenSSLVerifyMode": m.app.Spec.ControlPlane.SMTP.OpensslVerifyMode,
		"Sender":            m.app.Spec.ControlPlane.SMTP.Sender,
	}

	return d, nil
}

func (m *CpStateManager) Load() error {
	f := &desired.LoadFilter{DefaultLoader: true}

	if err := m.renderSmtpConfigs(); err != nil {
		return err
	}

	conf := desired.NewAssetsGroup(fs, fsRoot+"/conf/cm", m.Log(), f)
	if err := conf.LoadAssets(); err != nil {
		return err
	}
	m.AddToAssets(conf)

	rbac := desired.NewAssetsGroup(fs, fsRoot+"/conf/rbac", m.Log(), f)
	if err := rbac.LoadAssets(); err != nil {
		return err
	}
	m.AddToAssets(rbac)

	if m.app.Spec.ControlPlane.Hyper.Enabled {
		hyper := desired.NewAssetsGroup(fs, m.RootPath()+"/hyper", m.Log(), nil)
		if err := hyper.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(hyper)
	}

	if m.app.Spec.ControlPlane.CnvrgScheduler.Enabled {
		scheduler := desired.NewAssetsGroup(fs, m.RootPath()+"/scheduler", m.Log(), nil)
		if err := scheduler.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(scheduler)
	}

	if m.app.Spec.ControlPlane.Sidekiq.Enabled {
		if err := m.LoadKiqs("sidekiq", m.app.Spec.ControlPlane.Sidekiq.Hpa.Enabled); err != nil {
			return err
		}
	}

	if m.app.Spec.ControlPlane.Searchkiq.Enabled {
		if err := m.LoadKiqs("searchkiq", m.app.Spec.ControlPlane.Searchkiq.Hpa.Enabled); err != nil {
			return err
		}
	}

	if m.app.Spec.ControlPlane.Systemkiq.Enabled {
		if err := m.LoadKiqs("systemkiq", m.app.Spec.ControlPlane.Systemkiq.Hpa.Enabled); err != nil {
			return err
		}
	}

	if m.app.Spec.ControlPlane.WebApp.Enabled {
		webapp := desired.NewAssetsGroup(fs, m.RootPath()+"/webapp", m.Log(),
			&desired.LoadFilter{Ingress: &m.app.Spec.Networking.Ingress.Type, DefaultLoader: true})
		if err := webapp.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(webapp)
	}

	if m.app.Spec.ControlPlane.Nomex.Enabled {

		nomex := desired.NewAssetsGroup(fs, m.RootPath()+"/nomex", m.Log(), &desired.LoadFilter{DefaultLoader: true})
		if err := nomex.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(nomex)
	}

	return nil
}

func (m *CpStateManager) Apply() error {
	if err := m.Load(); err != nil {
		return err
	}

	if err := m.AssetsStateManager.Render(); err != nil {
		return err
	}
	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}
	return nil
}
