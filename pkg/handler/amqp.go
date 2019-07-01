package handler

import (
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/broker"
	"github.com/innobead/kubevent/pkg/engine"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type AmqpHandler struct {
	BaseHandler
	engine.ControllerEngineAwareType

	broker broker.BrokerOperation
}

func NewAmqpHandler(cfg config.AmqpBroker) *AmqpHandler {
	return &AmqpHandler{
		broker: broker.NewAmqpBroker(cfg),
	}
}

func (receiver *AmqpHandler) Start() error {
	if err := receiver.broker.Start(); err != nil {
		return err
	}

	return nil
}

func (receiver *AmqpHandler) Stop() error {
	if err := receiver.broker.Stop(); err != nil {
		return err
	}

	return nil
}

func (receiver *AmqpHandler) IsInitialized() bool {
	return receiver.IsInitialized()
}

func (receiver *AmqpHandler) Create(event event.CreateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Create(event, queue)
	}
}

func (receiver *AmqpHandler) Update(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Update(event, queue)
	}
}

func (receiver *AmqpHandler) Delete(event event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Delete(event, queue)
	}
}

func (receiver *AmqpHandler) Generic(event event.GenericEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Generic(event, queue)
	}
}
