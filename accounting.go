package main

import (
	"github.com/funzoneq/go-radius-dictionaries/erx"
	log "github.com/sirupsen/logrus"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)

func DebugRadiusPacket(r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	acctID := rfc2866.AcctSessionID_GetString(r.Packet)
	log.Debugf("Received code: %v for username: %v accounting ID: %v", r.Code, username, acctID)

	v4InputOctets := rfc2866.AcctInputOctets_Get(r.Packet)
	v4OutputOctets := rfc2866.AcctOutputOctets_Get(r.Packet)
	v4InputPackets := rfc2866.AcctInputPackets_Get(r.Packet)
	v4OutputPackets := rfc2866.AcctOutputPackets_Get(r.Packet)

	log.Infof("%v-%v: IPv4 Octets in: %d Octets out: %d Packets in: %d Packets out: %d", username, acctID, v4InputOctets, v4OutputOctets, v4InputPackets, v4OutputPackets)

	v6InputOctets := erx.ERXIPv6AcctInputOctets_Get(r.Packet)
	v6OutputOctets := erx.ERXIPv6AcctOutputOctets_Get(r.Packet)
	v6InputPackets := erx.ERXIPv6AcctInputPackets_Get(r.Packet)
	v6OutputPackets := erx.ERXIPv6AcctOutputPackets_Get(r.Packet)

	log.Infof("%v-%v: IPv6 Octets in: %d Octets out: %d Packets in: %d Packets out: %d", username, acctID, v6InputOctets, v6OutputOctets, v6InputPackets, v6OutputPackets)

	for i, attr := range r.Packet.Attributes {
		log.Debugf("Received attribute: type: %v attribute: %v i: %v", attr.Type, attr.Attribute, i)
	}
}

func AccountingHandler(w radius.ResponseWriter, r *radius.Request) {
	DebugRadiusPacket(r)

	if r.Code != radius.CodeAccountingRequest {
		return
	}

	log.Printf("Accounting packet received")
}

func AccountingServer() {
	AccountingServer := radius.PacketServer{
		Addr:         Config.AcctListenAddress,
		Handler:      radius.HandlerFunc(AccountingHandler),
		SecretSource: radius.StaticSecretSource([]byte(Config.RadiusSecret)),
	}

	log.Printf("Starting accounting server on %s", Config.AcctListenAddress)
	err := AccountingServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
