package gpu

import (
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/gpu/tmpl"

func nvidiaDp() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/nvidiadp/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
			TemplateData:   nil,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/nvidiadp/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DaemonSetGVR],
			Own:            true,
			TemplateData:   nil,
		},
	}
}

func NvidiaDpState(data interface{}) []*desired.State {
	nvidiaDp := nvidiaDp()
	nvidiaDp[0].TemplateData = data
	nvidiaDp[1].TemplateData = data
	return nvidiaDp
}
