package util

import (
	"os"
	"os/signal"
	"syscall"
)

func RegisterShutdownHook(callback func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopCh
		callback()
	}()
}
