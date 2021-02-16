package desired

import "k8s.io/apimachinery/pkg/runtime/schema"

type GVRName string

const (
	DeploymentGVR         GVRName = "DeploymentGVR"
	PvcGVR                GVRName = "PvcGVR"
	SecretGVR             GVRName = "SecretGVR"
	SvcGVR                GVRName = "SvcGVR"
	SaGVR                 GVRName = "SaGVR"
	CrdGVR                GVRName = "CrdGVR"
	IstioGVR              GVRName = "IstioGVR"
	ClusterRoleGVR        GVRName = "ClusterRoleGVR"
	ClusterRoleBindingGVR GVRName = "ClusterRoleBindingGVR"
	RoleGVR               GVRName = "RoleGVR"
	OcpRouteGVR           GVRName = "OcpRouteGVR"
	IstioVsGVR            GVRName = "IstioVsGVR"
)

var Kinds = map[GVRName]schema.GroupVersionKind{
	DeploymentGVR: schema.GroupVersionKind{
		Kind:    "Deployment",
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

	IstioVsGVR: schema.GroupVersionKind{
		Kind:    "VirtualService",
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	},
}
