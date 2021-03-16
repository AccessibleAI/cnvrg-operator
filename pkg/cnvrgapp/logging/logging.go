package logging

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/logging/tmpl"

var es = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/es/sts.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.StatefulSetGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/es/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
}

var fluentd = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/fluentd/clusterrole.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/fluentd/clusterrolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/fluentd/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/fluentd/daemonset.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DaemonSetGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/fluentd/cm.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
}

var elastAlert = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/elastalert/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.PvcGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/elastalert/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/elastalert/cm.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ConfigMapGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/elastalert/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

var kibana = []*desired.State{
	{
		Name:           "",
		TemplatePath:   path + "/kibana/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
	{
		Name:           "",
		TemplatePath:   path + "/kibana/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SvcGVR],
		Own:            true,
	},
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Es.Enabled == "true" {
		state = append(state, es...)
	}
	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Fluentd.Enabled == "true" {
		state = append(state, fluentd...)
	}
	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Elastalert.Enabled == "true" {
		state = append(state, elastAlert...)
	}
	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Kibana.Enabled == "true" {
		state = append(state, kibana...)
	}

	return state
}
