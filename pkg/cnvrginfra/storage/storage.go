package storage

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/storage/tmpl"

var hostPathState = []*desired.State{
	{
		TemplatePath:   path + "/hostpath/class.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.StorageClassGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/hostpath/clusterrole.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/hostpath/clusterrolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/hostpath/daemonset.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DaemonSetGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/hostpath/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
		Own:            true,
	},
}

var nfsClientState = []*desired.State{
	{
		
		TemplatePath:   path + "/nfsclient/class.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.StorageClassGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/nfsclient/clusterrole.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/nfsclient/clusterrolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/nfsclient/rolebinding.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleBindingGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/nfsclient/sa.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.SaGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/nfsclient/role.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.RoleGVR],
		Own:            true,
	},
	{
		
		TemplatePath:   path + "/nfsclient/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            desired.Kinds[desired.DeploymentGVR],
		Own:            true,
	},
}

func State(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if cnvrgInfra.Spec.Storage.Enabled == "true" && cnvrgInfra.Spec.Storage.Hostpath.Enabled == "true" {
		state = append(state, hostPathState...)
	}
	if cnvrgInfra.Spec.Storage.Enabled == "true" && cnvrgInfra.Spec.Storage.Nfs.Enabled == "true" {
		state = append(state, nfsClientState...)
	}
	return state
}
