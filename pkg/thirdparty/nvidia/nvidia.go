package nvidia

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const fsRoot = "tmpl"

//go:embed  tmpl/*
var fs embed.FS

type DevicePluginState struct {
	*desired.AssetsStateManager
}

type MetricsExporterState struct {
	*desired.AssetsStateManager
}

type RbacState struct {
	*desired.AssetsStateManager
	ctp *mlopsv1.CnvrgThirdParty
}

func NewNvidiaStateManager(ctp *mlopsv1.CnvrgThirdParty, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "nvidiaDevicePlugin")
	f := &desired.LoadFilter{DefaultLoader: true}
	asm := desired.NewAssetsStateManager(ctp, c, s, l, fs, fsRoot+"/deviceplugin", f)
	return &DevicePluginState{AssetsStateManager: asm}
}

func NewMetricsExporterState(ctp *mlopsv1.CnvrgThirdParty, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "nvidiaMetricsExporter")
	f := &desired.LoadFilter{DefaultLoader: true}
	asm := desired.NewAssetsStateManager(ctp, c, s, l, fs, fsRoot+"/dcgm", f)
	return &MetricsExporterState{AssetsStateManager: asm}
}

func NewNvidiaRbacState(ctp *mlopsv1.CnvrgThirdParty, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "nvidiaRbac")
	f := &desired.LoadFilter{DefaultLoader: true}
	asm := desired.NewAssetsStateManager(ctp, c, s, l, fs, fsRoot+"/rbac", f)
	return &RbacState{AssetsStateManager: asm, ctp: ctp}
}

func (m *RbacState) Load() error {
	if utils.IsOpenShift(m.C) {
		f := &desired.LoadFilter{AssetName: []string{"scc.tpl"}}
		scc := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
		if err := scc.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(scc)
	}
	return nil
}

func (m *RbacState) Apply() error {
	if err := m.Load(); err != nil {
		return err
	}

	if err := m.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}