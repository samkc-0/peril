package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to RabbitMQ")

	fmt.Println("Starting Peril server...")

	// Ctrl (or Cmd) + C to exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	fmt.Println("Press ^C to exit")
	<-sigs
	fmt.Println("Shutting down Peril server...")
}
