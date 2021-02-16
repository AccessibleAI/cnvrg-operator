package pg

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/pg/tmpl"

var state = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/secret.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	if cnvrgApp.Spec.Pg.Enabled == "true" {
		return state
	}
	return nil
}
