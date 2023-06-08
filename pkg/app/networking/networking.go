package networking

import (
	"embed"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const fsRoot = "tmpl"

//go:embed  tmpl/*
var fs embed.FS

type IstioGwState struct {
	*desired.AssetsStateManager
}

type HttpProxyState struct {
	*desired.AssetsStateManager
}

func NewIstioGwState(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "istioGw")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &IstioGwState{AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/ingress", f)}
}

func NewHttpProxyState(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "httpProxyState")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &HttpProxyState{AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/proxy", f)}
}

func DefaultNoProxy(clusterInternalDomain string) []string {
	return append([]string{
		"localhost",
		"127.0.0.1",
		".svc",
		fmt.Sprintf(".svc.%s", clusterInternalDomain),
		"kubernetes.default.svc",
		fmt.Sprintf("kubernetes.default.svc.%s", clusterInternalDomain)}, getK8sApiIP())
}

func getK8sApiIP() string {
	k8sApiDns := "kubernetes.default.svc"
	zap.S().Infof("getting K8s API (%s) IP", k8sApiDns)

	ips, err := net.LookupIP(k8sApiDns)
	if err != nil {
		zap.S().Errorf("%s: getting K8s api IP", err)
		return "unable-to-detect-k8s-api-ip"
	}
	if len(ips) < 1 {
		zap.S().Errorf("%s: getting K8s api IP", err)
		return "unable-to-detect-k8s-api-ip"
	}
	return ips[0].String()
}
