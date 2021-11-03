package hpu

import (
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/hpu/tmpl"

func hpuDp() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/hpudp/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DaemonSetGVK],
			Own:            true,
			TemplateData:   nil,
		},
	}
}

func HpuDpState(data interface{}) []*desired.State {
	hpuDp := hpuDp()
	hpuDp[0].TemplateData = data
	return hpuDp
}
