package main

import (
	"fmt"
	"time"

	"github.com/funzoneq/go-radius-dictionaries/erx"
	log "github.com/sirupsen/logrus"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

// const vlan_low = 2
// const vlan_high = 4094

func AuthHandler(w radius.ResponseWriter, r *radius.Request) {
	userInfo := UserInfo{}
	userInfo.startTime = time.Now()

	// Default response is denied.
	resp := r.Response(radius.CodeAccessReject)

	password := rfc2865.UserPassword_GetString(r.Packet)
	username := rfc2865.UserName_GetString(r.Packet)
	userInfo.Identifier = username

	// Match user information
	user := findSubscriber(subscribers, username, &userInfo)

	if password != Config.RadiusSecret {
		log.Printf("Login failed: %s", username)

		writeResponse(w, resp, userInfo)
		return
	} else if user == nil {
		log.Printf("User is not found: %s", username)

		if Config.CaptivePortalEnabled {
			resp = r.Response(radius.CodeAccessAccept)

			// Add Portal VRF
			err := erx.ERXVirtualRouterName_AddString(resp, "default:Portal")
			if err != nil {
				log.Print("Error adding VRF ", err)
			}

			log.Printf("User send to captive portal: %s", username)

			writeResponse(w, resp, userInfo)
			return
		}

		writeResponse(w, resp, userInfo)
		return
	} else if user.Status != "active" {
		log.Printf("User is not active: %s", user.Identifier)

		if Config.CaptivePortalEnabled {
			resp = r.Response(radius.CodeAccessAccept)

			// Add Portal VRF
			err := erx.ERXVirtualRouterName_AddString(resp, "default:Portal")
			if err != nil {
				log.Print("Error adding VRF ", err)
			}

			log.Printf("User send to captive portal: %s", username)

			writeResponse(w, resp, userInfo)
			return
		}

		writeResponse(w, resp, userInfo)
		return
		// TODO: Add check for vlan out of bounds
	} else {
		// Login succeeded
		resp = r.Response(radius.CodeAccessAccept)

		vrf := fmt.Sprintf("default:%s", user.VRF)
		filter_in := fmt.Sprintf("DEFAULT-FILTER-%s", user.SpeedUp)
		filter_out := fmt.Sprintf("DEFAULT-FILTER-%s", user.SpeedDown)

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

		// Add Ingress policy / policer
		err = erx.ERXInputInterfaceFilter_AddString(resp, filter_in)
		if err != nil {
			log.Printf("Error adding ingress interface filter %v", err)
		}
		log.Debugf("Added input filter %v for user %v", filter_in, username)

		// Add Egress policy / policer
		err = erx.ERXOutputInterfaceFilter_AddString(resp, filter_out)
		if err != nil {
			log.Errorf("Error adding egress interface filter %v", err)
		}
		log.Debugf("Added output filter %v for user %v", filter_out, username)

		log.Printf("User authenticated successfully %s", username)

		writeResponse(w, resp, userInfo)
		return
	}
}

func AllowAllAuthHandler(w radius.ResponseWriter, r *radius.Request) {
	userInfo := UserInfo{}
	userInfo.startTime = time.Now()

	password := rfc2865.UserPassword_GetString(r.Packet)
	username := rfc2865.UserName_GetString(r.Packet)

	if password != Config.RadiusSecret {
		log.Printf("Login failed: %s", username)

		writeResponse(w, r.Response(radius.CodeAccessReject), userInfo)
		return
	}

	resp := r.Response(radius.CodeAccessAccept)

	vrf := fmt.Sprintf("default:%s", Config.DefaultVRF)
	filter_in := fmt.Sprintf("DEFAULT-FILTER-%s", Config.DefaultUploadSpeed)
	filter_out := fmt.Sprintf("DEFAULT-FILTER-%s", Config.DefaultDownloadSpeed)

	// Add VRF / Routing Instance
	err := erx.ERXVirtualRouterName_AddString(resp, vrf)
	if err != nil {
		log.Print("Error adding VRF ", err)
	}

	// Add Ingress policy / policer
	err = erx.ERXInputInterfaceFilter_AddString(resp, filter_in)
	if err != nil {
		log.Printf("Error adding ingress interface filter %v", err)
	}
	log.Debugf("Added input filter %v for user %v", filter_in, username)

	// Add Egress policy / policer
	err = erx.ERXOutputInterfaceFilter_AddString(resp, filter_out)
	if err != nil {
		log.Errorf("Error adding egress interface filter %v", err)
	}
	log.Debugf("Added output filter %v for user %v", filter_out, username)

	log.Printf("User authenticated successfully %s", username)

	writeResponse(w, resp, userInfo)
}

func AuthServer(subs []Subscriber) {
	subscribers = subs

	handler := AuthHandler

	if !Config.AuthEnabled {
		handler = AllowAllAuthHandler
	}

	AuthServer := radius.PacketServer{
		Addr:         Config.AuthListenAddress,
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(Config.RadiusSecret)),
	}

	log.Printf("Starting authentication server on %s", Config.AuthListenAddress)
	err := AuthServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
