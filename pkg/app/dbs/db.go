package dbs

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS

func EsCreds(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/es/secret.tpl",
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

func PromDBCreds(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prom/creds.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			TemplateData:   data,
			Own:            true,
			Updatable:      false,
		},
	}
}

func prometheusState(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/prom/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/prom/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/prom/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/prom/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/prom/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/prom/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/prom/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/prom/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			TemplateData:   data,
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func esState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/es/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/es/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/es/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/es/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/es/sts.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.StatefulSetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/es/svc.tpl",
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

func PgCreds(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/pg/secret.tpl",
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

func pgState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/pg/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/pg/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/pg/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/pg/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/pg/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/pg/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/pg/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func RedisCreds(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/redis/secret.tpl",
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

func redisState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/redis/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/redis/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/redis/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/redis/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/redis/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/redis/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/redis/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/minio/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/minio/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/minio/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/minio/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/sh-dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/minio/dr.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioDestinationRuleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Fs:             &templatesContent,
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
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Fs:             &templatesContent,
		},
	}
}

func cvatState() []*desired.State {
	return []*desired.State{
		// cvat pg
		{
			TemplatePath:   path + "/cvat-pg/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-pg/pvc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PvcGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-pg/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-pg/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-pg/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-pg/secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-pg/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},

		// cvat redis
		{
			TemplatePath:   path + "/cvat-redis/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-redis/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-redis/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-redis/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-redis/secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/cvat-redis/svc.tpl",
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

func kibanaSvc(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			TemplateData:   data,
		},
	}
}

func kibanaIstioVs(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/kibana/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
			TemplateData:   data,
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

func elastAlertSvc(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/elastalert/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			TemplateData:   data,
		},
	}
}

func elastAlertIstioVs(data interface{}) []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/elastalert/vs.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IstioVsGVK],
			Own:            true,
			Updatable:      true,
			TemplateData:   data,
		},
	}
}

func elastAlertOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/elastalert/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Updatable:      true,
		},
	}
}

func elastAlertIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/elastalert/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Updatable:      true,
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

func AppDbsState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	// cvat pg and redis
	if cnvrgApp.Spec.Dbs.Cvat.Enabled {
		state = append(state, cvatState()...)
	}

	// pg
	if cnvrgApp.Spec.Dbs.Pg.Enabled {
		state = append(state, pgState()...)
	}

	// redis
	if cnvrgApp.Spec.Dbs.Redis.Enabled {
		state = append(state, redisState()...)
	}

	// minio
	if cnvrgApp.Spec.Dbs.Minio.Enabled && cnvrgApp.Spec.Dbs.Minio.SharedStorage.Enabled {
		state = append(state, sharedBackendMinio()...)
	} else if cnvrgApp.Spec.Dbs.Minio.Enabled {
		state = append(state, singleBackendMinio()...)
	}
	if cnvrgApp.Spec.Dbs.Minio.Enabled {
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
	if cnvrgApp.Spec.Dbs.Es.Enabled {
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

func ApplyAppPrometheus(app *mlopsv1.CnvrgApp, data interface{}) []*desired.State {
	if app.Spec.Dbs.Prom.Enabled {
		return prometheusState(data)
	}
	return nil
}

func KibanaState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State

	state = append(state, kibana()...)
	state = append(state, kibanaSvc(nil)...)

	if cnvrgApp.Spec.SSO.Enabled {
		state = append(state, kibanaOauthProxy()...)
	}

	switch cnvrgApp.Spec.Networking.Ingress.Type {
	case mlopsv1.IstioIngress:
		state = append(state, kibanaIstioVs(nil)...)
	case mlopsv1.NginxIngress:
		state = append(state, kibanaIngress()...)
	case mlopsv1.OpenShiftIngress:
		state = append(state, kibanaOcpRoute()...)
	}

	return state
}
