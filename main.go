package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	ConfigFile = "/etc/goradius/goradius.conf"
	Debug      = false
	Config     = &GoradiusConfig{} // Holds the GoRADIUS configuration data
)

var subscribers []Subscriber

func main() {
	flag.StringVar(&ConfigFile, "config", "/etc/goradius/goradius.conf", "Path to GoRADIUS config file")
	flag.BoolVar(&Debug, "debug", false, "Show debug logging")

	flag.Parse()

	if Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stdout)

	// Load config from file
	Config.ReadConfig(ConfigFile)

	// Load subscribers from file
	subscribers, err := loadSubscribers(Config.CustomerFile)
	if err != nil {
		log.Errorf("Error loading JSON: %v", err)
		return
	}

	go AccountingServer()
	AuthServer(subscribers)
}
