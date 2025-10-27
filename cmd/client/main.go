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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	gameState := gamelogic.NewGameState(username)

	fmt.Println("subscribing to pause events")
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueueType,
		handlerPause(gameState),
	); err != nil {
		fmt.Println(err)
	}

	fmt.Println("subscribing to move events")

	publishCh, err := conn.Channel()
	if err != nil {
		fmt.Println("error creating channel:")
		log.Fatal(err)
	}

	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.TransientQueueType,
		handlerMove(gameState, publishCh),
	); err != nil {
		fmt.Println(err)
	}

	fmt.Println("subscribing to war outcomes")
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		"war",
		routing.WarRecognitionPrefix+".*",
		pubsub.DurableQueueType,
		handlerWar(gameState),
	); err != nil {
		fmt.Println(err)
	}

	fmt.Println("entering REPL")
	for {
		words := gamelogic.GetInput()
		if words == nil || len(words) == 0 {
			continue
		}
		switch words[0] {
		case gamelogic.CmdClientSpawn:
			if err := gameState.CommandSpawn(words); err != nil {
				fmt.Println(err)
			}
		case gamelogic.CmdClientMove:
			move, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+username,
				move,
			); err != nil {
				fmt.Printf("publishing move failed: %s\n", err)
				continue
			}
			fmt.Println("moved units successfully")
		case gamelogic.CmdClientStatus:
			gameState.CommandStatus()
		case gamelogic.CmdClientHelp:
			gamelogic.PrintClientHelp()
		case gamelogic.CmdClientSpam:
			fmt.Println("Spamming not allowed yet!")
		case "exit":
			fallthrough
		case gamelogic.CmdClientQuit:
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("Error: invalid command '%s'\n", words[0])
		}
	}
}
