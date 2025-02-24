package main

import (
	"time"

	"github.com/funzoneq/go-radius-dictionaries/erx"
	log "github.com/sirupsen/logrus"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)

// Constants
const (
	acctStatisticsMode = 2    // 0 = disable; 1 = track time only; 2 = track traffic volume and time
	acctUpdateInterval = 3600 // Update interval in seconds
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
	userInfo := UserInfo{}
	userInfo.startTime = time.Now()
	DebugRadiusPacket(r)

	if r.Code != radius.CodeAccountingRequest {
		return
	}

	log.Printf("Accounting packet received")

	// Default response
	resp := r.Response(radius.CodeAccountingResponse)

	username := rfc2865.UserName_GetString(r.Packet)
	userInfo.Identifier = username
	acctID := rfc2866.AcctSessionID_GetString(r.Packet)
	accountingStatus := rfc2866.AcctStatusType_Get(r.Packet)

	res, err := ParseUsername(username)
	if err != nil {
		log.Errorf("Username parsing error: %v", err)
	} else {
		userInfo.site = res[1]
	}

	switch accountingStatus {
	// These types of packets get sent on starting/stopping the RADIUS service on the BNG
	case rfc2866.AcctStatusType_Value_AccountingOn:
		log.Debugf("Got Accounting-On accounting packet from %v", rfc2865.NASIdentifier_GetString(r.Packet))
		return
	case rfc2866.AcctStatusType_Value_AccountingOff:
		log.Debugf("Got Accounting-Off accounting packet from %v", rfc2865.NASIdentifier_GetString(r.Packet))
		return
	case rfc2866.AcctStatusType_Value_Failed:
		log.Errorf("Got Failed accounting packet from %v", rfc2865.NASIdentifier_GetString(r.Packet))
		return
	default:
		log.Debugf("Got accounting packet with type %v from %v", rfc2866.AcctStatusType_Strings[accountingStatus], rfc2865.NASIdentifier_GetString(r.Packet))
	}

	// Handling per-subscriber packets
	switch accountingStatus {
	case rfc2866.AcctStatusType_Value_Start:
		log.Infof("Start accounting for %v acctID %v", username, acctID)
	case rfc2866.AcctStatusType_Value_Stop:
		log.Infof("Stop accounting for %v acctID %v cause: %v", rfc2865.UserName_GetString(r.Packet), rfc2866.AcctSessionID_GetString(r.Packet), rfc2866.AcctTerminateCause_Get(r.Packet))
	case rfc2866.AcctStatusType_Value_InterimUpdate:
		log.Infof("Interim update for %v acctID %v", username, acctID)
	default:
		log.Debugf("Got accounting packet with type %v from %v", rfc2866.AcctStatusType_Strings[accountingStatus], rfc2865.NASIdentifier_GetString(r.Packet))
	}

	// Tag 0 (0x00) is the statistics mode
	err = erx.ERXServiceStatistics_Add(resp, 0x00, acctStatisticsMode)
	if err != nil {
		log.Errorf("Could not add Service Statistics: %v", err)
	}

	// Tag 0 (0x00) is the accounting update interval
	err = erx.ERXServiceAcctInterval_Add(resp, 0x00, acctUpdateInterval)
	if err != nil {
		log.Errorf("Could not add Service Accounting Interval: %v", err)
	}

	writeResponse(w, resp, userInfo)
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
