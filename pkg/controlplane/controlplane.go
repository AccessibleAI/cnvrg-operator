package controlplane

import (
	"embed"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const path = "tmpl"

//go:embed  tmpl/*
var templatesContent embed.FS

func rbacState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/conf/rbac/ccp-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/ccp-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/ccp-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/job-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/job-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/job-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/buildimage-job-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/buildimage-job-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/buildimage-job-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/spark-job-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
	}
}

func privilegedRbacState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/conf/rbac/privileged-job-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/rbac/privileged-job-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func controlPlaneConfigState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/conf/cm/config-base.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/cm/config-networking.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/cm/config-labels.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/cm/secret-base.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/conf/cm/secret-ldap.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/conf/cm/secret-object-storage.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/conf/cm/secret-smtp.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func webAppHpaState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/webapp/hpa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.HpaGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func sidekiqHpaState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/sidekiqs/sidekiq-hpa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.HpaGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func searchkiqHpaState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/sidekiqs/searchkiq-hpa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.HpaGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func systemkiqHpaState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/sidekiqs/systemkiq-hpa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.HpaGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func webAppState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/webapp/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/webapp/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/webapp/oauth.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/webapp/oauthtoken.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/webapp/pdb.tpl",
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

func webAppIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/webapp/vs.tpl",
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

func webAppOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/webapp/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func webAppIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/webapp/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func sidekiqState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/sidekiqs/sidekiq.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/sidekiqs/sidekiq-pdb.tpl",
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

func cnvrgRouter() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/router/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/router/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/router/svc.tpl",
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

func cnvrgRouterIstioVs() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/router/vs.tpl",
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

func cnvrgRouterOcpRoute() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/router/route.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.OcpRouteGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func cnvrgRouterIngress() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/router/ingress.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.IngressGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func searchkiqState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/sidekiqs/searchkiq.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/sidekiqs/searchkiq-pdb.tpl",
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

func systemkiqState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/sidekiqs/systemkiq.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/sidekiqs/systemkiq-pdb.tpl",
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

func hyperState() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/hyper/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/hyper/svc.tpl",
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

func cnvrgScheduler() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/scheduler/dep.tpl",
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

func cnvrgClusterProvisionerOperator() []*desired.State {
	return []*desired.State{
		{

			TemplatePath:   path + "/ccp/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/mgr-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/rb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/mgr-rb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{

			TemplatePath:   path + "/ccp/svc.tpl",
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

func ssoState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/webapp/oauth.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func mpiAppState() []*desired.State {
	return []*desired.State{

		{
			TemplatePath:   path + "/mpi/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/mpi/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/mpi/secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
		{
			TemplatePath:   path + "/mpi/dep.tpl",
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

func mpiInfraState() []*desired.State {
	return []*desired.State{
		{
			TemplatePath:   path + "/mpi/clusterrole.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ClusterRoleGVK],
			Own:            true,
			Updatable:      true,
			Fs:             &templatesContent,
		},
	}
}

func State(cnvrgApp *mlopsv1.CnvrgApp) []*desired.State {
	var state []*desired.State
	state = append(state, rbacState()...)

	if cnvrgApp.Spec.ControlPlane.BaseConfig.CnvrgPrivilegedJob {
		state = append(state, privilegedRbacState()...)
	}

	state = append(state, controlPlaneConfigState()...)

	if cnvrgApp.Spec.ControlPlane.WebApp.Enabled {
		state = append(state, webAppState()...)

		if cnvrgApp.Spec.ControlPlane.WebApp.Hpa.Enabled {
			state = append(state, webAppHpaState()...)
		}
		if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress {
			state = append(state, webAppIstioVs()...)
		}
		if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
			state = append(state, webAppOcpRoute()...)
		}
		if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NginxIngress {
			state = append(state, webAppIngress()...)
		}
	}

	if cnvrgApp.Spec.SSO.Enabled {
		state = append(state, ssoState()...)
	}

	if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled && cnvrgApp.Spec.ControlPlane.Sidekiq.Split {
		state = append(state, sidekiqState()...)

		if cnvrgApp.Spec.ControlPlane.Sidekiq.Hpa.Enabled {
			state = append(state, sidekiqHpaState()...)
		}
	}

	if cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled && cnvrgApp.Spec.ControlPlane.Sidekiq.Split {
		state = append(state, searchkiqState()...)

		if cnvrgApp.Spec.ControlPlane.Searchkiq.Hpa.Enabled {
			state = append(state, searchkiqHpaState()...)
		}
	}

	if cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled && cnvrgApp.Spec.ControlPlane.Sidekiq.Split {
		state = append(state, systemkiqState()...)

		if cnvrgApp.Spec.ControlPlane.Systemkiq.Hpa.Enabled {
			state = append(state, systemkiqHpaState()...)
		}
	}

	// if split stet to false -> all queues executed by sidekiq instance
	if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled && !cnvrgApp.Spec.ControlPlane.Sidekiq.Split {
		state = append(state, sidekiqState()...)
		if cnvrgApp.Spec.ControlPlane.Sidekiq.Hpa.Enabled {
			state = append(state, systemkiqHpaState()...)
		}
	}

	if cnvrgApp.Spec.ControlPlane.Hyper.Enabled {
		state = append(state, hyperState()...)
	}

	if cnvrgApp.Spec.ControlPlane.CnvrgScheduler.Enabled {
		state = append(state, cnvrgScheduler()...)
	}

	if cnvrgApp.Spec.ControlPlane.CnvrgClusterProvisionerOperator.Enabled {
		state = append(state, cnvrgClusterProvisionerOperator()...)
	}

	if cnvrgApp.Spec.ControlPlane.CnvrgRouter.Enabled {
		state = append(state, cnvrgRouter()...)
		if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress {
			state = append(state, cnvrgRouterIstioVs()...)
		}
		if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
			state = append(state, cnvrgRouterOcpRoute()...)
		}
		if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NginxIngress {
			state = append(state, cnvrgRouterIngress()...)
		}
	}

	if cnvrgApp.Spec.ControlPlane.Mpi.Enabled {
		state = append(state, mpiAppState()...)
	}
	return state
}

func MpiInfraState() []*desired.State {
	return mpiInfraState()
}

func Crds() (crds []*desired.State) {
	d, err := templatesContent.ReadDir(path + "/crds")
	if err != nil {
		zap.S().Error(err, "error loading control plane crds")
	}
	for _, f := range d {
		crd := &desired.State{

			TemplatePath:   path + "/crds/" + f.Name(),
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.CrdGVK],
			Own:            false,
			Updatable:      false,
			Fs:             &templatesContent,
		}
		crds = append(crds, crd)
	}

	return

}
