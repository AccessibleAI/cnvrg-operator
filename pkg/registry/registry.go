package registry

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/cnvrginfra/registry/tmpl"

var registryState = []*desired.State{
	{
		TemplatePath:   path + "/secret.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
}

func State(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	state = append(state, registryState...)
	return state
}
