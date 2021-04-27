package monitoring

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
)

const path = "/pkg/monitoring/tmpl"

func prometheusOperatorState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prometheus/operator/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/operator/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/operator/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/operator/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/operator/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/prometheus.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PrometheusGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/infra/rules.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PrometheusRuleGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/kube-state-metrics/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.SecretGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/grafana/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/grafana/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/grafana/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/grafana/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/grafana/dashboards.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/node-exporter/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/node-exporter/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/node-exporter/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/node-exporter/ds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DaemonSetGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/node-exporter/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            false,
		},
		{
			TemplatePath:   path + "/prometheus/instance/ccp/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            false,
		},
		{
			TemplatePath:   path + "/prometheus/instance/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/ccp/prometheus.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PrometheusGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/prometheus/instance/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.ConfigMapGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.IstioVsGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.OcpRouteGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.IstioVsGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.OcpRouteGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/controller-manager.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/coredns.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/kubelet.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/default-servicemonitors/scheduler.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.DaemonSetGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/dcgm-exporter/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/dcgm-exporter/servicemonitor.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},

		{
			TemplatePath:   path + "/dcgm-exporter/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
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
			GVR:            desired.Kinds[desired.SecretGVR],
			Own:            true,
		},
	}
}

func appServiceMonitors() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/default-servicemonitors/cnvrg-idle-metrics.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ServiceMonitorGVR],
			Own:            true,
		},
	}
}

func AppMonitoringState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	state = appServiceMonitors()

	if *cnvrgApp.Spec.Monitoring.Prometheus.Enabled {
		state = append(state, ccpPrometheusInstance()...)
	}
	if *cnvrgApp.Spec.Monitoring.Grafana.Enabled {
		state = append(state, grafanaState()...)
	}
	if *cnvrgApp.Spec.SSO.Enabled {
		state = append(state, promOauthProxy()...)
	}
	if *cnvrgApp.Spec.SSO.Enabled {
		state = append(state, grafanaOauthProxy()...)
	}

	if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
		if *cnvrgApp.Spec.Monitoring.Prometheus.Enabled {
			state = append(state, promIstioVs()...)
		}
		if *cnvrgApp.Spec.Monitoring.Grafana.Enabled {
			state = append(state, grafanaIstioVs()...)
		}
	}

	if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.OpenShiftIngress {
		if *cnvrgApp.Spec.Monitoring.Prometheus.Enabled {
			state = append(state, promOcpRoute()...)
		}
		if *cnvrgApp.Spec.Monitoring.Grafana.Enabled {
			state = append(state, grafanaOcpRoute()...)
		}
	}

	return state
}

func InfraMonitoringState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State

	if *infra.Spec.Monitoring.PrometheusOperator.Enabled {
		state = append(state, prometheusOperatorState()...)
	}
	if *infra.Spec.Monitoring.Prometheus.Enabled {
		state = append(state, infraPrometheusInstanceState()...)
	}
	if *infra.Spec.Monitoring.KubeStateMetrics.Enabled {
		state = append(state, kubeStateMetricsState()...)
	}
	if *infra.Spec.Monitoring.Grafana.Enabled {
		state = append(state, grafanaState()...)
	}
	if *infra.Spec.Monitoring.NodeExporter.Enabled {
		state = append(state, nodeExporterState()...)
	}
	if *infra.Spec.Monitoring.DefaultServiceMonitors.Enabled {
		state = append(state, defaultServiceMonitors()...)
	}
	if *infra.Spec.Monitoring.DcgmExporter.Enabled {
		state = append(state, dcgmExporter()...)
	}
	if *infra.Spec.SSO.Enabled {
		state = append(state, grafanaOauthProxy()...)
	}

	if infra.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
		if *infra.Spec.Monitoring.Prometheus.Enabled {
			state = append(state, promIstioVs()...)
		}
		if *infra.Spec.Monitoring.Grafana.Enabled {
			state = append(state, grafanaIstioVs()...)
		}
	}

	if infra.Spec.Networking.Ingress.IngressType == mlopsv1.OpenShiftIngress {
		if *infra.Spec.Monitoring.Prometheus.Enabled {
			state = append(state, promOcpRoute()...)
		}
		if *infra.Spec.Monitoring.Grafana.Enabled {
			state = append(state, grafanaOcpRoute()...)
		}
	}

	return state
}

func Crds() (crds []*desired.State) {
	err := pkger.Walk(path+"/crds", func(path string, info os.FileInfo, err error) error {
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
		zap.S().Error(err, "error loading prometheus crds")
	}
	return
}
