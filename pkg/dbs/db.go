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
		},
		{
			TemplatePath:   path + "/es/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/es/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/es/sts.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.StatefulSetGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/es/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
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
		},
		{
			TemplatePath:   path + "/pg/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/pg/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/pg/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/pg/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/pg/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/pg/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
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
		},
		{
			TemplatePath:   path + "/redis/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/redis/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/redis/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/redis/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/redis/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
	}
}

func singleBackendMinio() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/minio/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleGVR],
			Own:            true,
		},
		{
			TemplatePath:   path + "/minio/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.RoleBindingGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/minio/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/minio/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/minio/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
	}
}

func sharedBackendMinio() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/minio/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SaGVR],
			Own:            true,
		},

		{

			TemplatePath:   path + "/minio/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.PvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/minio/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.SvcGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/minio/sh-dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.DeploymentGVR],
			Own:            true,
		},
		{

			TemplatePath:   path + "/minio/dr.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVR:            desired.Kinds[desired.IstioDestinationRuleGVR],
			Own:            true,
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

func AppDbsState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	if *cnvrgApp.Spec.Dbs.Pg.Enabled {
		state = append(state, pgState()...)
	}

	if *cnvrgApp.Spec.Dbs.Redis.Enabled {
		state = append(state, redisState()...)
	}

	if *cnvrgApp.Spec.Dbs.Minio.Enabled && *cnvrgApp.Spec.Dbs.Minio.SharedStorage.Enabled {
		state = append(state, sharedBackendMinio()...)
	} else {
		state = append(state, singleBackendMinio()...)
	}
	if *cnvrgApp.Spec.Dbs.Es.Enabled {
		state = append(state, esState()...)
	}

	if *cnvrgApp.Spec.Dbs.Es.Enabled && cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
		state = append(state, esIstioVs()...)
	}

	if *cnvrgApp.Spec.Dbs.Es.Enabled && cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.OpenShiftIngress {
		state = append(state, esOcpRoute()...)
	}

	if *cnvrgApp.Spec.Dbs.Minio.Enabled && cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
		state = append(state, minioIstioVs()...)
	}

	if *cnvrgApp.Spec.Dbs.Minio.Enabled && cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.OpenShiftIngress {
		state = append(state, minioOcpRoute()...)
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
