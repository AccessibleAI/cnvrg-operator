package dbs

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/dbs/tmpl"

var esState = []*desired.State{
	{
		TemplatePath:   path + "/es/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/es/sts.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.StatefulSetGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/es/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/es/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
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

func AppDbsState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.Dbs.Pg.Enabled == "true" {
		state = append(state, pgState...)
	}

	if cnvrgApp.Spec.Dbs.Redis.Enabled == "true" {
		state = append(state, redisState...)
	}

	if cnvrgApp.Spec.Dbs.Minio.Enabled == "true" && cnvrgApp.Spec.Dbs.Minio.SharedStorage.Enabled == "true" {
		state = append(state, sharedBackendMinio...)
	} else {
		state = append(state, singleBackendMinio...)
	}

	if cnvrgApp.Spec.Dbs.Es.Enabled == "true" {
		state = append(state, esState...)
	}
	return state
}

func InfraDbsState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if infra.Spec.Dbs.Redis.Enabled == "true" {
		state = append(state, redisState...)
	}
	return state

}
