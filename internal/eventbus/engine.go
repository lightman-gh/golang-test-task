package eventbus

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPEngine struct {
	Queue   amqp.Queue
	Channel *amqp.Channel
}

func NewAMQPEngine(url string, queueName string) (*AMQPEngine, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("can not dial to: %s, err: %v", url, err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can not create a channel, err: %v", err)
	}

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("can not declar a new queue, err: %v", err)
	}

	return &AMQPEngine{
		Channel: channel,
		Queue:   queue,
	}, nil
}
