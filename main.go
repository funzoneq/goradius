package main

import (
	"log"
	"flag"
	"os"
)

var logfile = "/var/log/goradius.log"

func main() {
	flag.StringVar(&logfile, "logfile", "/var/log/goradius.log", "File to log to")
	flag.Parse()

	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	go AccountingServer()
	AuthServer()
}