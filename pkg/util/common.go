package util

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func RegisterShutdownHook(callback func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopCh

		logrus.Infoln("receiving shutdown hook")
		callback()
	}()
}
