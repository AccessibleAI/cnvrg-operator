package desired

import "k8s.io/apimachinery/pkg/runtime/schema"

type GVKName string

const (
	DeploymentGVK  GVKName = "DeploymentGVK"
	StatefulSetGVK GVKName = "StatefulSetGVK"
	JobGVK         GVKName = "JobGvk"
)

var Kinds = map[GVKName]schema.GroupVersionKind{
	DeploymentGVK: {
		Kind:    "Deployment",
		Group:   "apps",
		Version: "v1",
	},

	StatefulSetGVK: {
		Kind:    "StatefulSet",
		Group:   "apps",
		Version: "v1",
	},

	JobGVK: {
		Kind:    "Job",
		Group:   "batch",
		Version: "v1",
	},
}
