package main

import (
	"fmt"
	"peril/internal/gamelogic"
	"peril/internal/pubsub"
	"peril/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(
	gs *gamelogic.GameState,
	ch *amqp.Channel,
) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		switch gs.HandleMove(move) {
		case gamelogic.MoveOutcomeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			recognitionOfWar := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			if err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionPrefix+"."+gs.GetUsername(),
				recognitionOfWar,
			); err != nil {
				fmt.Println(err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			fmt.Println("unknown move outcome")
			return pubsub.NackDiscard
		}
	}
}
