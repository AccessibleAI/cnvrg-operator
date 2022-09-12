package utils

import (
	"context"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func IsOpenShift(c client.Client) bool {
	routes := &unstructured.UnstructuredList{}
	routes.SetGroupVersionKind(desired.Kinds["OcpRouteGVK"])
	if err := c.List(context.Background(), routes); err != nil {
		return false
	}
	return true
}
