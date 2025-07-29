package semail

import (
	"context"
	"fmt"
	"time"

	"github.com/kehl-gopher/logi/internal/jobs"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
)

func PublishToEmailQUeue(rq *rabbitmq.RabbitMQ, name, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := rq.DeclareQueue(ctx, name, routingKey, true, body)
	if err != nil {
		return err
	}

	// consume queue background tasks
	q := jobs.NewQueueProcessor(rq, name, true, routingKey)
	q.Start()

	fmt.Println("------------------_>")
	return err
}
