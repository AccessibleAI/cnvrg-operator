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
}

func IstioState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	if cnvrgApp.Spec.Networking.Enabled == "true" && cnvrgApp.Spec.Networking.Istio.Enabled == "true" {
		return istioState
	}
	return nil
}
