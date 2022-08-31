package registry

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

type State struct {
	*desired.AssetsStateManager
}

func NewRegistryStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "state")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &State{AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot, f)}
}
