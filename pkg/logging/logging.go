package logging

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/logging/tmpl"

func ElastCreds(data *desired.TemplateData) []*desired.State {
	return []*desired.State{
		{
			TemplateData:   data,
			TemplatePath:   path + "/elastalert/credsec.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
		},
	}
}

func ElastAlert() []*desired.State {
	return []*desired.State{
		{
			TemplatePath: path + "/elastalert/sa.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.SaGVK],
			Own:          true,
			Updatable:    false,
		},
		{
			TemplatePath: path + "/elastalert/authproxycm.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.ConfigMapGVK],
			Own:          true,
			Updatable:    true,
		},
		{
			TemplatePath: path + "/elastalert/role.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.RoleGVK],
			Own:          true,
			Updatable:    true,
		},
		{
			TemplatePath: path + "/elastalert/rolebinding.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.RoleBindingGVK],
			Own:          true,
			Updatable:    true,
		},
		{
			TemplatePath: path + "/elastalert/pvc.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.PvcGVK],
			Own:          true,
			Updatable:    false,
		},
		{

			TemplatePath: path + "/elastalert/svc.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.SvcGVK],
			Own:          true,
			Updatable:    true,
		},
		{

			TemplatePath: path + "/elastalert/cm.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.ConfigMapGVK],
			Own:          true,
			Updatable:    true,
		},
		{

			TemplatePath: path + "/elastalert/dep.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.DeploymentGVK],
			Own:          true,
			Updatable:    true,
		},
		{

			TemplatePath: path + "/elastalert/vs.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.IstioVsGVK],
			Own:          true,
			Updatable:    true,
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
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func kibanaIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath: path + "/kibana/vs.tpl",
			Obj:          &unstructured.Unstructured{},
			GVK:          desired.Kinds[desired.IstioVsGVK],
			Own:          true,
			Updatable:    true,
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
			GVK:            desired.Kinds[desired.OcpRouteGVK],
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
			GVK:            desired.Kinds[desired.IngressGVK],
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
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/kibana/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/kibana/credsec.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
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
			GVK:            desired.Kinds[desired.ConfigMapGVK],
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
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/fluentbit/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/fluentbit/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DaemonSetGVK],
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
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func CnvrgAppKibanaState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	state = append(state, kibana()...)

	if cnvrgApp.Spec.SSO.Enabled {
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

	return state
}

func InfraLoggingState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State

	if infra.Spec.Logging.Fluentbit.Enabled {
		state = append(state, fluentbitState()...)
	}

	return state
}

func FluentbitConfigurationState(data interface{}) []*desired.State {
	fluentbitConfigState := fluentbitConfigState()
	fluentbitConfigState[0].TemplateData = data
	return fluentbitConfigState
}
