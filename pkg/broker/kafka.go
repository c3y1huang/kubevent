package broker

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/innobead/kubevent/api/v1alpha1"
	"github.com/innobead/kubevent/pkg/util"
	"github.com/sirupsen/logrus"
	"sync"
)

type KafkaBroker struct {
	*v1alpha1.KafkaBroker

	producer sarama.AsyncProducer
	mtx      sync.Mutex
}

func NewKafkaBroker(broker *v1alpha1.KafkaBroker) *KafkaBroker {
	return &KafkaBroker{
		KafkaBroker: broker,
		mtx:         sync.Mutex{},
	}
}

func (k *KafkaBroker) Start() error {
	k.mtx.Lock()
	defer k.mtx.Unlock()

	if k.producer != nil {
		return nil
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	if k.TlsConfig != nil {
		if cfg, _ := util.CreateTLSConfig(k.TlsConfig); cfg != nil {
			config.Net.TLS.Config = cfg
			config.Net.TLS.Enable = true
		}
	}

	if p, err := sarama.NewAsyncProducer(k.KafkaBroker.Addresses, config); err != nil {
		return err
	} else {
		k.producer = p
	}

	return nil
}

func (k *KafkaBroker) Stop() error {
	k.mtx.Lock()
	defer k.mtx.Unlock()

	if k.producer == nil {
		return nil
	}

	if err := k.producer.Close(); err != nil {
		return err
	}

	k.producer = nil

	return nil
}

func (k *KafkaBroker) Send(msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	k.producer.Input() <- &sarama.ProducerMessage{
		Topic: k.Topic,
		Value: sarama.StringEncoder(body),
	}

loop:
	for {
		select {
		case s := <-k.producer.Successes():
			logrus.WithFields(
				logrus.Fields{
					"partition": s.Partition,
					"offset":    s.Offset,
					"message":   string(body),
				},
			).Traceln("event sent")

			break loop
		case err := <-k.producer.Errors():
			return err
		}
	}

	return nil
}
