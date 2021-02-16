package desired

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"text/template"
)

type State struct {
	Name           string
	TemplatePath   string
	Template       *template.Template
	ParsedTemplate string
	Obj            *unstructured.Unstructured
	GVR            schema.GroupVersionKind
}

var DeploymentGVR = schema.GroupVersionKind{
	Kind:    "Deployment",
	Group:   "apps",
	Version: "v1",
}

var PvcGVR = schema.GroupVersionKind{
	Kind:    "PersistentVolumeClaim",
	Group:   "",
	Version: "v1",
}

var SecretGVR = schema.GroupVersionKind{
	Kind:    "Secret",
	Group:   "",
	Version: "v1",
}

var SvcGVR = schema.GroupVersionKind{
	Kind:    "Service",
	Group:   "",
	Version: "v1",
}

var SaGVR = schema.GroupVersionKind{
	Kind:    "ServiceAccount",
	Group:   "",
	Version: "v1",
}

var CrdGVR = schema.GroupVersionKind{
	Kind:    "CustomResourceDefinition",
	Group:   "apiextensions.k8s.io",
	Version: "v1",
}

var IstioGVR = schema.GroupVersionKind{
	Kind:    "IstioOperator",
	Group:   "install.istio.io",
	Version: "v1alpha1",
}

var ClusterRoleGVR = schema.GroupVersionKind{
	Kind:    "ClusterRole",
	Group:   "rbac.authorization.k8s.io",
	Version: "v1",
}

var ClusterRoleBindingGVR = schema.GroupVersionKind{
	Kind:    "ClusterRoleBinding",
	Group:   "rbac.authorization.k8s.io",
	Version: "v1",
}

var RoleGVR = schema.GroupVersionKind{
	Kind:    "Role",
	Group:   "rbac.authorization.k8s.io",
	Version: "v1",
}
