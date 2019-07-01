package broker

import (
	"encoding/json"
	"github.com/innobead/kubevent/internal/config"
	er "github.com/innobead/kubevent/pkg/error"
	"github.com/innobead/kubevent/pkg/util"
	"github.com/streadway/amqp"
	"sync"
)

type AmqpBroker struct {
	config.AmqpBroker
	conn *amqp.Connection
	mtx  sync.Mutex
}

func NewAmqpBroker(cfg config.AmqpBroker) BrokerOperation {
	return &AmqpBroker{
		AmqpBroker: cfg,
	}
}

func (receiver *AmqpBroker) Start() error {
	receiver.mtx.Lock()
	defer receiver.mtx.Unlock()

	if receiver.conn != nil {
		return nil
	}

	var err error

	if tlsConfig, _ := util.CreateTLSConfig(receiver.Tls); tlsConfig != nil {
		receiver.conn, err = amqp.DialTLS(
			receiver.Uri,
			tlsConfig,
		)
		if err != nil {
			return err
		}

		return nil
	}

	receiver.conn, err = amqp.Dial(receiver.Uri)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *AmqpBroker) Stop() error {
	receiver.mtx.Lock()
	defer receiver.mtx.Unlock()

	if receiver.conn == nil {
		return nil
	}

	_ = receiver.conn.Close()
	receiver.conn = nil

	return nil
}

func (receiver *AmqpBroker) IsInitialized() bool {
	return receiver.conn != nil
}

func (receiver *AmqpBroker) Send(msg interface{}) error {
	if !receiver.IsInitialized() {
		return er.NotInitialized
	}

	ch, err := receiver.conn.Channel()
	if err != nil {
		if err == amqp.ErrClosed {
			_ = receiver.Stop()
		}

		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		receiver.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		receiver.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	return err
}
