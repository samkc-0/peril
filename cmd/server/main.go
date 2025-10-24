package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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

	if err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
		fmt.Printf("failed to publish json:\n%v\n", err)
	}

	fmt.Println("Connected to RabbitMQ")

	fmt.Println("Starting Peril server...")

	// Ctrl (or Cmd) + C to exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	fmt.Println("Press ^C to exit")
	<-sigs
	fmt.Println("Shutting down Peril server...")
}
