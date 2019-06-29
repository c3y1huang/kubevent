package handler

import (
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/broker"
	"github.com/innobead/kubevent/pkg/broker/message"
	"github.com/innobead/kubevent/pkg/engine"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type AmqpEventHandler struct {
	BaseHandler
	engine.ControllerEngineAwareType

	broker broker.Operation
}

func NewAmqpEventHandler(cfg config.AmqpSink) *AmqpEventHandler {
	handler := AmqpEventHandler{
		broker: &message.AmqpBroker{AmqpSink: cfg},
	}

	return &handler
}

func (receiver *AmqpEventHandler) Start() error {
	if err := receiver.broker.Start(); err != nil {
		return err
	}

	return nil
}

func (receiver *AmqpEventHandler) Stop() error {
	if err := receiver.broker.Stop(); err != nil {
		return err
	}

	return nil
}

func (receiver *AmqpEventHandler) Create(event event.CreateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Create(event, queue)
	}
}

func (receiver *AmqpEventHandler) Update(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Update(event, queue)
	}
}

func (receiver *AmqpEventHandler) Delete(event event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Delete(event, queue)
	}
}

func (receiver *AmqpEventHandler) Generic(event event.GenericEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Generic(event, queue)
	}
}
