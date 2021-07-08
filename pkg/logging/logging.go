package logging

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
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
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/elastalert/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/elastalert/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/elastalert/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
			TemplateData:   nil,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/elastalert/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			TemplateData:   nil,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/elastalert/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
			TemplateData:   nil,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/elastalert/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			TemplateData:   nil,
			Updatable:      true,
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
			Updatable:      true,
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
			Updatable:      true,
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
			Updatable:      true,
		},
	}
}

func kibanaIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IngressGVR],
			Own:            true,
			Updatable:      true,
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
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/kibana/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
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
			Updatable:      true,
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
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/fluentbit/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DaemonSetGVR],
			Own:            true,
			Updatable:      true,
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

		switch cnvrgApp.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, kibanaIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, kibanaIngress()...)
		case mlopsv1.OpenShiftIngress:
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
