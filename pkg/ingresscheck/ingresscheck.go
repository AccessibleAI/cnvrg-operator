package ingresscheck

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/ingresscheck/tmpl"

func ingressCheck() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/job.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.JobGVK],
			Own:            true,
			Updatable:      false,
		},
	}
}

func IngressCheckState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.IngressCheck.Enabled {
		state = append(state, ingressCheck()...)
	}

	return state
}
