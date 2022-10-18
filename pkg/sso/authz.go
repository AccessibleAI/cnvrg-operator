package sso

import (
	"errors"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AuthzStateManager struct {
	*desired.AssetsStateManager
	app *mlopsv1.CnvrgApp
}

func NewAuthzStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "authz")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &AuthzStateManager{
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/authz", f),
		app:                app,
	}
}

func (a *AuthzStateManager) Render() error {
	assets := []string{"dep.tpl", "svc.tpl"}
	f := &desired.LoadFilter{AssetName: assets}
	authz := desired.NewAssetsGroup(fs, a.RootPath(), a.Log(), f)

	if err := authz.LoadAssets(); err != nil {
		return err
	}

	if err := authz.Render(a.data()); err != nil {
		return err
	}

	a.AddToState(authz)

	return nil

}

func (a *AuthzStateManager) data() map[string]string {
	return map[string]string{
		"Namespace":   a.app.Namespace,
		"Image":       a.app.Spec.SSO.Authz.Image,
		"IngressType": a.ingressType(),
	}
}

func (a *AuthzStateManager) ingressType() string {
	switch a.app.Spec.Networking.Ingress.Type {
	case mlopsv1.IstioIngress:
		return "vs"
	case mlopsv1.NginxIngress:
		return "ingress"
	case mlopsv1.OpenShiftIngress:
		return "route"
	}
	a.Log().Error(errors.New("wrong ingress type for authz server"), "")
	return ""
}

func (a *AuthzStateManager) Apply() error {
	if err := a.Render(); err != nil {
		return err
	}

	if err := a.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}
