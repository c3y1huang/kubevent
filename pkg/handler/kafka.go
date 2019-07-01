package handler

import (
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/broker"
	"github.com/innobead/kubevent/pkg/engine"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type KafkaHandler struct {
	BaseHandler
	engine.ControllerEngineAwareType

	broker broker.BrokerOperation
}

func NewKafkaHandler(cfg config.KafkaBroker) *KafkaHandler {
	return &KafkaHandler{
		broker: broker.NewKafkaBroker(cfg),
	}
}

func (receiver *KafkaHandler) Start() error {
	if err := receiver.broker.Start(); err != nil {
		return err
	}

	return nil
}

func (receiver *KafkaHandler) Stop() error {
	if err := receiver.broker.Stop(); err != nil {
		return err
	}

	return nil
}

func (receiver *KafkaHandler) Create(event event.CreateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Create(event, queue)
	}
}

func (receiver *KafkaHandler) Update(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Update(event, queue)
	}
}

func (receiver *KafkaHandler) Delete(event event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Delete(event, queue)
	}
}

func (receiver *KafkaHandler) Generic(event event.GenericEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Generic(event, queue)
	}
}
