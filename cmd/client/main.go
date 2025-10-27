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

	fmt.Println("Starting Peril client...")

	gameState := gamelogic.NewGameState(username)

	// subscribe to pause events
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

	// subscribe to other players' move events
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.TransientQueueType,
		func(move gamelogic.ArmyMove) pubsub.AckType {
			defer fmt.Print("> ")
			switch gameState.HandleMove(move) {
			case gamelogic.MoveOutcomeSafe:
				fallthrough
			case gamelogic.MoveOutcomeMakeWar:
				ch, err := conn.Channel()
				defer ch.Close()
				if err != nil {
					fmt.Println(err)
				}
				if err := pubsub.PublishJSON(
					ch,
					routing.ExchangePerilTopic,
					routing.WarRecognitionPrefix+"."+username,
					move,
				); err != nil {
					fmt.Println(err)
				}
				return pubsub.NackRequeue
			case gamelogic.MoveOutcomeSamePlayer:
				fallthrough
			default:
				return pubsub.NackDiscard
			}
		},
	); err != nil {
		fmt.Println(err)
	}

	// subscribe to all war outcomes
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.WarRecognitionPrefix,
		routing.WarRecognitionPrefix+".*",
		pubsub.DurableQueueType,
		handlerWar(gameState),
	); err != nil {
		fmt.Println(err)
	}

	moveChan, err := conn.Channel()
	if err != nil {
		fmt.Println("error creating channel:")
		log.Fatal(err)
	}

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
				moveChan,
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
