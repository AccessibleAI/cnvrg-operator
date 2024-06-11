package dbs

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MinioStateManager struct {
	*desired.AssetsStateManager
	pgSecret *desired.AssetsGroup
}

func NewMinioStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "minio")
	f := &desired.LoadFilter{Ingress: &app.Spec.Networking.Ingress.Type, DefaultLoader: true}
	asm := desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/minio", f)
	return &ElasticStateManager{AssetsStateManager: asm, app: app}
}
