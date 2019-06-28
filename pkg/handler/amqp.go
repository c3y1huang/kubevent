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
	engine.ControllerEngineAwareType
	broker broker.Operation
}

func NewAmqpEventHandler(cfg config.AmqpSink) *AmqpEventHandler {
	handler := AmqpEventHandler{
		broker: &message.AmqpBroker{AmqpSink: cfg},
	}

	return &handler
}

func (receiver *AmqpEventHandler) Create(event event.CreateEvent, _ workqueue.RateLimitingInterface) {
	sendEvent(receiver.broker, event)
}

func (receiver *AmqpEventHandler) Update(event event.UpdateEvent, _ workqueue.RateLimitingInterface) {
	sendEvent(receiver.broker, event)
}

func (receiver *AmqpEventHandler) Delete(event event.DeleteEvent, _ workqueue.RateLimitingInterface) {
	sendEvent(receiver.broker, event)
}

func (receiver *AmqpEventHandler) Generic(event event.GenericEvent, _ workqueue.RateLimitingInterface) {
	sendEvent(receiver.broker, event)
}
