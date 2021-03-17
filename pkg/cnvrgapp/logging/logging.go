package logging

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/cnvrgapp/logging/tmpl"

var es = []*desired.State{
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
}

var elastAlert = []*desired.State{
	{

		TemplatePath:   path + "/elastalert/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/elastalert/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/elastalert/cm.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/elastalert/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var kibana = []*desired.State{
	{

		TemplatePath:   path + "/kibana/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{

		TemplatePath:   path + "/kibana/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Es.Enabled == "true" {
		state = append(state, es...)
	}
	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Elastalert.Enabled == "true" {
		state = append(state, elastAlert...)
	}
	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Kibana.Enabled == "true" {
		state = append(state, kibana...)
	}

	return state
}
