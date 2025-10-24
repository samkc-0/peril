package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	channel, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueueType,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting Peril client...")

	fmt.Println("Press Ctrl+C to exit...")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
