package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	DurableQueueType SimpleQueueType = iota
	TransientQueueType
)

func PublishJSON[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	payload, err := json.Marshal(val)

	if err != nil {
		return err
	}

	if err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	); err != nil {
		return err
	}
	return nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange, queueName, key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			var msg T
			if err := json.Unmarshal(delivery.Body, &msg); err == nil {
				handler(msg)
			}
			delivery.Ack(false)
		}
	}()
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange, queueName, key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
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
		return nil, amqp.Queue{}, err
	}

	if err = channel.QueueBind(queueName, key, exchange, noWait, nil); err != nil {
		return nil, amqp.Queue{}, err
	}
	return channel, queue, nil

}
