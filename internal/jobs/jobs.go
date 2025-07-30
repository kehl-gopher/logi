package jobs

import (
	"runtime"
	"sync"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
)

type QueueProcessor struct {
	qm         *sync.Mutex
	RM         *rabbitmq.RabbitMQ
	workers    int
	name       string
	exchange   string
	durable    bool
	routingKey string
	conf       *config.AppConfig
	log        *utils.Log
}

func NewQueueProcessor(r *rabbitmq.RabbitMQ, name string, durable bool, routingKey, exchange string, log *utils.Log, conf *config.AppConfig) *QueueProcessor {
	workers := runtime.NumCPU() * 2
	return &QueueProcessor{
		qm:         new(sync.Mutex),
		RM:         r,
		workers:    workers,
		name:       name,
		durable:    durable,
		routingKey: routingKey,
		exchange:   exchange,
		log:        log,
		conf:       conf,
	}
}

func (rq *QueueProcessor) Processor() {

}
