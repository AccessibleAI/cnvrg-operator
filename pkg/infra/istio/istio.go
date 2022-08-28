package istio

import (
	"embed"
	"fmt"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS // TODO: this is bad, but I've to hurry up

func istioInstanceState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/istio/instance/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent, // TODO: this is bad, but I've to hurry up
		},
		{

			TemplatePath:   path + "/istio/instance/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/istio/instance/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/istio/instance/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/istio/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/istio/instance/instance.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func IstioCrds() (crds []*desired.State) {
	d, err := templatesContent.ReadDir(path + "/istio/crds")
	if err != nil {
		zap.S().Error(err, "error loading istio crds")
	}
	for _, f := range d {
		crd := &desired.State{

			TemplatePath:   path + "/istio/crds/" + f.Name(),
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.CrdGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		}
		crds = append(crds, crd)
	}

	return
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

func IstioState() []*desired.State {
	var state []*desired.State
	state = append(state, istioInstanceState()...)
	return state
}
