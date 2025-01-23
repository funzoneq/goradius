package main

import (
	"os"
	"log"

	"layeh.com/radius"
)

func AccountingHandler(w radius.ResponseWriter, r *radius.Request) {
	log.Printf("Accounting packet received")
}

func AccountingServer() {
	AccountingServer := radius.PacketServer{
		Addr: 		  ":1813",
		Handler:      radius.HandlerFunc(AccountingHandler),
		SecretSource: radius.StaticSecretSource([]byte(os.Getenv("RADIUS_SECRET"))),
	}

	log.Printf("Starting accounting server on :1813")
	err := AccountingServer.ListenAndServe();
	if err != nil {
		log.Fatal(err)
	}
}