package nvidia

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
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

func NewNvidiaStateManager(ctp *mlopsv1.CnvrgThirdParty, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "nvidia")
	f := &desired.LoadFilter{DefaultLoader: true}
	asm := desired.NewAssetsStateManager(ctp, c, s, l, fs, fsRoot+"/deviceplugin", f)
	return &DevicePluginState{AssetsStateManager: asm}
}
