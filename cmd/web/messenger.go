package main

import (
	"fmt"
	"log"

	"github.com/gregdel/pushover"
)

func SendMessage(userId string, appId string, hashValue string, nonceValue string, timeToSolve float64) {
	app := pushover.New(userId)
	recipient := pushover.NewRecipient(appId)
	message := pushover.NewMessage("Execution Time: " + fmt.Sprintf("%f", timeToSolve) + " Hash Value: " + hashValue + " Nonce Value: " + nonceValue)
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.Println("Response:", response)

}
