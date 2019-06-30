package handler

import (
	"github.com/innobead/kubevent/pkg/broker"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

type Operation interface {
	Start() error
	Stop() error
}

type Base struct {
	handler.EnqueueRequestForObject
}

func sendEvent(broker broker.Operation, e interface{}) error {
	eventType := reflect.TypeOf(e).Name()

	var objName interface{}

	switch e := e.(type) {
	case event.CreateEvent:
		objName = types.NamespacedName{Namespace: e.Meta.GetNamespace(), Name: e.Meta.GetName()}

	case event.UpdateEvent:
		objName = types.NamespacedName{Namespace: e.MetaOld.GetNamespace(), Name: e.MetaOld.GetName()}

	case event.DeleteEvent:
		objName = types.NamespacedName{Namespace: e.Meta.GetNamespace(), Name: e.Meta.GetName()}

	case event.GenericEvent:
		objName = types.NamespacedName{Namespace: e.Meta.GetNamespace(), Name: e.Meta.GetName()}
	}

	log := log.WithFields(log.Fields{
		"type": eventType,
		"name": objName,
	})

	if err := broker.Send(e); err != nil {
		log.Errorf("Failed to send event, %v", err)
		return err
	} else {
		log.Debug("Sending event")
	}

	return nil
}
