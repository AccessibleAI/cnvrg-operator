package reloader

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS

func reloaderState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func State(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if cnvrgInfra.Spec.ConfigReloader.Enabled {
		state = append(state, reloaderState()...)
	}
	return state
}
