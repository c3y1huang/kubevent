package util

import (
	"github.com/innobead/kubevent/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func ToSchemaGroupVersionKinds(kinds []*v1alpha1.GroupVersionKind) []schema.GroupVersionKind {
	var gvks []schema.GroupVersionKind
	for _, kind := range kinds {
		gvks = append(gvks, schema.GroupVersionKind{
			Group:   kind.Group,
			Version: kind.Version,
			Kind:    kind.Kind,
		})
	}

	return gvks
}
