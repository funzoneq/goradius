package main

import (
	"os"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/funzoneq/go-radius-dictionaries/erx"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

const vlan_low = 2
const vlan_high = 4094

type UserIdentifier struct {
	Router string
	Interface string
	Outertag int
	Innertag int
}

func matchUserIdentifier(Username string) UserIdentifier {
	// Regex edgegw01-somesite.ps36:2017-10001
	valid_username := regexp.MustCompile(`(?P<router>[\w-]+)\.(?P<intf>ps\d+)\:(?P<outertag>\d+)\-(?P<innertag>\d+)`)
	res := valid_username.FindStringSubmatch(Username)

	outertag, _ := strconv.Atoi(res[3])
	innertag, _ := strconv.Atoi(res[4])

	user := UserIdentifier{
		Router: res[1],
		Interface: res[2],
		Outertag: outertag,
		Innertag: innertag,
	}

	return user
}

func AuthHandler(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	// Default response is denied.
	resp := r.Response(radius.CodeAccessReject)

	// Match user information
	user := matchUserIdentifier(username)

	if password != os.Getenv("RADIUS_SECRET") {
		log.Printf("Login failed: %s", username)
		err := w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
	} else if (user.Innertag < vlan_low || user.Innertag > vlan_high || user.Outertag < vlan_low || user.Outertag > vlan_high) {
		log.Printf("Vlan out of bounds: %d-%d", user.Outertag, user.Innertag)
		err := w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Login succeeded
		resp = r.Response(radius.CodeAccessAccept)

		vrf := "default:BNG-Users"
		v4_filter_in := "DEFAULT-FILTER-ACCEPT-V4"
		v4_filter_out := "DEFAULT-FILTER-ACCEPT-V4"
		v6_filter_in := "DEFAULT-FILTER-ACCEPT-V6"
		v6_filter_out := "DEFAULT-FILTER-ACCEPT-V6"

		// Add VRF / Routing Instance
		err := erx.ERXVirtualRouterName_AddString(resp, vrf)
		if err != nil {
			log.Print("Error adding VRF ", err)
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

func AuthServer() {
	AuthServer := radius.PacketServer{
		Addr: 		  ":1812",
		Handler:      radius.HandlerFunc(AuthHandler),
		SecretSource: radius.StaticSecretSource([]byte(os.Getenv("RADIUS_SECRET"))),
	}

	log.Printf("Starting authentication server on :1812")
	err := AuthServer.ListenAndServe();
	if err != nil {
		log.Fatal(err)
	}
}