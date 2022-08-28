package istio

import (
	"embed"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS // TODO: this is bad, but I've to hurry up

func istioInstanceState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/instance/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/instance/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/instance/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/instance/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/instance/instance.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		},
	}
}

func Crds() (crds []*desired.State) {
	d, err := templatesContent.ReadDir(path + "/crds")
	if err != nil {
		zap.S().Error(err, "error loading istio crds")
	}
	for _, f := range d {
		crd := &desired.State{

			TemplatePath:   path + "/crds/" + f.Name(),
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.CrdGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		}
		crds = append(crds, crd)
	}

	return
}

func ApplyState() []*desired.State {
	var state []*desired.State
	state = append(state, istioInstanceState()...)
	return state
}
