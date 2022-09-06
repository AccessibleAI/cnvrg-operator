package controlplane

import (
	"embed"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
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

func NewControlPlaneCredsStateManager(c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "controlPlaneCrds")
	asm := desired.NewAssetsStateManager(nil, c, s, l, fs, fsRoot+"/crds", nil)
	return &CpCrdsStateManager{AssetsStateManager: asm}
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

func (m *CpStateManager) Load() error {
	f := &desired.LoadFilter{DefaultLoader: true}

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

	if m.app.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
		assetName := []string{"ocp-scc.tpl"}
		ocpScc := desired.NewAssetsGroup(fs, fsRoot+"/conf/rbac", m.Log(), &desired.LoadFilter{AssetName: assetName})
		if err := ocpScc.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(ocpScc)
	}

	if m.app.Spec.ControlPlane.Hyper.Enabled {
		hyper := desired.NewAssetsGroup(fs, m.RootPath()+"/hyper", m.Log(), nil)
		if err := hyper.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(hyper)
	}

	if m.app.Spec.ControlPlane.Mpi.Enabled {
		mpi := desired.NewAssetsGroup(fs, m.RootPath()+"/mpi", m.Log(), nil)
		if err := mpi.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(mpi)
	}

	if m.app.Spec.ControlPlane.CnvrgRouter.Enabled {
		router := desired.NewAssetsGroup(fs, m.RootPath()+"/router", m.Log(), nil)
		if err := router.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(router)
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
