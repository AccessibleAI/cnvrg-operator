package ingress

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/cnvrgapp/ingress/tmpl"

var istioGwState = []*desired.State{
	{

		TemplatePath:   path + "/istio/gw/gw.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioGwGVR],
		Own:            true,
	},
}

var istioVsState = []*desired.State{
	{
		TemplatePath:   path + "/istio/vs/app.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/istio/vs/es.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/istio/vs/grafana.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/istio/vs/kibana.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/istio/vs/minio.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/istio/vs/prom.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	if cnvrgApp.Spec.Ingress.Enabled == "false" {
		return state
	}
	if cnvrgApp.Spec.Ingress.IngressType == mlopsv1.OpenShiftIngress {

	}
	if cnvrgApp.Spec.Ingress.IngressType == mlopsv1.IstioIngress {
		state = append(state, istioGwState...)
		state = append(state, istioVsState...)
	}
	return state
}
