package metagpu

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed  tmpl/*
var fs embed.FS

const fsRoot = "tmpl"

type DevicePluginState struct {
	*desired.AssetsStateManager
}

type DevicePluginSccState struct {
	*desired.AssetsStateManager
}

func NewDevicePluginState(ctp *mlopsv1.CnvrgThirdParty, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "metaGpuDevicePlugin")
	f := &desired.LoadFilter{DefaultLoader: true}
	asm := desired.NewAssetsStateManager(ctp, c, s, l, fs, fsRoot, f)
	dps := &DevicePluginState{AssetsStateManager: asm}
	if utils.IsOpenShift(c) {
		f := &desired.LoadFilter{AssetName: []string{"scc.tpl"}}
		scc := desired.NewAssetsGroup(fs, fsRoot, dps.Log(), f)
		if err := scc.LoadAssets(); err != nil {
			dps.Log().Error(err, "error loading scc for metagpu")
		} else {
			dps.AddToAssets(scc)
		}
	}
	return dps
}
