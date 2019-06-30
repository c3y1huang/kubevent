package message

import (
	"encoding/json"
	"github.com/innobead/kubevent/internal/config"
	"github.com/innobead/kubevent/pkg/broker"
	er "github.com/innobead/kubevent/pkg/error"
	"github.com/innobead/kubevent/pkg/util"
	"github.com/streadway/amqp"
)

type AmqpBroker struct {
	config.AmqpBroker
	conn *amqp.Connection
}

func NewAmqpBroker(cfg config.AmqpBroker) broker.Operation {
	return &AmqpBroker{
		AmqpBroker: cfg,
	}
}

func (receiver *AmqpBroker) Start() error {
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
	if receiver.conn == nil {
		return nil
	}

	if err := receiver.conn.Close(); err != nil {
		return err
	}
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
