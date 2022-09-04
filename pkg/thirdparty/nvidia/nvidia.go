package nvidia

import (
	"context"
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	if m.isOpenshift() {
		assetName := "scc.tpl"
		f := &desired.LoadFilter{AssetName: &assetName}
		scc := desired.NewAssetsGroup(fs, m.RootPath(), m.Log(), f)
		if err := scc.LoadAssets(); err != nil {
			return err
		}
		m.AddToAssets(scc)
	}
	return nil
}

func (m *RbacState) isOpenshift() bool {
	routes := &unstructured.UnstructuredList{}
	routes.SetGroupVersionKind(desired.Kinds["OcpRouteGVK"])
	if err := m.C.List(context.Background(), routes); err != nil {
		return false
	}
	return true
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
