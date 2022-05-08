package jwks

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/jwks/tmpl"

func jwks() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func JwksState(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State

	if cnvrgInfra.Spec.Jwks.Enabled {
		state = append(state, jwks()...)
	}

	return state
}
