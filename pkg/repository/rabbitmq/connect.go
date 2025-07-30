package rabbitmq

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/rabbitmq/amqp091-go"
)

func (rq *RabbitMQ) RabbitMQConn() error {
	const maxRetries = 5
	var backoff = time.Second
	const maxBackoff = time.Second * 5
	var ind int = 0
	var conn *amqp091.Connection
	var ch *amqp091.Channel
	var err error

	for ; ind <= 5; ind++ {
		conn, err = amqp091.DialConfig(rq.conf.CONN_STR, amqp091.Config{
			Heartbeat: 10 * time.Second,
		})

		if err != nil {
			utils.PrintLog(rq.log, fmt.Sprintf("RabbitMQ connection attempt %d/%d failed: %v", ind, maxRetries, err), utils.WarnLevel)
			if ind < maxRetries {
				utils.PrintLog(rq.log, fmt.Sprintf("Retrying in %v...", backoff), utils.InfoLevel)
				time.Sleep(backoff)
				backoff = time.Duration(float64(backoff) * (1.5 + rand.Float64()*0.5))

				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
			continue
		}

		ch, err = conn.Channel()

		if err != nil {
			utils.PrintLog(rq.log, fmt.Sprintf("Failed to create channel (attempt %d/%d): %v", ind, maxRetries, err), utils.WarnLevel)
			conn.Close()

			if ind < maxRetries {
				backoff = time.Duration(float64(backoff) * 0.5)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}

			continue
		}
		break
	}

	if err != nil {
		utils.PrintLog(rq.log, fmt.Sprintf("Failed to establish RabbitMQ connection after %d attempts: %v", maxRetries, err), utils.ErrorLevel)
		return fmt.Errorf("rabbitmq connection failed after %d attempts: %w", maxRetries, err)
	}

	rq.setupConnectionMonitoring()

	rq.updateChannel(ch)
	rq.updateConn(conn)

	rq.isConn = true

	utils.PrintLog(rq.log, "RabbitMQ connection successful", utils.DebugLevel)
	return nil
}

func (rq *RabbitMQ) setupConnectionMonitoring() {
	go func() {
		connClose := rq.notifyConnClosed
		chClose := rq.notifyChanClosed

		select {
		case err := <-connClose:
			if err != nil {
				utils.PrintLog(rq.log, fmt.Sprintf("RabbitMQ connection closed: %v", err), utils.ErrorLevel)
				rq.isConn = false
				rq.RabbitMQConn()
			}
		case err := <-chClose:
			if err != nil {
				utils.PrintLog(rq.log, fmt.Sprintf("RabbitMQ channel closed: %v", err), utils.ErrorLevel)
				rq.isConn = false
				rq.RabbitMQConn()
			}
		}
	}()
}

func (rq *RabbitMQ) IsHealthy() bool {
	if !rq.isConn {
		return false
	}

	if rq.conn != nil && rq.conn.IsClosed() {
		return false
	}

	if rq.ch != nil && rq.ch.IsClosed() {
		return false
	}

	return true
}

func (rq *RabbitMQ) EnsureConnection() error {
	if rq.IsHealthy() {
		return nil
	}
	utils.PrintLog(rq.log, "Connection unhealthy, attempting to reconnect...", utils.InfoLevel)
	return rq.RabbitMQConn()
}

func (rq *RabbitMQ) updateConn(conn *amqp091.Connection) {
	rq.conn = conn
	rq.notifyConnClosed = make(chan *amqp091.Error)
	rq.conn.NotifyClose(rq.notifyConnClosed)
}

func (rq *RabbitMQ) updateChannel(chh *amqp091.Channel) {
	rq.ch = chh
	rq.notifyChanClosed = make(chan *amqp091.Error)
	rq.ch.NotifyClose(rq.notifyChanClosed)
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

func (rq *RabbitMQ) DeclareQueue(ctx context.Context, queueName, routingKey, exchange string, durable bool, body []byte) error {
	err := rq.ch.ExchangeDeclare(exchange, "topic", durable, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)

	}
	// Declare queue
	q, err := rq.ch.QueueDeclare(queueName, durable, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = rq.ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)

	}
	err = rq.PublishQueue(ctx, exchange, routingKey, body)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return err
}

func (rq *RabbitMQ) ConsumeQueue(ctx context.Context, name, exchange, routingKey string, durable bool) (<-chan amqp091.Delivery, error) {
	err := rq.ch.ExchangeDeclare(exchange, "topic", durable, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}
	q, err := rq.ch.QueueDeclare(name, durable, false, false, false, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	err = rq.ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
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
