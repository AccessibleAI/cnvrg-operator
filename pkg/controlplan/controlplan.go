package controlplan

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/controlplan/tmpl"

var registryState = []*desired.State{
	{

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

		TemplatePath:   path + "/conf/rbac/role.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/conf/rbac/rolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleBindingGVR],
		Own:            true,
	},
	{

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

		TemplatePath:   path + "/conf/cm/config-base.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/conf/cm/config-networking.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/conf/cm/secret-base.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/conf/cm/secret-ldap.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/conf/cm/secret-object-storage.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/conf/cm/secret-sso.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
}

var webAppState = []*desired.State{
	{

		TemplatePath:   path + "/webapp/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/webapp/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/webapp/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
}

var sidekiqState = []*desired.State{
	{

		TemplatePath:   path + "/sidekiqs/sidekiq.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var searchkiqState = []*desired.State{
	{

		TemplatePath:   path + "/sidekiqs/searchkiq.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var systemkiqState = []*desired.State{
	{

		TemplatePath:   path + "/sidekiqs/systemkiq.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var hyperState = []*desired.State{
	{

		TemplatePath:   path + "/hyper/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/hyper/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
}

var pgState = []*desired.State{
	{

		TemplatePath:   path + "/pg/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/pg/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/pg/secret.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/pg/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/pg/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
}

var redisState = []*desired.State{
	{

		TemplatePath:   path + "/redis/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/redis/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/redis/cm.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/redis/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var singleBackendMinio = []*desired.State{

	{

		TemplatePath:   path + "/minio/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/minio/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/minio/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/minio/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
}

var sharedBackendMinio = []*desired.State{

	{

		TemplatePath:   path + "/minio/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},

	{

		TemplatePath:   path + "/minio/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/minio/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/minio/sh-dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/minio/dr.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioDestinationRuleGVR],
		Own:            true,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	state = append(state, registryState...)
	state = append(state, rbacState...)
	state = append(state, controlPlanConfigState...)

	if cnvrgApp.Spec.ControlPlan.Pg.Enabled == "true" {
		state = append(state, pgState...)
	}

	if cnvrgApp.Spec.ControlPlan.Redis.Enabled == "true" {
		state = append(state, redisState...)
	}

	if cnvrgApp.Spec.ControlPlan.Minio.Enabled == "true" && cnvrgApp.Spec.ControlPlan.Minio.SharedStorage.Enabled == "true" {
		state = append(state, sharedBackendMinio...)
	} else {
		state = append(state, singleBackendMinio...)
	}

	if cnvrgApp.Spec.ControlPlan.WebApp.Enabled == "true" {
		state = append(state, webAppState...)
	}

	if cnvrgApp.Spec.ControlPlan.Sidekiq.Enabled == "true" && cnvrgApp.Spec.ControlPlan.Sidekiq.Split == "true" {
		state = append(state, sidekiqState...)
	}

	if cnvrgApp.Spec.ControlPlan.Searchkiq.Enabled == "true" && cnvrgApp.Spec.ControlPlan.Sidekiq.Split == "true" {
		state = append(state, searchkiqState...)
	}

	if cnvrgApp.Spec.ControlPlan.Systemkiq.Enabled == "true" && cnvrgApp.Spec.ControlPlan.Sidekiq.Split == "true" {
		state = append(state, systemkiqState...)
	}

	if cnvrgApp.Spec.ControlPlan.Sidekiq.Enabled == "true" && cnvrgApp.Spec.ControlPlan.Sidekiq.Split == "false" {
		state = append(state, sidekiqState...)
	}

	if cnvrgApp.Spec.ControlPlan.Hyper.Enabled == "true" {
		state = append(state, hyperState...)
	}
	return state
}
