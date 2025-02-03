package main

import (
	"encoding/json"
	"os"
)

type Subscriber struct {
	Identifier   string   `json:"identifier"`
	Status       string   `json:"status"`
	SpeedUp      string   `json:"speed_up"`
	SpeedDown    string   `json:"speed_down"`
	VRF          string   `json:"vrf"`
	StaticRoutes []string `json:"static_routes"`
}

// Load subscribers from JSON file
func loadSubscribers(filename string) ([]Subscriber, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var subscribers []Subscriber
	err = json.Unmarshal(data, &subscribers)
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

// Find a subscriber by identifier
func findSubscriber(subscribers []Subscriber, identifier string) *Subscriber {
	for _, sub := range subscribers {
		if sub.Identifier == identifier {
			return &sub
		}
	}
	return nil
}
