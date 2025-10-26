package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	DurableQueueType SimpleQueueType = iota
	TransientQueueType
)

func PublishJSON[T any](ch *amqp091.Channel, exchange string, key string, val T) error {
	payload, err := json.Marshal(val)

	if err != nil {
		return err
	}

	ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)

	return nil
}

func SubscribeJSON[T any](
	conn *amqp091.Connection,
	exchange, queueName, key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	return nil
}

func DeclareAndBind(
	conn *amqp091.Connection,
	exchange, queueName, key string,
	queueType SimpleQueueType,
) (*amqp091.Channel, amqp091.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	durable := (queueType == DurableQueueType)
	autoDelete := (queueType == TransientQueueType)
	exclusive := (queueType == TransientQueueType)
	noWait := false

	queue, err := channel.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)

	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	channel.QueueBind(queueName, exchange, key, noWait, nil)
	return channel, queue, nil

}
