package priorityclass

import (
	"embed"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS

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
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/jobs-class.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PriorityClassGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func State() []*desired.State {
	return priorityClassState()
}
