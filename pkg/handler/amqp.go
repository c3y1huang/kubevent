package handler

import (
	"github.com/innobead/kubevent/api/v1alpha1"
	"github.com/innobead/kubevent/pkg/broker"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type AmqpHandler struct {
	EventHandler
	broker broker.Operation
}

func NewAmqpHandler(ab *v1alpha1.AMQPBroker) *AmqpHandler {
	return &AmqpHandler{
		broker: broker.NewAMQPBroker(ab),
	}
}

func (a *AmqpHandler) Start() error {
	if err := a.broker.Start(); err != nil {
		return err
	}

	return nil
}

func (a *AmqpHandler) Stop() error {
	if err := a.broker.Stop(); err != nil {
		return err
	}

	return nil
}

func (a *AmqpHandler) Create(event event.CreateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(a.broker, event); err != nil {
		if queue != nil {
			a.EnqueueRequestForObject.Create(event, queue)
		}
	}
}

func (a *AmqpHandler) Update(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(a.broker, event); err != nil {
		if queue != nil {
			a.EnqueueRequestForObject.Update(event, queue)
		}
	}
}

func (a *AmqpHandler) Delete(event event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(a.broker, event); err != nil {
		if queue != nil {
			a.EnqueueRequestForObject.Delete(event, queue)
		}
	}
}

func (a *AmqpHandler) Generic(event event.GenericEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(a.broker, event); err != nil {
		if queue != nil {
			a.EnqueueRequestForObject.Generic(event, queue)
		}
	}
}
