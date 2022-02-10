package controlplane

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
)

const path = "/pkg/controlplane/tmpl"

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
		},
		{

			TemplatePath:   path + "/conf/rbac/ccp-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/rbac/ccp-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/conf/rbac/job-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/rbac/job-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/rbac/job-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/conf/rbac/buildimage-job-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/rbac/buildimage-job-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/rbac/buildimage-job-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/conf/rbac/spark-job-sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      false,
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
		},
		{

			TemplatePath:   path + "/conf/rbac/privileged-job-rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{

			TemplatePath:   path + "/conf/cm/config-networking.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/cm/config-labels.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/conf/cm/secret-base.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
		},
		{

			TemplatePath:   path + "/conf/cm/secret-ldap.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/conf/cm/secret-object-storage.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/conf/cm/secret-smtp.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{
			TemplatePath:   path + "/webapp/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/webapp/oauth.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/webapp/oauthtoken.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      false,
		},
		{
			TemplatePath:   path + "/webapp/pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{
			TemplatePath:   path + "/sidekiqs/sidekiq-pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{
			TemplatePath:   path + "/router/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/router/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{
			TemplatePath:   path + "/sidekiqs/searchkiq-pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{
			TemplatePath:   path + "/sidekiqs/systemkiq-pdb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.PodDisruptionBudgetGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{

			TemplatePath:   path + "/hyper/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{

			TemplatePath:   path + "/ccp/sa.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SaGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/ccp/role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/ccp/mgr-role.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/ccp/rb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/ccp/mgr-rb.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/ccp/cm.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.ConfigMapGVK],
			Own:            true,
			Updatable:      true,
		},
		{

			TemplatePath:   path + "/ccp/svc.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SvcGVK],
			Own:            true,
			Updatable:      true,
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
		},
		{
			TemplatePath:   path + "/mpi/rolebinding.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.RoleBindingGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/mpi/secret.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.SecretGVK],
			Own:            true,
			Updatable:      true,
		},
		{
			TemplatePath:   path + "/mpi/dep.tpl",
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.DeploymentGVK],
			Own:            true,
			Updatable:      true,
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
	err := pkger.Walk(path+"/crds", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		updatable := true
		// mpi crds can't be updatable, b/c of the crd version (v1beta1)
		// do not change the mpi crd file name!
		if info.Name() == "mpijobs.yaml" {
			updatable = false
		}
		crd := &desired.State{
			TemplatePath:   path,
			Template:       nil,
			ParsedTemplate: "",
			Obj:            &unstructured.Unstructured{},
			GVK:            desired.Kinds[desired.CrdGVK],
			Own:            false,
			Updatable:      updatable,
		}
		crds = append(crds, crd)
		return nil
	})
	if err != nil {
		zap.S().Error(err, "error loading control plane crds")
	}
	return
}
