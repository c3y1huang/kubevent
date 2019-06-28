package main

import (
	"github.com/innobead/kubevent/internal/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}
