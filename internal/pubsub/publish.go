package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp091.Channel, exchange string, key string, val T) error {
	payload, err := json.Marshal(val)

	if err != nil {
		return err
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp091.Publishing{ContentType: "application/json", Body: payload})

	return nil
}
