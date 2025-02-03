package main

import (
	"fmt"
	"os"

	"github.com/funzoneq/go-radius-dictionaries/erx"
	log "github.com/sirupsen/logrus"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

// const vlan_low = 2
// const vlan_high = 4094
var subscribers []Subscriber

func AuthHandler(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	// Default response is denied.
	resp := r.Response(radius.CodeAccessReject)

	// Match user information
	user := findSubscriber(subscribers, username)

	if password != os.Getenv("RADIUS_SECRET") {
		log.Printf("Login failed: %s", username)
		err := w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
	} else if user == nil {
		log.Printf("User is not found: %s", username)
		err := w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
	} else if user.Status != "active" {
		log.Printf("User is not active: %s", user.Identifier)
		err := w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: Add check for vlan out of bounds
	} else {
		// Login succeeded
		resp = r.Response(radius.CodeAccessAccept)

		vrf := fmt.Sprintf("default:%s", user.VRF)
		v4_filter_in := "DEFAULT-FILTER-ACCEPT-V4"
		v4_filter_out := "DEFAULT-FILTER-ACCEPT-V4"
		v6_filter_in := "DEFAULT-FILTER-ACCEPT-V6"
		v6_filter_out := "DEFAULT-FILTER-ACCEPT-V6"

		// Add VRF / Routing Instance
		err := erx.ERXVirtualRouterName_AddString(resp, vrf)
		if err != nil {
			log.Print("Error adding VRF ", err)
		}

		for _, staticRoute := range user.StaticRoutes {
			err = rfc2865.FramedRoute_AddString(resp, staticRoute)
			if err != nil {
				log.Errorf("Error adding static route %v", err)
			}
			log.Debugf("Added static route %v for user %v", staticRoute, username)
		}

		// Add IPv6 Ingress policy / policer
		err = erx.ERXIPv6IngressPolicyName_AddString(resp, v6_filter_in)
		if err != nil {
			log.Print("Error adding IPv6 ingress ", err)
		}

		// Add IPv6 Egress policy / policer
		err = erx.ERXIPv6EgressPolicyName_AddString(resp, v6_filter_out)
		if err != nil {
			log.Print("Error adding IPv6 egress ", err)
		}

		// Add IPv4 Ingress policy / policer
		err = erx.ERXIngressPolicyName_AddString(resp, v4_filter_in)
		if err != nil {
			log.Print("Error adding IPv4 ingress ", err)
		}

		// Add IPv4 Egress policy / policer
		err = erx.ERXEgressPolicyName_AddString(resp, v4_filter_out)
		if err != nil {
			log.Print("Error adding IPv4 egress ", err)
		}

		log.Printf("User authenticated successfully %s", username)

		err = w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func AuthServer(subs []Subscriber) {
	subscribers = subs

	AuthServer := radius.PacketServer{
		Addr:         ":1812",
		Handler:      radius.HandlerFunc(AuthHandler),
		SecretSource: radius.StaticSecretSource([]byte(os.Getenv("RADIUS_SECRET"))),
	}

	log.Printf("Starting authentication server on :1812")
	err := AuthServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
