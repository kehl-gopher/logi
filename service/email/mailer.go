package semail

import (
	"context"
	"time"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
)

func PublishToEmailQUeue(rq *rabbitmq.RabbitMQ, name, routingKey, exchange string, body []byte, log *utils.Log, conf *config.AppConfig) error {

	time.Sleep(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := rq.DeclareQueue(ctx, name, routingKey, exchange, true, body)
	if err != nil {
		return err
	}

	return err
}
