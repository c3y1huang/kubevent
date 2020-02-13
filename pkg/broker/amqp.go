package broker

import (
	"encoding/json"
	"github.com/innobead/kubevent/api/v1alpha1"
	"github.com/innobead/kubevent/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

type AMQPBroker struct {
	*v1alpha1.AMQPBroker

	conn *amqp.Connection
	mtx  sync.Mutex
}

func NewAMQPBroker(broker *v1alpha1.AMQPBroker) *AMQPBroker {
	return &AMQPBroker{
		AMQPBroker: broker,
		mtx:        sync.Mutex{},
	}
}

func (a *AMQPBroker) Start() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.conn != nil {
		return nil
	}

	if a.TlsConfig == nil {
		if conn, err := amqp.Dial(a.Addresses[0]); err != nil {
			return err
		} else {
			a.conn = conn
		}

		return nil
	}

	if cfg, err := util.CreateTLSConfig(a.TlsConfig); err != nil {
		if conn, err := amqp.DialTLS(
			a.Addresses[0],
			cfg,
		); err != nil {
			return err
		} else {
			a.conn = conn
		}
	}

	return nil
}

func (a *AMQPBroker) Stop() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.conn == nil {
		return nil
	}

	if err := a.conn.Close(); err != nil {
		logrus.WithError(err).Errorln("failed to close AMQP broker channel")
		return err
	}

	a.conn = nil

	return nil
}

func (a *AMQPBroker) Send(msg interface{}) error {
	ch, err := a.conn.Channel()
	if err != nil {
		if err == amqp.ErrClosed {
			_ = a.Stop()
		}

		return err
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			logrus.WithError(err).Errorln("failed to close AMQP broker channel")
		}
	}()

	if err := ch.ExchangeDeclare(
		a.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		a.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}

	return err
}
