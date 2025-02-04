package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	configFile = "/etc/goradius/goradius.conf"
	debug      = false
	Config     = &GoradiusConfig{} // Holds the GoRADIUS configuration data
)

var subscribers []Subscriber

func main() {
	flag.StringVar(&configFile, "config", "/etc/goradius/goradius.conf", "Path to GoRADIUS config file")
	flag.BoolVar(&debug, "debug", false, "Show debug logging")

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stdout)

	// Load config from file
	Config.ReadConfig(configFile)

	// Load subscribers from file
	subscribers, err := loadSubscribers(Config.CustomerFile)
	if err != nil {
		log.Errorf("Error loading JSON: %v", err)
		return
	}

	go AccountingServer()
	AuthServer(subscribers)
}
