package registry

import (
	"embed"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS

func registryState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func State(data interface{}) []*desired.State {
	registry := registryState()
	registry[0].TemplateData = data
	return registry
}
