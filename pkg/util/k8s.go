package util

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetK8sClient(cfg *rest.Config) (kubernetes.Interface, error) {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return clientset, err
}
