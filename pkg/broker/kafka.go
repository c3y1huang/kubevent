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

	producer sarama.SyncProducer
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

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 10
	kafkaConfig.Producer.Return.Successes = true

	if k.TlsConfig != nil {
		if cfg, _ := util.CreateTLSConfig(k.TlsConfig); cfg != nil {
			kafkaConfig.Net.TLS.Config = cfg
			kafkaConfig.Net.TLS.Enable = true
		}
	}

	if p, err := sarama.NewSyncProducer(k.KafkaBroker.Addresses, kafkaConfig); err != nil {
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

	part, offset, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.Topic,
		Value: sarama.StringEncoder(body),
	})
	if err != nil {
		return err
	}
	logrus.WithFields(
		logrus.Fields{
			"partition": part,
			"offset":    offset,
			"message":   string(body),
		},
	).Debugln("message sent")

	return nil
}
