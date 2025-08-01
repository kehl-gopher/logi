package jobs

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/mailer"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

type ConsumerManager struct {
	qm             *sync.Mutex
	RM             *rabbitmq.RabbitMQ
	conf           *config.AppConfig
	log            *utils.Log
	workers        int
	queueProcessor []QueueProcessor
	ctx            context.Context
	cancelFunc     context.CancelFunc
	wg             *sync.WaitGroup
	done           chan struct{}
}

func NewConsumerManager(log *utils.Log, conf *config.Config) *ConsumerManager {
	qm := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU() * 2
	done := make(chan struct{})
	qp := make([]QueueProcessor, 0)

	ctx, cancel := context.WithCancel(context.Background())
	return &ConsumerManager{
		qm:             qm,
		wg:             wg,
		done:           done,
		log:            log,
		ctx:            ctx,
		cancelFunc:     cancel,
		queueProcessor: qp,
		workers:        workers,
		conf:           &conf.APP_CONFIG,
	}
}

// I need this to start a processor engine that continously run and
// ensure the queues is being read when data is inserted...
// it listens to events and processes them accordingly
// hopefully you work as expected.. LOL

func (c *ConsumerManager) AddProcessor(q QueueProcessor) {
	c.queueProcessor = append(c.queueProcessor, q)
}

func (c *ConsumerManager) processWorker(p QueueProcessor) {
	defer c.wg.Done()

	backoff := time.Second
	maxBackoff := time.Second * 10

	for {
		select {
		case <-c.ctx.Done():
			utils.PrintLog(c.log, fmt.Sprintf("stopping processor for queue %s", p.name), utils.DebugLevel)
			return
		default:
			if err := c.runProcessor(p); err != nil {
				select {
				case <-time.After(backoff):
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				case <-c.ctx.Done():
					return
				}
			} else {
				backoff = maxBackoff
			}
		}
	}

}

func (c *ConsumerManager) runProcessor(p QueueProcessor) error {
	msg, err := p.rq.ConsumeQueue(context.Background(), p.name, p.exchange, p.routingKey, p.durable)

	if err != nil {
		return fmt.Errorf("failed to consume queue: %w", err)
	}
	var i int = 1
	for ; i <= p.workers; i++ {
		c.wg.Add(1)
		go func(id int) {
			defer c.wg.Done()
			for {
				select {
				case dev, ok := <-msg:
					if !ok {
						utils.PrintLog(c.log, fmt.Sprintf("Queue %s worker %d: channel closed", p.name, id), utils.WarnLevel)
						return
					} else {
						if err := c.processMessage(dev, c.conf, c.log); err != nil {
							utils.PrintLog(c.log, fmt.Sprintf("Worker %d failed to process message: %v", id, err), utils.ErrorLevel)
							dev.Nack(false, true)
						} else {
							dev.Ack(false)
						}
					}
				case <-c.done:
					return
				}
			}
		}(i)
	}

	<-c.ctx.Done()

	go func() {
		c.wg.Wait()
		close(c.done)
	}()

	select {
	case <-c.done:
	case <-time.After(60 * time.Second):
		utils.PrintLog(c.log, fmt.Sprintf("Workers for queue %s didn't finish within 30s", p.name), utils.WarnLevel)
	}
	return nil
}

func (c *ConsumerManager) cancel() {
	c.cancelFunc()
}

func (c *ConsumerManager) Stop() {
	fmt.Println("Stopping consumer manager...")
	c.cancel()
	c.wg.Wait()
	fmt.Println("Consumer manager stopped")
}
func (c *ConsumerManager) Start() {
	for _, p := range c.queueProcessor {
		c.wg.Add(1)
		go c.processWorker(p)
	}
	fmt.Printf("Started %d queue processors\n", len(c.queueProcessor))
}

func (c *ConsumerManager) processMessage(dev amqp091.Delivery, conf *config.AppConfig, lg *utils.Log) error {
	key := dev.RoutingKey
	switch key {
	case "email.welcome":
		var ej mailer.EmailJOB
		err := utils.UnmarshalJSON(dev.Body, &ej)
		if err != nil {
			return err
		}
		c.qm.Lock()
		if err := ej.SendWelcomeEmails(conf, lg); err != nil {
			c.qm.Unlock()
			return err
		}
		c.qm.Unlock()
	case "email.verify":
		var ej mailer.EmailJOB
		err := utils.UnmarshalJSON(dev.Body, &ej)
		if err != nil {
			return err
		}
		c.qm.Lock()
		if err := ej.SendVerificationMail(conf, lg); err != nil {
			c.qm.Unlock()
			return err
		}
		c.qm.Unlock()
	case "email.forgot_password":
		var ej mailer.EmailJOB
		err := utils.UnmarshalJSON(dev.Body, &ej)
		if err != nil {
			return err
		}
		c.qm.Lock()
		if err := ej.SendForgotPassword(conf, lg); err != nil {
			c.qm.Unlock()
			return err
		}
		c.qm.Unlock()
	}
	return nil
}
