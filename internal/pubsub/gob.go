package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"peril/internal/routing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func encodeGob[T any](val T) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(val); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decodeGob[T any](data []byte) (T, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var result T
	if err := decoder.Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func PublishGob[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	payload, err := encodeGob(val)
	if err != nil {
		return err
	}

	if err = ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false, false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        payload,
		},
	); err != nil {
		return err
	}
	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange, queueName, key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	ch.Qos(10, 0, true)
	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			if msg, err := decodeGob[T](delivery.Body); err == nil {
				switch handler(msg) {
				case Ack:
					delivery.Ack(false)
				case NackRequeue:
					delivery.Nack(false, true)
				case NackDiscard:
					delivery.Nack(false, false)
				default:
					delivery.Nack(false, false)
				}
			} else {
				fmt.Println(err)
			}
		}
	}()
	return nil
}

func PublishGameLog(ch *amqp.Channel, message, username string) error {
	err := PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+username, routing.GameLog{CurrentTime: time.Now(), Message: message, Username: username})
	if err != nil {
		return err
	}
	return nil
}
