package logging

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/logging/tmpl"

func elastAlert() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/elastalert/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/elastalert/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/elastalert/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/elastalert/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
			TemplateData:   nil,
		},
		{

			TemplatePath:   path + "/elastalert/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			TemplateData:   nil,
		},
		{

			TemplatePath:   path + "/elastalert/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
			TemplateData:   nil,
		},
		{

			TemplatePath:   path + "/elastalert/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			TemplateData:   nil,
		},
	}
}

func KibanaConfSecret(data desired.TemplateData) []*desired.State {
	return []*desired.State{
		{
			TemplateData:   data,
			TemplatePath:   path + "/kibana/secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SecretGVR],
			Own:            true,
		},
	}
}

func kibanaIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioVsGVR],
			Own:            true,
		},
	}
}

func kibanaOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.OcpRouteGVR],
			Own:            true,
		},
	}
}

func kibana() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kibana/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kibana/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kibana/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kibana/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
	}
}

func fluentbitConfigState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/fluentbit/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
			Override:       true,
		},
	}
}

func fluentbitState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/fluentbit/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/fluentbit/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/fluentbit/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/fluentbit/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DaemonSetGVR],
			Own:            true,
		},
	}
}

func kibanaOauthProxy() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/oauth.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SecretGVR],
			Own:            true,
		},
	}
}

func CnvrgAppLoggingState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if *cnvrgApp.Spec.Logging.Elastalert.Enabled {
		state = append(state, elastAlert()...)
	}
	if *cnvrgApp.Spec.Logging.Kibana.Enabled {
		state = append(state, kibana()...)

		if *cnvrgApp.Spec.SSO.Enabled {
			state = append(state, kibanaOauthProxy()...)
		}

		if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
			state = append(state, kibanaIstioVs()...)
		}

		if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.OpenShiftIngress {
			state = append(state, kibanaOcpRoute()...)
		}
	}

	return state
}

func InfraLoggingState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State

	if *infra.Spec.Logging.Fluentbit.Enabled {
		state = append(state, fluentbitState()...)
	}

	return state
}

func FluentbitConfigurationState(data interface{}) []*desired.State {
	fluentbitConfigState := fluentbitConfigState()
	fluentbitConfigState[0].TemplateData = data
	return fluentbitConfigState
}
