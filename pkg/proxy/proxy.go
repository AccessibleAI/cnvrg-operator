package proxy

import (
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/proxy/tmpl"

func proxyState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/proxy-cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func State() []*desired.State {
	return proxyState()
}
