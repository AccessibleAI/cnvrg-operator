package dbs

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PgStateManager struct {
	*desired.AssetsStateManager
	app      *mlopsv1.CnvrgApp
	pgSecret *desired.AssetsGroup
}

func NewPgStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "pg")
	f := &desired.LoadFilter{DefaultLoader: true}
	asm := desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/pg", f)
	return &PgStateManager{AssetsStateManager: asm, app: app}
}

func (m *PgStateManager) Load() error {
	assetName := "secret.tpl"
	f := &desired.LoadFilter{AssetName: &assetName}
	m.pgSecret = desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
	if err := m.pgSecret.LoadAssets(); err != nil {
		return err
	}
	return nil
}

func (m *PgStateManager) Render() error {
	pgSecretData := map[string]interface{}{
		"Namespace":          m.app.Namespace,
		"CredsRef":           m.app.Spec.Dbs.Pg.CredsRef,
		"Annotations":        m.app.Spec.Annotations,
		"Labels":             m.app.Spec.Labels,
		"MaxConnections":     m.app.Spec.Dbs.Pg.MaxConnections,
		"SharedBuffers":      m.app.Spec.Dbs.Pg.SharedBuffers,
		"EffectiveCacheSize": m.app.Spec.Dbs.Pg.EffectiveCacheSize,
		"SvcName":            m.app.Spec.Dbs.Pg.SvcName,
	}
	if err := m.pgSecret.Render(pgSecretData); err != nil {
		return err
	}
	m.AssetsStateManager.AddToState(m.pgSecret)
	return nil
}

func (m *PgStateManager) Apply() error {

	if err := m.Load(); err != nil {
		return err
	}

	if err := m.Render(); err != nil {
		return err
	}

	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}
