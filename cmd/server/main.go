package main

import (
	"log"

	"ronche.se/moneytracker"
)

func main() {

	//s := mock.NewMockStore()

	//db.Populate(s)

	srv := moneytracker.NewServer(nil)

	log.Fatal(srv.Start("localhost:3245"))

}
