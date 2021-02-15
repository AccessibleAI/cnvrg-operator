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
	Group:   "",
	Version: "apps/v1",
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
