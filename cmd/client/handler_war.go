package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"peril/internal/gamelogic"
	"peril/internal/pubsub"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(rw)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			fallthrough
		case gamelogic.WarOutcomeYouWon:
			message := getVictoryMessage(winner, loser)
			if err := pubsub.PublishGameLog(ch, message, gs.GetUsername()); err != nil {
				fmt.Println(err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			message := getDrawMessage(winner, loser)
			if err := pubsub.PublishGameLog(ch, message, gs.GetUsername()); err != nil {
				fmt.Println(err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Println("unknown war outcome")
			return pubsub.NackDiscard
		}
	}
}

func getVictoryMessage(winner, loser string) string {
	result := fmt.Sprintf("%s won a war against %s", winner, loser)
	return result
}

func getDrawMessage(winner, loser string) string {
	result := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
	return result
}
