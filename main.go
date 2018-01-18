package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ronche.se/moneytracker/handler"
	"ronche.se/moneytracker/model/sqlite"
)

func main() {

	srv, err := sqlite.New("./db.sqlite", true)
	if err != nil {
		log.Fatalln("impossible to connect to db")
	}
	defer func() {
		if err := srv.Close(); err != nil {
			log.Fatalf("impossible to close connection: %v", err)
		}
	}()

	mux, err := handler.HTMLHandler(srv, "handler/templates")
	if err != nil {
		log.Fatalf("impossible to create HTMLHandler: %v", err)
	}

	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, mux))

}
