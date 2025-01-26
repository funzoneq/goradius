package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	go AccountingServer()
	AuthServer()
}