package handler

import (
	"github.com/innobead/kubevent/pkg/broker"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func sendEvent(broker broker.Operation, event interface{}) {
	eventType := reflect.TypeOf(event).Name()

	if err := broker.Send(event); err != nil {
		log.Errorf("Failed to send msg from %s, %v", eventType, event)
	}
}
