package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"layeh.com/radius"
)

type UserInfo struct {
	Identifier string
	site       string
	startTime  time.Time
}

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
func findSubscriber(subscribers []Subscriber, identifier string, userInfo *UserInfo) *Subscriber {
	username := ""
	res, err := ParseUsername(identifier)
	if err != nil {
		log.Errorf("Username parsing error: %v", err)
	} else {
		username = res[0]
		userInfo.site = res[1]
	}

	for _, sub := range subscribers {
		if sub.Identifier == username {
			return &sub
		}
	}
	return nil
}

func ParseUsername(username string) ([]string, error) {
	// Regex vallejo.ps1:22-23@dhcpv4
	validUsername := regexp.MustCompile(`(?P<routerID>[\w]+)\.(?P<intf>[\w\/-]+)\:(?:(?P<svlanID>\d{1,4})-)?(?P<cvlanID>\d{1,4})`)
	matches := validUsername.FindStringSubmatch(username)
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not match username with regex: %v", username)
	} else if len(matches) != len(validUsername.SubexpNames()) {
		return nil, fmt.Errorf("could not fully parse username: %v, only matched on %d out of %d", username, len(matches), len(validUsername.SubexpNames()))
	}
	return matches, nil
}

func writeResponse(w radius.ResponseWriter, resp *radius.Packet, userInfo UserInfo) {
	err := w.Write(resp)
	if err != nil {
		log.Fatal(err)
	}

	// Save metrics for request
	labels := prometheus.Labels{"identifier": userInfo.Identifier, "site": userInfo.site}

	switch resp.Code {
	case radius.CodeAccessAccept:
		authAcceptedRequests.With(labels).Add(1)
		authDuration.WithLabelValues("accept").Observe(time.Since(userInfo.startTime).Seconds())

	case radius.CodeAccessReject:
		authRejectedRequests.With(labels).Add(1)
		authDuration.WithLabelValues("reject").Observe(time.Since(userInfo.startTime).Seconds())

	case radius.CodeAccountingResponse:
		acctRequests.With(labels).Add(1)
		acctDuration.Observe(time.Since(userInfo.startTime).Seconds())
	}
}
