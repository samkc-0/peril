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

	_, _, err = pubsub.DeclareAndBind(
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

	gameState := gamelogic.NewGameState(username)
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
			if _, err := gameState.CommandMove(words); err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("moved units successfully")
		case gamelogic.CmdClientStatus:
			gameState.CommandStatus()
		case gamelogic.CmdClientHelp:
			gamelogic.PrintClientHelp()
		case gamelogic.CmdClientSpam:
			fmt.Println("Spamming not allowed yet!")
		case gamelogic.CmdClientQuit:
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("Error: invalid command '%s'\n", words[0])
		}
	}
}
