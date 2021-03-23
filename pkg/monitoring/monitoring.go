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

var prometheusOperatorState = []*desired.State{
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

var infraPrometheusInstanceState = []*desired.State{
	{
		TemplatePath:   path + "/prometheus/cluster-res/kubelet.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ServiceMonitorGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/cluster-res/clusterrole.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/cluster-res/clusterrolebinding.tpl",
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
		TemplatePath:   path + "/prometheus/instance/prometheus.tpl",
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
		TemplatePath:   path + "/prometheus/instance/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/instance/rules.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PrometheusRuleGVR],
		Own:            true,
	},
}

var kubeStateMetricsState = []*desired.State{
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

var grafanaState = []*desired.State{
	{
		TemplatePath:   path + "/grafana/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
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
		TemplatePath:   path + "/grafana/datasource.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
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
	{
		TemplatePath:   path + "/grafana/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
}

var nodeExporterState = []*desired.State{
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

var appPrometheusInstance = []*desired.State{
	{
		TemplatePath:   path + "/prometheus/instance/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/instance/prometheus.tpl",
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
		TemplatePath:   path + "/prometheus/instance/vs.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.IstioVsGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/namespace-res/role.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/namespace-res/rolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleBindingGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/namespace-res/staticconfig.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SecretGVR],
		Own:            true,
	},
	{
		TemplatePath:   path + "/prometheus/instance/rules.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PrometheusRuleGVR],
		Own:            true,
	},
}

func AppMonitoringState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {

	var state []*desired.State
	if cnvrgApp.Spec.Monitoring.Enabled == "true" && cnvrgApp.Spec.Monitoring.Prometheus.Enabled == "true" {
		state = append(state, appPrometheusInstance...)
	}
	if cnvrgApp.Spec.Monitoring.Enabled == "true" && cnvrgApp.Spec.Monitoring.Grafana.Enabled == "true" {
		state = append(state, grafanaState...)
	}

	return state
}

func InfraMonitoringState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if infra.Spec.Monitoring.Enabled == "true" && infra.Spec.Monitoring.PrometheusOperator.Enabled == "true" {
		state = append(state, prometheusOperatorState...)
	}
	if infra.Spec.Monitoring.Enabled == "true" && infra.Spec.Monitoring.Prometheus.Enabled == "true" {
		state = append(state, infraPrometheusInstanceState...)
	}
	if infra.Spec.Monitoring.Enabled == "true" && infra.Spec.Monitoring.KubeStateMetrics.Enabled == "true" {
		state = append(state, kubeStateMetricsState...)
	}
	if infra.Spec.Monitoring.Enabled == "true" && infra.Spec.Monitoring.Grafana.Enabled == "true" {
		state = append(state, grafanaState...)
	}
	if infra.Spec.Monitoring.Enabled == "true" && infra.Spec.Monitoring.NodeExporter.Enabled == "true" {
		state = append(state, nodeExporterState...)
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
