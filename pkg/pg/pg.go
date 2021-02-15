package pg

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/pg/tmpl"

var state = []*desired.State{
	{
		Name:           "pg-pvc",
		TemplatePath:   path + "/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.PvcGVR,
	},
	{
		Name:           "pg-dep",
		TemplatePath:   path + "/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.DeploymentGVR,
	},
	{
		Name:           "pg-secret",
		TemplatePath:   path + "/secret.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.SecretGVR,
	},
	{
		Name:           "pg-svc",
		TemplatePath:   path + "/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.PvcGVR,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	if cnvrgApp.Spec.Pg.Enabled == "true" {
		return state
	}
	return nil
}
