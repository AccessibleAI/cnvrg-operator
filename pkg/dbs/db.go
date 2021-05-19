package dbs

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "/pkg/dbs/tmpl"

func esState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/es/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/es/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/es/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/es/sts.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.StatefulSetGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/es/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func pgState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/pg/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/pg/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/pg/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/pg/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/pg/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/pg/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/pg/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PodDisruptionBudgetGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func redisState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/redis/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/redis/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/redis/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/redis/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/redis/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/redis/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/redis/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PodDisruptionBudgetGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func singleBackendMinio() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PodDisruptionBudgetGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/minio/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/minio/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/minio/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/minio/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/minio/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/minio/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func sharedBackendMinio() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PodDisruptionBudgetGVR],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/minio/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/minio/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            false,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/minio/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/minio/sh-dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/minio/dr.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioDestinationRuleGVR],
			Own:            true,
			Updatable:      true,
		},
	}
}

func esIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/es/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioVsGVR],
			Own:            true,
		},
	}
}

func esIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/es/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IngressGVR],
			Own:            true,
		},
	}
}

func esOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/es/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.OcpRouteGVR],
			Own:            true,
		},
	}
}

func minioIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioVsGVR],
			Own:            true,
		},
	}
}

func minioOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.OcpRouteGVR],
			Own:            true,
		},
	}
}

func minioIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IngressGVR],
			Own:            true,
		},
	}
}

func AppDbsState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	// pg
	if *cnvrgApp.Spec.Dbs.Pg.Enabled {
		state = append(state, pgState()...)
	}

	// redis
	if *cnvrgApp.Spec.Dbs.Redis.Enabled {
		state = append(state, redisState()...)
	}

	// minio
	if *cnvrgApp.Spec.Dbs.Minio.Enabled && *cnvrgApp.Spec.Dbs.Minio.SharedStorage.Enabled {
		state = append(state, sharedBackendMinio()...)
	} else if *cnvrgApp.Spec.Dbs.Minio.Enabled {
		state = append(state, singleBackendMinio()...)
	}
	if *cnvrgApp.Spec.Dbs.Minio.Enabled {
		switch cnvrgApp.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, minioIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, minioIngress()...)
		case mlopsv1.OpenShiftIngress:
			state = append(state, minioOcpRoute()...)
		}
	}

	// elasticsearch
	if *cnvrgApp.Spec.Dbs.Es.Enabled {
		state = append(state, esState()...)
		switch cnvrgApp.Spec.Networking.Ingress.Type {
		case mlopsv1.IstioIngress:
			state = append(state, esIstioVs()...)
		case mlopsv1.NginxIngress:
			state = append(state, esIngress()...)
		case mlopsv1.OpenShiftIngress:
			state = append(state, esOcpRoute()...)
		}
	}

	return state
}

func InfraDbsState(infra *mlopsv1.CnvrgInfra) []*desired.State {
	var state []*desired.State
	if *infra.Spec.Dbs.Redis.Enabled {
		state = append(state, redisState()...)
	}
	return state
}
