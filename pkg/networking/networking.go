package networking

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
)

const path = "/pkg/networking/tmpl"

var istioState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/istio/clusterrole.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/clusterrolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
	},
	{
		Name:           "",
		TemplatePath:   path + "/istio/instance.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioGVR],
	},
}

var ocpRoutesState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/route/app.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.OcpRouteGVR],
	},
}

var istioVsState = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/vs/app.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	if cnvrgApp.Spec.Networking.Enabled == "false" {
		return state
	}
	if cnvrgApp.Spec.Networking.Istio.Enabled == "true" {
		state = append(state, istioState...)
	}
	if cnvrgApp.Spec.Networking.IngressType == mlopsv1.IstioIngress {
		// set istio VSs and GWs here
	}
	if cnvrgApp.Spec.Networking.IngressType == mlopsv1.OpenShiftIngress {
		state = append(state, ocpRoutesState...)
	}
	if cnvrgApp.Spec.Networking.IngressType == mlopsv1.IstioIngress {
		state = append(state, istioVsState...)
	}
	return state
}

func LoadCrds() (crds []*desired.State) {
	err := pkger.Walk(path+"/istio/crds", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		crd := &desired.State{
			Name:           "",
			TemplatePath:   path,
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.CrdGVR],
		}
		if err := crd.GenerateDeployable(nil); err != nil {
			zap.S().Error(err, "error loading istio crds")

		}
		if err := crd.Apply(); err != nil {
			zap.S().Error(err, "error applying crd")
		}
		crds = append(crds, crd)
		return nil
	})
	if err != nil {
		zap.S().Error(err, "error loading istio crds")
	}
	return
}
