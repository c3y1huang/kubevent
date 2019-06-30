package stream

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/broker"
	"github.com/innobead/kubevent/pkg/util"
)

type KafkaBroker struct {
	config.KafkaBroker
	producer sarama.SyncProducer
}

func NewKafkaBroker(cfg config.KafkaBroker) broker.Operation {
	return &KafkaBroker{
		KafkaBroker: cfg,
	}
}

func (receiver *KafkaBroker) Start() error {
	if receiver.producer != nil {
		return nil
	}

	var err error

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Return.Successes = true

	if tlsCfg, _ := util.CreateTLSConfig(receiver.Tls); tlsCfg != nil {
		cfg.Net.TLS.Config = tlsCfg
		cfg.Net.TLS.Enable = true
	}

	receiver.producer, err = sarama.NewSyncProducer(receiver.KafkaBroker.Servers, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *KafkaBroker) Stop() error {
	if receiver.producer == nil {
		return nil
	}

	_ = receiver.producer.Close()
	return nil
}

func (receiver *KafkaBroker) IsInitialized() bool {
	return receiver.producer != nil
}

func (receiver *KafkaBroker) Send(msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, _, err = receiver.producer.SendMessage(&sarama.ProducerMessage{
		Topic: receiver.Topic,
		Value: sarama.StringEncoder(body),
	})
	if err != nil {
		return err
	}

	return nil
}
