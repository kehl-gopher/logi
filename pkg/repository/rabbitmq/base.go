package rabbitmq

import (
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conf   *config.RabbitMQ
	log    *utils.Log
	conn   *amqp091.Connection
	isConn bool
	ch     *amqp091.Channel
}

func NewMQManager(conf *config.RabbitMQ, log *utils.Log) *RabbitMQ {
	re := &RabbitMQ{conf: conf, log: log}
	re.RabbitMQConn()
	return re
}
