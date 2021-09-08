package storage

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/storage/tmpl"

func hostPathState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/hostpath/class.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.StorageClassGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/hostpath/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/hostpath/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/hostpath/daemonset.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DaemonSetGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/hostpath/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            false,
			Updatable:      false,
		},
	}
}

func nfsClientState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/nfsclient/class.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.StorageClassGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/nfsclient/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/nfsclient/clusterrolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.ClusterRoleBindingGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/nfsclient/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/nfsclient/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/nfsclient/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/nfsclient/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            false,
			Updatable:      false,
		},
	}
}

func State(cnvrgInfra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if cnvrgInfra.Spec.Storage.Hostpath.Enabled {
		state = append(state, hostPathState()...)
	}
	if cnvrgInfra.Spec.Storage.Nfs.Enabled {
		state = append(state, nfsClientState()...)
	}
	return state
}
