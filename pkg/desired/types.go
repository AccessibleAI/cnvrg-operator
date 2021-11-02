package desired

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"text/template"
)

type TemplateData struct {
	Namespace string
	Data      map[string]interface{}
}

type State struct {
	TemplatePath   string
	Template       *template.Template
	ParsedTemplate string
	Obj            *unstructured.Unstructured
	GVK            schema.GroupVersionKind
	Own            bool
	Override       bool
	Updatable      bool
	TemplateData   interface{}
}
