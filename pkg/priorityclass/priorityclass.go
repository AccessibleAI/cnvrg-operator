package priorityclass

import (
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/priorityclass/tmpl"

func priorityClassState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/apps-class.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PriorityClassGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/jobs-class.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PriorityClassGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func State() []*desired.State {
	return priorityClassState()
}
