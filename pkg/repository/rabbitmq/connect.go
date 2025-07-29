package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/rabbitmq/amqp091-go"
)

func (rq *RabbitMQ) RabbitMQConn() error {
	conn, err := amqp091.DialConfig(rq.conf.CONN_STR, amqp091.Config{
		Heartbeat: 5 * time.Second,
	})

	if err != nil {
		utils.PrintLog(rq.log, fmt.Sprintf("failed to establish connection: %v", err), utils.ErrorLevel)
		return err
	}

	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	rq.updateConn(conn)
	rq.updateChannel(ch)
	rq.isConn = true
	utils.PrintLog(rq.log, "RabbitMQ connection successful", utils.DebugLevel)
	return nil
}

func (rq *RabbitMQ) updateConn(conn *amqp091.Connection) {
	rq.conn = conn
}

func (rq *RabbitMQ) updateChannel(chh *amqp091.Channel) {
	rq.ch = chh
}

func (rq *RabbitMQ) PubslishQueue(ctx context.Context, exchange, routingKey, body string, mandatory, immediate bool) error {
	err := rq.ch.PublishWithContext(ctx,
		exchange,
		routingKey,
		mandatory,
		immediate,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)

	if err != nil {
		utils.PrintLog(rq.log, fmt.Sprintf("failed to publish queue message: %s", err.Error()), utils.ErrorLevel)
		return err
	}
	return nil
}

func (rq *RabbitMQ) Close() error {
	if !rq.isConn {
		return fmt.Errorf("no rabbitmq connection established")
	}

	err := rq.conn.Close()

	if err != nil {
		return err
	}
	err = rq.ch.Close()
	if err != nil {
		return err
	}

	return nil
}
