package networking

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/networking/tmpl"

var istioState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/istio/crds.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.CrdGVR,
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/clusterrole.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.ClusterRoleGVR,
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/clusterrolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.ClusterRoleBindingGVR,
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.DeploymentGVR,
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.SaGVR,
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.SvcGVR,
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/instance.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.IstioGVR,
	},
}

var ocpRoutesState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/route/app.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.OcpRouteGVR,
	},
}

var istioVsState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/vs/app.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.IstioVsGVR,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	if cnvrgApp.Spec.Networking.Enabled == "true" && cnvrgApp.Spec.Networking.Istio.Enabled == "true" {
		state = append(state, istioState...)
	}
	if cnvrgApp.Spec.Networking.IngressType == mlopsv1.OpenShiftIngress {
		state = append(state, ocpRoutesState...)
	}
	if cnvrgApp.Spec.Networking.IngressType == mlopsv1.IstioIngress {
		state = append(state, istioVsState...)
	}
	return state
}
