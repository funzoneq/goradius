package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	// Load subscribers from file
	subscribers, err := loadSubscribers("customers.json")
	if err != nil {
		fmt.Println("Error loading JSON:", err)
		return
	}

	go AccountingServer()
	AuthServer(subscribers)
}
