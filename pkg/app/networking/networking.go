package networking

import (
	"embed"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS // TODO: this is bad, but I've to hurry up

func ingressState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/ingress/gw.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioGwGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func proxyState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/proxy/proxy-cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func AppNetworkingState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress && cnvrgApp.Spec.Networking.Ingress.IstioGwEnabled {
		state = append(state, ingressState()...)
	}

	if cnvrgApp.Spec.Networking.Proxy.Enabled {
		state = append(state, proxyState()...)
	}

	return state
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
