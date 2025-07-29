package jobs

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
)

type QueueProcessor struct {
	RM         *rabbitmq.RabbitMQ
	workers    int
	name       string
	receiver   interface{}
	durable    bool
	routingKey string
}

func NewQueueProcessor(r *rabbitmq.RabbitMQ, name string, durable bool, routingKey string) *QueueProcessor {
	workers := runtime.NumCPU() * 2
	return &QueueProcessor{RM: r, workers: workers, name: name, durable: durable, routingKey: routingKey}
}

func (q *QueueProcessor) ProcessQueue() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msgs, err := q.RM.ConsumeQueue(ctx, q.name, q.routingKey, q.durable)
	if err != nil {
		log.Printf("failed to read data from queue: %v", err)
		panic(err)
	}
	for i := 1; i <= q.workers; i++ {
		go func() {
			for d := range msgs {
				fmt.Printf("Worker processing message: %s\n", string(d.Body))
				d.Ack(false)
			}
		}()
	}
}

func (q *QueueProcessor) Start() {
	fmt.Println("start queue... jor")
	go q.ProcessQueue()
}
