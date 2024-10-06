package sso

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type CentralStateManager struct {
	*desired.AssetsStateManager
	app    *mlopsv1.CnvrgApp
	client client.Client
}

func NewCentralStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "ssoCentral")
	f := &desired.LoadFilter{DefaultLoader: true, Ingress: &app.Spec.Networking.Ingress.Type}
	return &CentralStateManager{
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/central", f),
		app:                app,
	}
}

func (c *CentralStateManager) renderSsoConfigs() error {
	assets := []string{"proxycfg.tpl"}
	f := &desired.LoadFilter{AssetName: assets}
	cfg := desired.NewAssetsGroup(fs, c.RootPath(), c.Log(), f)
	if err := cfg.LoadAssets(); err != nil {
		return err
	}

	configData, err := c.proxyCfgData()
	if err != nil {
		return err
	}

	if err = cfg.Render(configData); err != nil {
		return err
	}

	c.AddToState(cfg)

	return nil
}

func (c *CentralStateManager) renderDeploymentAndSvc() error {
	assets := []string{"dep.tpl", "svc.tpl"}
	f := &desired.LoadFilter{AssetName: assets}
	dep := desired.NewAssetsGroup(fs, c.RootPath(), c.Log(), f)

	if err := dep.LoadAssets(); err != nil {
		return err
	}

	if err := dep.Render(c.depData()); err != nil {
		return err
	}

	c.AddToState(dep)

	return nil
}

func (c *CentralStateManager) depData() map[string]interface{} {
	return map[string]interface{}{
		"Namespace":   c.app.Namespace,
		"SsoDomainId": strings.Split(c.app.Spec.ClusterDomain, ".")[0],
		"Spec":        c.app.Spec,
		"AppUrl": fmt.Sprintf("%s://%s%s.%s", c.schema(),
			c.app.Spec.ControlPlane.WebApp.SvcName,
			c.app.Spec.Networking.ClusterDomainPrefix.Prefix,
			c.app.Spec.ClusterDomain,
		),
	}
}

func (c *CentralStateManager) proxyCfgData() (map[string]interface{}, error) {
	var groups []string
	if c.app.Spec.SSO.Central.GroupsAuth {
		groups = append(groups, c.domainId())
	}

	var clientId, clientSecret string

	// if credentials secret ref is set, get clientId and clientSecret from the secret
	if c.app.Spec.SSO.Central.CredentialsSecretRef != "" {
		credentialsSecret := &corev1.Secret{}
		err := c.client.Get(context.Background(), client.ObjectKey{Namespace: c.app.Namespace, Name: c.app.Spec.SSO.Central.CredentialsSecretRef}, credentialsSecret)
		if err != nil {
			return nil, err
		}

		if _, ok := credentialsSecret.Data["ClientId"]; !ok {
			return nil, fmt.Errorf("credentialSecretRef configured for SSO, but clientId not found in secret %s", c.app.Spec.SSO.Central.CredentialsSecretRef)
		}
		if _, ok := credentialsSecret.Data["ClientSecret"]; !ok {
			return nil, fmt.Errorf("credentialSecretRef configured for SSO, but clientSecret not found in secret %s", c.app.Spec.SSO.Central.CredentialsSecretRef)
		}

		clientId = string(credentialsSecret.Data["clientId"])
		clientSecret = string(credentialsSecret.Data["clientSecret"])
	}

	if c.app.Spec.SSO.Central.ClientID != "" {
		clientId = c.app.Spec.SSO.Central.ClientID
	}

	if c.app.Spec.SSO.Central.ClientSecret != "" {
		clientSecret = c.app.Spec.SSO.Central.ClientSecret
	}

	d := map[string]interface{}{
		"Namespace":    c.app.Namespace,
		"EmailDomain":  c.app.Spec.SSO.Central.EmailDomain,
		"Provider":     c.app.Spec.SSO.Central.Provider,
		"ClientId":     clientId,
		"ClientSecret": clientSecret,
		"RedirectUrl": fmt.Sprintf("%s://%s%s.%s/oauth2/callback",
			c.schema(),
			c.app.Spec.SSO.Central.SvcName,
			c.app.Spec.Networking.ClusterDomainPrefix.Prefix,
			c.app.Spec.ClusterDomain),
		"OidcIssuerURL":                    c.app.Spec.SSO.Central.OidcIssuerURL,
		"Scope":                            c.app.Spec.SSO.Central.Scope,
		"InsecureOidcAllowUnverifiedEmail": c.app.Spec.SSO.Central.InsecureOidcAllowUnverifiedEmail,
		"SslInsecureSkipVerify":            c.app.Spec.SSO.Central.SslInsecureSkipVerify,
		"WhitelistDomain":                  c.app.Spec.SSO.Central.WhitelistDomain,
		"CookieDomain":                     c.app.Spec.SSO.Central.CookieDomain,
		"ExtraJwtIssuer":                   c.jwksUrlWithAudience(),
		"Groups":                           groups,
	}
	return d, nil
}

func (c *CentralStateManager) domainId() string {
	return strings.Split(c.app.Spec.ClusterDomain, ".")[0]
}

func (c *CentralStateManager) schema() string {
	schema := "http"
	if c.app.Spec.Networking.HTTPS.Enabled {
		schema = "https"
	}
	return schema
}

func (c *CentralStateManager) jwksUrlWithAudience() string {

	return fmt.Sprintf("%s=cnvrg-tenant", c.app.Spec.SSO.Central.JwksURL)
}

func (c *CentralStateManager) Apply() error {
	if err := c.renderDeploymentAndSvc(); err != nil {
		return err
	}

	if err := c.renderSsoConfigs(); err != nil {
		return err
	}

	if err := c.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}
