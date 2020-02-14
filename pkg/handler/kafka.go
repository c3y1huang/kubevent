package handler

import (
	"github.com/innobead/kubevent/api/v1alpha1"
	"github.com/innobead/kubevent/pkg/broker"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type KafkaHandler struct {
	EventHandler
	broker broker.Operation
}

func NewKafkaHandler(kb *v1alpha1.KafkaBroker) *KafkaHandler {
	return &KafkaHandler{
		broker: broker.NewKafkaBroker(kb),
	}
}

func (k *KafkaHandler) Start() error {
	if err := k.broker.Start(); err != nil {
		return err
	}
	return nil
}

func (k *KafkaHandler) Stop() error {
	if err := k.broker.Stop(); err != nil {
		return err
	}
	return nil
}

func (k *KafkaHandler) Create(event event.CreateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(k.broker, event); err != nil {
		logrus.WithError(err).Errorln("failed to send event")
		return
	}

	k.EnqueueRequestForObject.Create(event, queue)
}

func (k *KafkaHandler) Update(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(k.broker, event); err != nil {
		logrus.WithError(err).Errorln("failed to send event")
		return
	}

	k.EnqueueRequestForObject.Update(event, queue)
}

func (k *KafkaHandler) Delete(event event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(k.broker, event); err != nil {
		logrus.WithError(err).Errorln("failed to send event")
		return
	}

	k.EnqueueRequestForObject.Delete(event, queue)
}

func (k *KafkaHandler) Generic(event event.GenericEvent, queue workqueue.RateLimitingInterface) {
	if err := sendEvent(k.broker, event); err != nil {
		logrus.WithError(err).Errorln("failed to send event")
		return
	}

	k.EnqueueRequestForObject.Generic(event, queue)
}
