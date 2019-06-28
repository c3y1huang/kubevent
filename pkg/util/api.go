package util

import (
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

func CreateCRD(group string, version string, resource runtime.Object) func(c apiextclient.Interface) error {
	name := reflect.TypeOf(resource).Elem().Name()

	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:    group,
			Versions: []apiextv1beta1.CustomResourceDefinitionVersion{{Name: version}},
			Scope:    apiextv1beta1.NamespaceScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Plural: name + "s",
				Kind:   name,
			},
		},
	}

	return func(c apiextclient.Interface) error {
		_, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)

		if err != nil && apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
}

func AddToScheme(group string, version string, resource runtime.Object) func(s *runtime.Scheme) error {
	var SchemeGroupVersion = schema.GroupVersion{
		Group:   group,
		Version: version,
	}

	addKnownTypes := func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(
			SchemeGroupVersion,
			resource,
		)

		metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
		return nil
	}

	bu := runtime.NewSchemeBuilder(addKnownTypes)
	return bu.AddToScheme
}
