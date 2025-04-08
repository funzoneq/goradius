package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type GoradiusConfig struct {
	RadiusSecret         string
	AuthListenAddress    string
	AcctListenAddress    string
	MetricsListenAddress string
	CustomerFile         string
	CaptivePortalEnabled bool
	AuthEnabled 		 bool
	DefaultVRF			 string
	DefaultUploadSpeed	 string
	DefaultDownloadSpeed string
}

// ReadConfig reads and parses the RADIUS configuration from JSON format from the given filepath, and fails on invalid configuration values.
func (c *GoradiusConfig) ReadConfig(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Errorf("Could not open GoRADIUS config file, will use defaults %v", err)
		return
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		log.Fatalf("Problem parsing GoRADIUS config file %v", err)
	}

	if len(c.RadiusSecret) == 0 {
		log.Fatal("RADIUS secret is empty!")
	}
}
