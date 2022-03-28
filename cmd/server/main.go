package main

import (
	"log"

	"ronche.se/moneytracker"
	"ronche.se/moneytracker/db"
	"ronche.se/moneytracker/db/mock"
)

const apiKey = "keyAuR8F3wLAUXZAL"

func main() {

	s := mock.NewMockStore()

	db.Populate(s)

	srv := moneytracker.NewServer(s)

	log.Fatal(srv.Start("localhost:3245"))

}
