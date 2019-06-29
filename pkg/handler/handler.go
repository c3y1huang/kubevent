package handler

import (
	"github.com/innobead/kubevent/pkg/broker"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	controllerruntime "sigs.k8s.io/controller-runtime/pkg/event"
)

type Operation interface {
	Start() error
	Stop() error
}

func sendEvent(broker broker.Operation, event interface{}) {
	eventType := reflect.TypeOf(event).Name()

	var eventMeta interface{}

	switch e := event.(type) {
	case controllerruntime.CreateEvent:
		eventMeta = types.NamespacedName{Namespace: e.Meta.GetNamespace(), Name: e.Meta.GetName()}

	case controllerruntime.UpdateEvent:
		eventMeta = types.NamespacedName{Namespace: e.MetaOld.GetNamespace(), Name: e.MetaOld.GetName()}

	case controllerruntime.DeleteEvent:
		eventMeta = types.NamespacedName{Namespace: e.Meta.GetNamespace(), Name: e.Meta.GetName()}

	case controllerruntime.GenericEvent:
		eventMeta = types.NamespacedName{Namespace: e.Meta.GetNamespace(), Name: e.Meta.GetName()}
	}

	if err := broker.Send(event); err != nil {
		log.Errorf("Failed to send msg from %s, %v", eventType, eventMeta)
	}
}
