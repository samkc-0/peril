package gamelogic

import (
	"fmt"
	"log"
	"os"
	"peril/internal/routing"
	"time"
)

const logsFile = "game.log"

const writeToDiskSleep = 1 * time.Second

func WriteLog(gamelog routing.GameLog) error {
	log.Printf("received game log...")
	time.Sleep(writeToDiskSleep)

	f, err := os.OpenFile(logsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open logs file: %w", err)
	}
	defer f.Close()

	logEntry := fmt.Sprintf("%v %v: %v\n", gamelog.CurrentTime.Format(time.RFC3339), gamelog.Username, gamelog.Message)
	_, err = f.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("could not write to logs file: %w", err)
	}
	return nil
}
