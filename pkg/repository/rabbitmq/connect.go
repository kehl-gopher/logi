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

func (rq *RabbitMQ) PublishQueue(ctx context.Context, exchange, routingKey string, body []byte) error {
	err := rq.ch.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
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

func (rq *RabbitMQ) DeclareQueue(ctx context.Context, name string, routingKey string, durable bool, body []byte) error {
	q, err := rq.ch.QueueDeclare(routingKey, durable, false, false, false, nil)
	if err != nil {
		utils.PrintLog(rq.log, fmt.Sprintf("failed to declare queue: %v", err), utils.ErrorLevel)
		return err
	}
	err = rq.PublishQueue(ctx, name, q.Name, body)
	if err != nil {
		return err
	}
	return err
}

func (rq *RabbitMQ) ConsumeQueue(ctx context.Context, name string, routingKey string, durable bool) (<-chan amqp091.Delivery, error) {
	q, err := rq.ch.QueueDeclare(routingKey, durable, false, false, false, nil)

	if err != nil {
		return nil, err
	}

	dev, err := rq.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

func (rq *RabbitMQ) RQChan() *amqp091.Channel {
	return rq.ch
}
