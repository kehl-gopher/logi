package jobs

import (
	"runtime"

	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
)

// exchange
const (
	EMAIL_EXCHANGE = "email-exchange"
)

// queues
const (
	EMAIL_QUEUE = "email-queue"
)

type QueueProcessor struct {
	name       string
	exchange   string
	durable    bool
	routingKey string
	workers    int
	rq         *rabbitmq.RabbitMQ
}

func NewQueueProcessor(name, exchange, routingKey string, durable bool, rq *rabbitmq.RabbitMQ) QueueProcessor {
	workers := runtime.NumCPU() * 2
	return QueueProcessor{
		name:       name,
		exchange:   exchange,
		routingKey: routingKey,
		durable:    durable,
		workers:    workers,
		rq:         rq,
	}
}
