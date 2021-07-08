package networking

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net"
	"os"
)

const path = "/pkg/networking/tmpl"

func istioInstanceState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/istio/instance/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            false,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/istio/instance/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            false,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/istio/instance/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            false,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/istio/instance/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            false,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/istio/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            false,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/istio/instance/instance.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioGVR],
			Own:            false,
			Updatable:      true,
		},
	}
}

func ingressState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/ingress/gw.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioGwGVR],
			Own:            true,
			Updatable:      true,
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
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func InfraNetworkingState(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State

	if *cnvrgInfra.Spec.Networking.Istio.Enabled {
		state = append(state, istioInstanceState()...)
	}

	if cnvrgInfra.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress && *cnvrgInfra.Spec.Networking.Ingress.IstioGwEnabled {
		state = append(state, ingressState()...)
	}

	if *cnvrgInfra.Spec.Networking.Proxy.Enabled {
		state = append(state, proxyState()...)
	}

	return state
}

func CnvrgAppNetworkingState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress && *cnvrgApp.Spec.Networking.Ingress.IstioGwEnabled {
		state = append(state, ingressState()...)
	}

	if *cnvrgApp.Spec.Networking.Proxy.Enabled {
		state = append(state, proxyState()...)
	}

	return state
}

func IstioCrds() (crds []*desired.State) {
	err := pkger.Walk(path+"/istio/crds", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		crd := &desired.State{

			TemplatePath:   path,
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.CrdGVR],
			Own:            false,
			Updatable:      false,
		}
		crds = append(crds, crd)
		return nil
	})
	if err != nil {
		zap.S().Error(err, "error loading istio crds")
	}
	return
}

func DefaultNoProxy() []string {
	return append([]string{
		"localhost",
		"127.0.0.1",
		".svc",
		".svc.cluster.local",
		"kubernetes.default.svc",
		"kubernetes.default.svc.cluster.local"}, getK8sApiIP())
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
