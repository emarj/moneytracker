package main

import (
	"log"

	"ronche.se/moneytracker"
	"ronche.se/moneytracker/store/sqlite"
)

func main() {

	s := sqlite.New("./db.sqlite", true)

	err := s.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	srv := moneytracker.NewServer(s)

	log.Fatal(srv.Start("localhost:3245"))

}
