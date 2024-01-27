package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"gamelight/internal/types"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type AMQPConsumer struct {
	*AMQPEngine

	Messages <-chan amqp.Delivery
}

func NewAMQPConsumer(url string, queueName string) (*AMQPConsumer, error) {
	engine, err := NewAMQPEngine(url, queueName)
	if err != nil {
		return nil, err
	}

	messages, err := engine.Channel.Consume(
		engine.Queue.Name,
		"message_consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("can not start consuming, err: %v", err)
	}

	return &AMQPConsumer{
		engine,
		messages,
	}, nil
}

type handleMessage func(ctx context.Context, event *types.Message) error

func (consumer *AMQPConsumer) Consume(ctx context.Context, handle handleMessage) error {
	for {
		select {
		case msg := <-consumer.Messages:
			var message types.Message

			if err := json.Unmarshal(msg.Body, &message); err != nil {
				if err = msg.Nack(false, true); err != nil {
					return fmt.Errorf("can not nack message, err: %v", err)
				}
			}

			if err := handle(ctx, &message); err != nil {
				logrus.Errorf("can not process message from the queue, err: %s. Requeue", err.Error())

				if err = msg.Nack(false, true); err != nil {
					return fmt.Errorf("can not nack message, err: %v", err)
				}

				continue
			}

			if err := msg.Ack(true); err != nil {
				return fmt.Errorf("can not ack message, err: %v", err)
			}
		}
	}
}
