package main

import (
	"fmt"
	"log"
	"peril/internal/gamelogic"
	"peril/internal/pubsub"
	"peril/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to create a channel:\n%v\n", err)

	}
	fmt.Println("Connected to RabbitMQ")

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.DurableQueueType)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Declared and bound queue game_logs")

	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if words == nil || len(words) == 0 {
			continue
		}
		if words[0] == gamelogic.CmdServerPause {
			fmt.Println("Sending pause command...")
			if err := pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			); err != nil {
				fmt.Printf("FAILED:\n%v\n", err)
			}

		} else if words[0] == gamelogic.CmdServerResume {
			fmt.Println("Sending resume command...")

			if err := pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			); err != nil {
				fmt.Printf("FAILED:\n%v\n", err)
			}

		} else if words[0] == gamelogic.CmdServerQuit {
			fmt.Println("Goodbye!")
			break
		} else if words[0] == gamelogic.CmdServerHelp {
			gamelogic.PrintServerHelp()
		} else {
			fmt.Printf("unknown command: %s\n", words[0])
		}
	}
}
