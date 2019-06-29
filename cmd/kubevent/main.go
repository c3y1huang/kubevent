package main

import (
	"github.com/innobead/kubevent/internal/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
