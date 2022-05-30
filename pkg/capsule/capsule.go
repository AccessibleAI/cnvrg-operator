package capsule

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS // TODO: this is bat, but I've to hurry up (as usual :( )

func capsuleState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent, // TODO: this is bat, but I've to hurry up (as usual :( )
		},
		{
			TemplatePath:   path + "/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
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

func State(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if infra.Spec.Capsule.Enabled {
		state = append(state, capsuleState()...)
	}
	return state
}
