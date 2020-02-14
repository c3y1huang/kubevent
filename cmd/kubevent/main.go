package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		logrus.WithError(err).Fatalln("failed to start Kubevent CLI")
	}
}
