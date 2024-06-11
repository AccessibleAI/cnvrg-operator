package dbs

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RedisStateManager struct {
	*desired.AssetsStateManager
}

func NewRedisStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "redis")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &RedisStateManager{AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/redis", f)}
}
