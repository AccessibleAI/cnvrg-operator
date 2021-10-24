package desired

import "k8s.io/apimachinery/pkg/runtime/schema"

type GVRName string

const (
	DeploymentGVR           GVRName = "DeploymentGVR"
	StatefulSetGVR          GVRName = "StatefulSetGVR"
	DaemonSetGVR            GVRName = "DaemonSetGVR"
	ConfigMapGVR            GVRName = "ConfigMapGVR"
	PvcGVR                  GVRName = "PvcGVR"
	SecretGVR               GVRName = "SecretGVR"
	SvcGVR                  GVRName = "SvcGVR"
	SaGVR                   GVRName = "SaGVR"
	CrdGVR                  GVRName = "CrdGVR"
	PrometheusGVR           GVRName = "PrometheusGVR"
	ServiceMonitorGVR       GVRName = "ServiceMonitorGVR"
	PrometheusRuleGVR       GVRName = "PrometheusRuleGVR"
	ClusterRoleGVR          GVRName = "ClusterRoleGVR"
	ClusterRoleBindingGVR   GVRName = "ClusterRoleBindingGVR"
	RoleGVR                 GVRName = "RoleGVR"
	RoleBindingGVR          GVRName = "RoleBindingGVR"
	OcpRouteGVR             GVRName = "OcpRouteGVR"
	IngressGVR              GVRName = "IngressGVR"
	IstioVsGVR              GVRName = "IstioVsGVR"
	IstioGVR                GVRName = "IstioGVR"
	IstioDestinationRuleGVR GVRName = "IstioDestinationRule"
	IstioGwGVR              GVRName = "IstioGwGVR"
	StorageClassGVR         GVRName = "StorageClassGVR"
	PodDisruptionBudgetGVR  GVRName = "PodDisruptionBudget"
	HpaGVR                  GVRName = "HpaGVR"
	PriorityClassGCR        GVRName = "PriorityClassGCR"
)

var Kinds = map[GVRName]schema.GroupVersionKind{
	DeploymentGVR: schema.GroupVersionKind{
		Kind:    "Deployment",
		Group:   "apps",
		Version: "v1",
	},

	StatefulSetGVR: schema.GroupVersionKind{
		Kind:    "StatefulSet",
		Group:   "apps",
		Version: "v1",
	},

	DaemonSetGVR: schema.GroupVersionKind{
		Kind:    "DaemonSet",
		Group:   "apps",
		Version: "v1",
	},

	PvcGVR: schema.GroupVersionKind{
		Kind:    "PersistentVolumeClaim",
		Group:   "",
		Version: "v1",
	},

	SecretGVR: schema.GroupVersionKind{
		Kind:    "Secret",
		Group:   "",
		Version: "v1",
	},

	SvcGVR: schema.GroupVersionKind{
		Kind:    "Service",
		Group:   "",
		Version: "v1",
	},

	SaGVR: schema.GroupVersionKind{
		Kind:    "ServiceAccount",
		Group:   "",
		Version: "v1",
	},

	CrdGVR: schema.GroupVersionKind{
		Kind:    "CustomResourceDefinition",
		Group:   "apiextensions.k8s.io",
		Version: "v1",
	},

	IstioGVR: schema.GroupVersionKind{
		Kind:    "IstioOperator",
		Group:   "install.istio.io",
		Version: "v1alpha1",
	},
	IstioVsGVR: schema.GroupVersionKind{
		Kind:    "VirtualService",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},

	IstioDestinationRuleGVR: schema.GroupVersionKind{
		Kind:    "DestinationRule",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},

	IstioGwGVR: schema.GroupVersionKind{
		Kind:    "Gateway",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},

	ClusterRoleGVR: schema.GroupVersionKind{
		Kind:    "ClusterRole",
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
	},

	ClusterRoleBindingGVR: schema.GroupVersionKind{
		Kind:    "ClusterRoleBinding",
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
	},

	RoleGVR: schema.GroupVersionKind{
		Kind:    "Role",
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
	},

	OcpRouteGVR: schema.GroupVersionKind{
		Kind:    "Route",
		Group:   "route.openshift.io",
		Version: "v1",
	},

	RoleBindingGVR: schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "RoleBinding",
	},
	ConfigMapGVR: schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	},
	StorageClassGVR: schema.GroupVersionKind{
		Group:   "storage.k8s.io",
		Version: "v1",
		Kind:    "StorageClass",
	},
	PrometheusGVR: schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "Prometheus",
	},
	ServiceMonitorGVR: schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "ServiceMonitor",
	},
	PrometheusRuleGVR: schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "PrometheusRule",
	},
	IngressGVR: schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1beta1",
		Kind:    "Ingress",
	},
	PodDisruptionBudgetGVR: schema.GroupVersionKind{
		Group:   "policy",
		Version: "v1beta1",
		Kind:    "PodDisruptionBudget",
	},
	HpaGVR: schema.GroupVersionKind{
		Group:   "autoscaling",
		Version: "v2beta2",
		Kind:    "HorizontalPodAutoscaler",
	},
	PriorityClassGCR: schema.GroupVersionKind{
		Group:   "scheduling.k8s.io",
		Version: "v1",
		Kind:    "PriorityClass",
	},
}
