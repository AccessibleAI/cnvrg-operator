package sso

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/sso/tmpl"

func pki(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/private-key-secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/public-key-secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
		},
	}
}

func SsoState(cnvrgApp *mlopsv1.CnvrgApp, data interface{}) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.SSO.Pki.Enabled {
		state = append(state, pki(data)...)
	}

	return state
}
