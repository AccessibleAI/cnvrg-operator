package services_check

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/services-check/tmpl"

func servicesCheck() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/job.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.JobGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func ServicesCheckState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.ServicesCheck.Enabled {
		state = append(state, servicesCheck()...)
	}

	return state
}
