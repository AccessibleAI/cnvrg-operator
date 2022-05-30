package monitoring

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS // TODO: this is bat, but I've to hurry up

func PromCreds(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/credsec.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			TemplateData:   data,
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
	}
}

func PromUpstreamCreds(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/ccp/upstream.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			TemplateData:   data,
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
	}
}

func prometheusOperatorState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/operator/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/operator/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/operator/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/operator/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/operator/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func infraPrometheusInstanceState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/infra/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/prometheus.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PrometheusGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/rules.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PrometheusRuleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func kubeStateMetricsState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kube-state-metrics/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func GrafanaDSState(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplateData:   data,
			TemplatePath:   path + "/grafana/datasource.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func grafanaState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/grafana/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/grafana/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/grafana/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/grafana/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/grafana/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/grafana/dashboards.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/grafana/credsec.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}

}

func nodeExporterState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/node-exporter/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/node-exporter/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/node-exporter/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/node-exporter/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/node-exporter/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DaemonSetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/node-exporter/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func ccpPrometheusInstance() []*desired.State {
	return []*desired.State{

		{
			TemplatePath:   path + "/prometheus/instance/ccp/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/ccp/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            false,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/ccp/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/ccp/prometheus.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PrometheusGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/prometheus/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func promOauthProxy() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/prom-auth-proxy-cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func promIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func promIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func promOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/instance/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func grafanaIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/grafana/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func grafanaIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/grafana/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func grafanaOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/grafana/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func defaultServiceMonitors() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/default-servicemonitors/apiserver.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/controller-manager.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/coredns.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/kubelet.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/scheduler.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/metagpu.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func dcgmExporter() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/dcgm-exporter/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DaemonSetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/dcgm-exporter/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/dcgm-exporter/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/dcgm-exporter/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func habanaExporter() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/habana-exporter/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DaemonSetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/habana-exporter/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/habana-exporter/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		{
			TemplatePath:   path + "/habana-exporter/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func grafanaOauthProxy() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/grafana/oauth.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func cnvrgIdleMetricsServiceMonitors() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/cnvrg-servicemonitors/idle-metrics.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ServiceMonitorGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func AppMonitoringState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.Monitoring.Prometheus.Enabled {
		state = append(state, promOauthProxy()...)
		state = append(state, ccpPrometheusInstance()...)
		switch cnvrgApp.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, promIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, promIngress()...)
		case mlopsv1.OpenShiftIngress:
			state = append(state, promOcpRoute()...)
		}
	}

	if cnvrgApp.Spec.Monitoring.Grafana.Enabled {
		state = append(state, grafanaState()...)
		if cnvrgApp.Spec.SSO.Enabled {
			state = append(state, grafanaOauthProxy()...)
		}
		switch cnvrgApp.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, grafanaIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, grafanaIngress()...)
		case mlopsv1.OpenShiftIngress:
			state = append(state, grafanaOcpRoute()...)
		}
	}

	if cnvrgApp.Spec.Monitoring.CnvrgIdleMetricsExporter.Enabled {
		state = append(state, cnvrgIdleMetricsServiceMonitors()...)
	}

	return state
}

func InfraMonitoringState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State

	if infra.Spec.Monitoring.PrometheusOperator.Enabled {
		state = append(state, prometheusOperatorState()...)
	}

	if infra.Spec.Monitoring.Prometheus.Enabled {
		state = append(state, promOauthProxy()...)
		state = append(state, infraPrometheusInstanceState()...)

		switch infra.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, promIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, promIngress()...)
		case mlopsv1.OpenShiftIngress:
			state = append(state, promOcpRoute()...)
		}

	}
	if infra.Spec.Monitoring.KubeStateMetrics.Enabled {
		state = append(state, kubeStateMetricsState()...)
	}

	if infra.Spec.Monitoring.Grafana.Enabled {
		state = append(state, grafanaState()...)
		if infra.Spec.SSO.Enabled {
			state = append(state, grafanaOauthProxy()...)
		}
		switch infra.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, grafanaIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, grafanaIngress()...)
		case mlopsv1.OpenShiftIngress:
			state = append(state, grafanaOcpRoute()...)
		}
	}

	if infra.Spec.Monitoring.NodeExporter.Enabled {
		state = append(state, nodeExporterState()...)
	}

	if infra.Spec.Monitoring.DefaultServiceMonitors.Enabled {
		state = append(state, defaultServiceMonitors()...)
	}

	if infra.Spec.Monitoring.DcgmExporter.Enabled {
		state = append(state, dcgmExporter()...)
	}

	if infra.Spec.Monitoring.HabanaExporter.Enabled {
		state = append(state, habanaExporter()...)
	}

	if infra.Spec.Monitoring.CnvrgIdleMetricsExporter.Enabled {
		state = append(state, cnvrgIdleMetricsServiceMonitors()...)
	}

	return state
}

func Crds() (crds []*desired.State) {
	d, err := templatesContent.ReadDir(path + "/crds")
	if err != nil {
		zap.S().Error(err, "error loading prometheus crds")
	}
	for _, f := range d {
		crd := &desired.State{

			TemplatePath:   path + "/crds/" + f.Name(),
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.CrdGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		}
		crds = append(crds, crd)
	}
	return
}
