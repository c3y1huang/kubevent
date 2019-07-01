package broker

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/util"
	"sync"
)

type KafkaBroker struct {
	config.KafkaBroker
	producer sarama.SyncProducer
	mtx      sync.Mutex
}

func NewKafkaBroker(cfg config.KafkaBroker) BrokerOperation {
	return &KafkaBroker{
		KafkaBroker: cfg,
	}
}

func (receiver *KafkaBroker) Start() error {
	receiver.mtx.Lock()
	defer receiver.mtx.Unlock()

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
	receiver.mtx.Lock()
	defer receiver.mtx.Unlock()

	if receiver.producer == nil {
		return nil
	}

	_ = receiver.producer.Close()
	receiver.producer = nil

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
