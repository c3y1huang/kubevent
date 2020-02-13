package main

import (
	_ "github.com/innobead/kubevent/cmd/kubevent/version"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := NewKubeventCmd().Execute(); err != nil {
		logrus.WithError(err).Fatalln("Failed to start kubevent process")
	}
}
