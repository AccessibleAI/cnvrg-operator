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
}

func State(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {

	var state []*desired.State

	if cnvrgInfra.Spec.Monitoring.Enabled == "true" && cnvrgInfra.Spec.Monitoring.PrometheusOperator.Enabled == "true" {
		state = append(state, prometheusOperatorState...)
	}

	//if cnvrgInfra.Spec.Monitoring.Enabled == "true" && cnvrgInfra.Spec.Monitoring.Prometheus.Enabled == "true" {
	//	state = append(state, prometheusInstanceState...)
	//}

	//if cnvrgInfra.Spec.Monitoring.Enabled == "true" && cnvrgInfra.Spec.Monitoring.KubeletServiceMonitor == "true" {
	//	state = append(state, kubeletServiceMonitorInstanceState...)
	//}

	if cnvrgInfra.Spec.Monitoring.Enabled == "true" && cnvrgInfra.Spec.Monitoring.KubeStateMetrics.Enabled == "true" {
		state = append(state, kubeStateMetricsState...)
	}

	if cnvrgInfra.Spec.Monitoring.Enabled == "true" {
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
