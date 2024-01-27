package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"gamelight/internal/types"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPProducer struct {
	*AMQPEngine
}

func NewAMQPProducer(url string, queueName string) (*AMQPProducer, error) {
	engine, err := NewAMQPEngine(url, queueName)
	if err != nil {
		return nil, err
	}

	return &AMQPProducer{
		AMQPEngine: engine,
	}, nil
}

func (producer *AMQPProducer) Produce(ctx context.Context, event *types.Message) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("can not marshal json, err: %v", err)
	}

	err = producer.Channel.PublishWithContext(
		ctxTimeout,
		"",
		producer.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to push to queue, err: %v", err)
	}

	return nil
}
