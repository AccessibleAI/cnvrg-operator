package desired

import "k8s.io/apimachinery/pkg/runtime/schema"

type GVKName string

const (
	DeploymentGVK           GVKName = "DeploymentGVK"
	StatefulSetGVK          GVKName = "StatefulSetGVK"
	DaemonSetGVK            GVKName = "DaemonSetGVK"
	ConfigMapGVK            GVKName = "ConfigMapGVK"
	PvcGVK                  GVKName = "PvcGVK"
	SecretGVK               GVKName = "SecretGVK"
	SvcGVK                  GVKName = "SvcGVK"
	SaGVK                   GVKName = "SaGVK"
	CrdGVK                  GVKName = "CrdGVK"
	PrometheusGVK           GVKName = "PrometheusGVK"
	ServiceMonitorGVK       GVKName = "ServiceMonitorGVK"
	PrometheusRuleGVK       GVKName = "PrometheusRuleGVK"
	ClusterRoleGVK          GVKName = "ClusterRoleGVK"
	ClusterRoleBindingGVK   GVKName = "ClusterRoleBindingGVK"
	RoleGVK                 GVKName = "RoleGVK"
	RoleBindingGVK          GVKName = "RoleBindingGVK"
	OcpRouteGVK             GVKName = "OcpRouteGVK"
	IngressGVK              GVKName = "IngressGVK"
	IstioVsGVK              GVKName = "IstioVsGVK"
	IstioGVK                GVKName = "IstioGVK"
	IstioDestinationRuleGVK GVKName = "IstioDestinationRule"
	IstioGwGVK              GVKName = "IstioGwGVK"
	StorageClassGVK         GVKName = "StorageClassGVK"
	PodDisruptionBudgetGVK  GVKName = "PodDisruptionBudget"
	HpaGVK                  GVKName = "HpaGVK"
	PriorityClassGVK        GVKName = "PriorityClassGVK"
	JobGVK                  GVKName = "JobGvk"
)

var Kinds = map[GVKName]schema.GroupVersionKind{
	DeploymentGVK: schema.GroupVersionKind{
		Kind:    "Deployment",
		Group:   "apps",
		Version: "v1",
	},

	StatefulSetGVK: schema.GroupVersionKind{
		Kind:    "StatefulSet",
		Group:   "apps",
		Version: "v1",
	},

	DaemonSetGVK: schema.GroupVersionKind{
		Kind:    "DaemonSet",
		Group:   "apps",
		Version: "v1",
	},

	PvcGVK: schema.GroupVersionKind{
		Kind:    "PersistentVolumeClaim",
		Group:   "",
		Version: "v1",
	},

	SecretGVK: schema.GroupVersionKind{
		Kind:    "Secret",
		Group:   "",
		Version: "v1",
	},

	SvcGVK: schema.GroupVersionKind{
		Kind:    "Service",
		Group:   "",
		Version: "v1",
	},

	SaGVK: schema.GroupVersionKind{
		Kind:    "ServiceAccount",
		Group:   "",
		Version: "v1",
	},

	CrdGVK: schema.GroupVersionKind{
		Kind:    "CustomResourceDefinition",
		Group:   "apiextensions.k8s.io",
		Version: "v1",
	},

	IstioGVK: schema.GroupVersionKind{
		Kind:    "IstioOperator",
		Group:   "install.istio.io",
		Version: "v1alpha1",
	},
	IstioVsGVK: schema.GroupVersionKind{
		Kind:    "VirtualService",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},

	IstioDestinationRuleGVK: schema.GroupVersionKind{
		Kind:    "DestinationRule",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},

	IstioGwGVK: schema.GroupVersionKind{
		Kind:    "Gateway",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},

	ClusterRoleGVK: schema.GroupVersionKind{
		Kind:    "ClusterRole",
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
	},

	ClusterRoleBindingGVK: schema.GroupVersionKind{
		Kind:    "ClusterRoleBinding",
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
	},

	RoleGVK: schema.GroupVersionKind{
		Kind:    "Role",
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
	},

	OcpRouteGVK: schema.GroupVersionKind{
		Kind:    "Route",
		Group:   "route.openshift.io",
		Version: "v1",
	},

	RoleBindingGVK: schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "RoleBinding",
	},
	ConfigMapGVK: schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	},
	StorageClassGVK: schema.GroupVersionKind{
		Group:   "storage.k8s.io",
		Version: "v1",
		Kind:    "StorageClass",
	},
	PrometheusGVK: schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "Prometheus",
	},
	ServiceMonitorGVK: schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "ServiceMonitor",
	},
	PrometheusRuleGVK: schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "PrometheusRule",
	},
	IngressGVK: schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "Ingress",
	},
	PodDisruptionBudgetGVK: schema.GroupVersionKind{
		Group:   "policy",
		Version: "v1beta1",
		Kind:    "PodDisruptionBudget",
	},
	HpaGVK: schema.GroupVersionKind{
		Group:   "autoscaling",
		Version: "v2beta2",
		Kind:    "HorizontalPodAutoscaler",
	},
	PriorityClassGVK: schema.GroupVersionKind{
		Group:   "scheduling.k8s.io",
		Version: "v1",
		Kind:    "PriorityClass",
	},
	JobGVK: schema.GroupVersionKind{
		Kind:    "Job",
		Group:   "batch",
		Version: "v1",
	},
}
