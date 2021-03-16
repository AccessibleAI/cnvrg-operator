package minio

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/minio/tmpl"

var singleBackendMinio = []*desired.State{

	{

		TemplatePath:   path + "/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var sharedBackendMinio = []*desired.State{

	{

		TemplatePath:   path + "/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/sh-dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/dr.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioDestinationRuleGVR],
		Own:            true,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	if cnvrgApp.Spec.Minio.Enabled == "true" && cnvrgApp.Spec.Minio.SharedStorage.Enabled == "true" {
		return sharedBackendMinio
	}
	return singleBackendMinio
}
