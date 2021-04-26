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

func istioInstanceState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/istio/instance/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            false,
		},
		{

			TemplatePath:   path + "/istio/instance/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/istio/instance/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/istio/instance/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/istio/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/istio/instance/instance.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioGVR],
			Own:            true,
		},
	}
}

func ingressState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/ingress/gw.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioGwGVR],
			Own:            false,
		},
	}
}

func IstioInstanceState(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if cnvrgInfra.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
		if *cnvrgInfra.Spec.Networking.Istio.Enabled {
			state = append(state, istioInstanceState()...)
		}
		if cnvrgInfra.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
			state = append(state, ingressState()...)
		}
	}
	return state
}

func CnvrgAppNetworkingState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
		state = append(state, ingressState()...)
	}

	return state
}

func IstioCrds() (crds []*desired.State) {
	err := pkger.Walk(path+"/istio/crds", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		crd := &desired.State{

			TemplatePath:   path,
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.CrdGVR],
			Own:            false,
		}
		crds = append(crds, crd)
		return nil
	})
	if err != nil {
		zap.S().Error(err, "error loading istio crds")
	}
	return
}
