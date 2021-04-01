package desired

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"text/template"
)

type State struct {
	TemplatePath   string
	Template       *template.Template
	ParsedTemplate string
	Obj            *unstructured.Unstructured
	GVR            schema.GroupVersionKind
	Own            bool
	Override       bool
	TemplateData   interface{}
}
