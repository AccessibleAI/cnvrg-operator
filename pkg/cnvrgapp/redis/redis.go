package redis

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/cnvrgapp/redis/tmpl"

var state = []*desired.State{
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

		TemplatePath:   path + "/cm.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
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

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	if cnvrgApp.Spec.Redis.Enabled == "true" {
		return state
	}
	return nil
}
