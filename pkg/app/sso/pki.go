package sso

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"embed"
	"encoding/pem"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const fsRoot = "tmpl"

//go:embed  tmpl/*
var fs embed.FS

type PkiStateManager struct {
	*desired.AssetsStateManager
	app        *mlopsv1.CnvrgApp
	PrivateKey string
	PublicKey  string
}

func NewPkiStateManager(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "ssoPki")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &PkiStateManager{
		AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/pki", f),
		app:                app,
	}
}

func (p *PkiStateManager) domainId() string {
	return strings.Split(p.app.Spec.ClusterDomain, ".")[0]
}

func (p *PkiStateManager) generate() error {
	pkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	// private key
	privatePemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pkey)}
	p.PrivateKey = string(pem.EncodeToMemory(privatePemBlock))

	// public key
	b, err := x509.MarshalPKIXPublicKey(&pkey.PublicKey)
	if err != nil {
		return err
	}
	publicPemBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: b}
	p.PublicKey = string(pem.EncodeToMemory(publicPemBlock))

	return nil
}

func (p *PkiStateManager) RenderDeployment() error {
	assets := []string{"private-key-secret.tpl", "public-key-secret.tpl"}
	f := &desired.LoadFilter{AssetName: assets}
	pki := desired.NewAssetsGroup(fs, p.RootPath(), p.Log(), f)
	if err := pki.LoadAssets(); err != nil {
		return err
	}

	if err := p.generate(); err != nil {
		return err
	}

	data := map[string]interface{}{
		"Namespace":        p.app.Namespace,
		"PrivateKeySecret": p.app.Spec.SSO.Pki.PrivateKeySecret,
		"PrivateKey":       p.PrivateKey,
		"PublicKeySecret":  p.app.Spec.SSO.Pki.PublicKeySecret,
		"PublicKey":        p.PublicKey,
		"DomainID":         p.domainId(),
	}

	if err := pki.Render(data); err != nil {
		return err
	}

	p.AddToState(pki)

	return nil
}

func (p *PkiStateManager) Apply() error {

	if err := p.RenderDeployment(); err != nil {
		return err
	}

	if err := p.AssetsStateManager.Apply(); err != nil {
		return err
	}

	return nil
}
