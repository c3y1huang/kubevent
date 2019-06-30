package handler

import (
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/broker"
	"github.com/innobead/kubevent/pkg/broker/stream"
	"github.com/innobead/kubevent/pkg/engine"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type Kafka struct {
	Base
	engine.ControllerEngineAwareType

	broker broker.Operation
}

func NewKafka(cfg config.KafkaBroker) *Kafka {
	return &Kafka{
		broker: stream.NewKafkaBroker(cfg),
	}
}

func (receiver *Kafka) Start() error {
	if err := receiver.broker.Start(); err != nil {
		return err
	}

	return nil
}

func (receiver *Kafka) Stop() error {
	if err := receiver.broker.Stop(); err != nil {
		return err
	}

	return nil
}

func (receiver *Kafka) Create(event event.CreateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Create(event, queue)
	}
}

func (receiver *Kafka) Update(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Update(event, queue)
	}
}

func (receiver *Kafka) Delete(event event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Delete(event, queue)
	}
}

func (receiver *Kafka) Generic(event event.GenericEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(receiver.broker, event); err != nil {
		receiver.EnqueueRequestForObject.Generic(event, queue)
	}
}
