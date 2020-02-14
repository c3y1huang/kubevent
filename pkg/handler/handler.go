package handler

import (
	"errors"
	"github.com/innobead/kubevent/api/v1alpha1"
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

type EventHandler struct {
	handler.EventHandler
	handler.EnqueueRequestForObject
}

func sendEvent(broker broker.Operation, e interface{}) error {
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

	logger := log.WithFields(log.Fields{
		"type": reflect.TypeOf(e).Name(),
		"name": objName,
	})

	logger.Debugln("sending event")

	if err := broker.Send(e); err != nil {
		logger.Errorf("failed to send event, %v", err)
		return err
	}

	return nil
}

func CreateEventBrokerHandler(spec *v1alpha1.EventBrokerSpec) (handler.EventHandler, Operation, error) {
	switch {
	case spec.Kafka != nil:
		h := NewKafkaHandler(spec.Kafka)
		return h, h, nil

	case spec.AMQP != nil:
		h := NewAmqpHandler(spec.AMQP)
		return h, h, nil
	}

	return nil, nil, errors.New("no event broker configured")
}
