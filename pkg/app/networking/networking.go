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

func NewIstioGwState(app *mlopsv1.CnvrgApp, c client.Client, s *runtime.Scheme, log logr.Logger) desired.StateManager {
	l := log.WithValues("stateManager", "istioGw")
	f := &desired.LoadFilter{DefaultLoader: true}
	return &IstioGwState{AssetsStateManager: desired.NewAssetsStateManager(app, c, s, l, fs, fsRoot+"/ingress", f)}
}

//
//func ingressState() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/ingress/gw.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IstioGwGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func proxyState() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/proxy/proxy-cm.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.ConfigMapGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func AppNetworkingState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
//	var state []*desired.State
//
//	if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress && cnvrgApp.Spec.Networking.Ingress.IstioGwEnabled {
//		state = append(state, ingressState()...)
//	}
//
//	if cnvrgApp.Spec.Networking.Proxy.Enabled {
//		state = append(state, proxyState()...)
//	}
//
//	return state
//}
//

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
