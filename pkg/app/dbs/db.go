package dbs

import (
	"embed"
)

const fsRoot = "tmpl"

//go:embed  tmpl/*
var fs embed.FS

//func EsCreds(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/elastic/secret.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			TemplateData:   data,
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func PromDBCreds(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/prom/creds.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			TemplateData:   data,
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func prometheusState(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/prom/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/cm.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.ConfigMapGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/pvc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PvcGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/route.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.OcpRouteGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/prom/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			TemplateData:   data,
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func esState() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/elastic/cm.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.ConfigMapGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/elastic/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/elastic/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/elastic/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/elastic/sts.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.StatefulSetGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/elastic/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func PgCreds(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/pg/secret.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			TemplateData:   data,
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func pgState() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/pg/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/pg/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/pg/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/pg/pvc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PvcGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/pg/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/pg/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/pg/pdb.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func RedisCreds(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/redis/secret.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			TemplateData:   data,
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func redisState() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/redis/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/redis/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/redis/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/redis/pvc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PvcGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/redis/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/redis/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/redis/pdb.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func singleBackendMinio() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/minio/pdb.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/minio/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/minio/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/minio/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/pvc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PvcGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func sharedBackendMinio() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/minio/pdb.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/minio/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/pvc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PvcGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/sh-dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//
//			TemplatePath:   path + "/minio/dr.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IstioDestinationRuleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func esIstioVs() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/elastic/vs.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IstioVsGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func esIngress() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/elastic/ingress.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IngressGVK],
//			Own:            true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func esOcpRoute() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/elastic/route.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.OcpRouteGVK],
//			Own:            true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func minioIstioVs() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/minio/vs.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IstioVsGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func minioOcpRoute() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/minio/route.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.OcpRouteGVK],
//			Own:            true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func minioIngress() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/minio/ingress.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IngressGVK],
//			Own:            true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func cvatState() []*desired.State {
//	return []*desired.State{
//		// cvat pg
//		{
//			TemplatePath:   path + "/cvat-pg/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-pg/pvc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.PvcGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-pg/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-pg/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-pg/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-pg/secret.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-pg/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//
//		// cvat redis
//		{
//			TemplatePath:   path + "/cvat-redis/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-redis/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-redis/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-redis/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-redis/secret.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/cvat-redis/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func KibanaConfSecret(data desired.TemplateData) []*desired.State {
//	return []*desired.State{
//		{
//			TemplateData:   data,
//			TemplatePath:   path + "/es/kibana/secret.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func kibanaSvc(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/kibana/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			TemplateData:   data,
//		},
//	}
//}
//
//func kibanaIstioVs(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/kibana/vs.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IstioVsGVK],
//			Own:            true,
//			Updatable:      true,
//			TemplateData:   data,
//		},
//	}
//}
//
//func kibanaOcpRoute() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/kibana/route.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.OcpRouteGVK],
//			Own:            true,
//			Updatable:      true,
//		},
//	}
//}
//
//func kibanaIngress() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/kibana/ingress.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IngressGVK],
//			Own:            true,
//			Updatable:      true,
//		},
//	}
//}
//
//func kibana() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/kibana/sa.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SaGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/kibana/role.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/kibana/rolebinding.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.RoleBindingGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/kibana/dep.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.DeploymentGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//		{
//			TemplatePath:   path + "/es/kibana/credsec.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func ElastCreds(data *desired.TemplateData) []*desired.State {
//	return []*desired.State{
//		{
//			TemplateData:   data,
//			TemplatePath:   path + "/elastalert/credsec.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			Own:            true,
//			Updatable:      false,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func elastAlertSvc(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/elastalert/svc.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SvcGVK],
//			Own:            true,
//			Updatable:      true,
//			TemplateData:   data,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func elastAlertIstioVs(data interface{}) []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/elastalert/vs.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IstioVsGVK],
//			Own:            true,
//			Updatable:      true,
//			TemplateData:   data,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func elastAlertOcpRoute() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/elastalert/route.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.OcpRouteGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func elastAlertIngress() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/elastalert/ingress.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.IngressGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func ElastAlert() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath: path + "/elastalert/sa.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.SaGVK],
//			Own:          true,
//			Updatable:    false,
//			Fs:           &templatesContent,
//		},
//		{
//			TemplatePath: path + "/elastalert/authproxycm.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.ConfigMapGVK],
//			Own:          true,
//			Updatable:    true,
//			Fs:           &templatesContent,
//		},
//		{
//			TemplatePath: path + "/elastalert/role.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.RoleGVK],
//			Own:          true,
//			Updatable:    true,
//			Fs:           &templatesContent,
//		},
//		{
//			TemplatePath: path + "/elastalert/rolebinding.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.RoleBindingGVK],
//			Own:          true,
//			Updatable:    true,
//			Fs:           &templatesContent,
//		},
//		{
//			TemplatePath: path + "/elastalert/pvc.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.PvcGVK],
//			Own:          true,
//			Updatable:    false,
//			Fs:           &templatesContent,
//		},
//		{
//
//			TemplatePath: path + "/elastalert/cm.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.ConfigMapGVK],
//			Own:          true,
//			Updatable:    true,
//			Fs:           &templatesContent,
//		},
//		{
//
//			TemplatePath: path + "/elastalert/dep.tpl",
//			Obj:          &unstructured.Unstructured{},
//			GVK:          desired.Kinds[desired.DeploymentGVK],
//			Own:          true,
//			Updatable:    true,
//			Fs:           &templatesContent,
//		},
//	}
//}
//
//func kibanaOauthProxy() []*desired.State {
//	return []*desired.State{
//		{
//			TemplatePath:   path + "/es/kibana/oauth.tpl",
//			Template:       nil,
//			ParsedTemplate: "",
//			Obj:            &unstructured.Unstructured{},
//			GVK:            desired.Kinds[desired.SecretGVK],
//			Own:            true,
//			Updatable:      true,
//			Fs:             &templatesContent,
//		},
//	}
//}
//
//func AppDbsState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
//	var state []*desired.State
//
//	// cvat pg and redis
//	if cnvrgApp.Spec.Dbs.Cvat.Enabled {
//		state = append(state, cvatState()...)
//	}
//
//	// pg
//	if cnvrgApp.Spec.Dbs.Pg.Enabled {
//		state = append(state, pgState()...)
//	}
//
//	// redis
//	if cnvrgApp.Spec.Dbs.Redis.Enabled {
//		state = append(state, redisState()...)
//	}
//
//	// minio
//	if cnvrgApp.Spec.Dbs.Minio.Enabled && cnvrgApp.Spec.Dbs.Minio.SharedStorage.Enabled {
//		state = append(state, sharedBackendMinio()...)
//	} else if cnvrgApp.Spec.Dbs.Minio.Enabled {
//		state = append(state, singleBackendMinio()...)
//	}
//	if cnvrgApp.Spec.Dbs.Minio.Enabled {
//		switch cnvrgApp.Spec.Networking.Ingress.Type {
//		case mlopsv1.IstioIngress:
//			state = append(state, minioIstioVs()...)
//		case mlopsv1.NginxIngress:
//			state = append(state, minioIngress()...)
//		case mlopsv1.OpenShiftIngress:
//			state = append(state, minioOcpRoute()...)
//		}
//	}
//
//	// elasticsearch
//	if cnvrgApp.Spec.Dbs.Es.Enabled {
//		state = append(state, esState()...)
//		switch cnvrgApp.Spec.Networking.Ingress.Type {
//		case mlopsv1.IstioIngress:
//			state = append(state, esIstioVs()...)
//		case mlopsv1.NginxIngress:
//			state = append(state, esIngress()...)
//		case mlopsv1.OpenShiftIngress:
//			state = append(state, esOcpRoute()...)
//		}
//	}
//
//	return state
//}
//
//func ApplyAppPrometheus(app *mlopsv1.CnvrgApp, data interface{}) []*desired.State {
//	if app.Spec.Dbs.Prom.Enabled {
//		return prometheusState(data)
//	}
//	return nil
//}
//
//func KibanaState(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
//	var state []*desired.State
//
//	state = append(state, kibana()...)
//	state = append(state, kibanaSvc(nil)...)
//
//	if cnvrgApp.Spec.SSO.Enabled {
//		state = append(state, kibanaOauthProxy()...)
//	}
//
//	switch cnvrgApp.Spec.Networking.Ingress.Type {
//	case mlopsv1.IstioIngress:
//		state = append(state, kibanaIstioVs(nil)...)
//	case mlopsv1.NginxIngress:
//		state = append(state, kibanaIngress()...)
//	case mlopsv1.OpenShiftIngress:
//		state = append(state, kibanaOcpRoute()...)
//	}
//
//	return state
//}
