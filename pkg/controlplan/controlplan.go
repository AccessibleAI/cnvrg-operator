package controlplan

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/controlplan/tmpl"

var registryState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/conf/registry/secret.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
}

var rbacState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/conf/rbac/role.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/conf/rbac/rolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleBindingGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/conf/rbac/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
		Own:            true,
	},
}

var controlPlanConfigState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/conf/cm/control-plan-core.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
}


var webAppState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/webapp/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/webapp/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	//state = append(state, registryState...)
	//state = append(state, rbacState...)
	state = append(state, controlPlanConfigState...)
	if cnvrgApp.Spec.ControlPlan.WebApp.Enabled == "true" {
		state = append(state, webAppState...)
	}
	return state
}
