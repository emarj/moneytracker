package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ronche.se/moneytracker/handler"
	"ronche.se/moneytracker/model/sqlite"
	"ronche.se/moneytracker/sheet"
)

func main() {

	sheetsSrv, err := sheet.New("client_secret_v1.json", "ux-dSXlaWBUq2itP4QAAuiQW")

	dbSrv, err := sqlite.New("./db.sqlite", true)
	if err != nil {
		log.Fatalln("impossible to connect to db")
	}
	defer func() {
		if err := dbSrv.Close(); err != nil {
			log.Fatalf("impossible to close connection: %v", err)
		}
	}()

	mux, err := handler.HTMLHandler(dbSrv, sheetsSrv, "handler/templates")
	if err != nil {
		log.Fatalf("impossible to create HTMLHandler: %v", err)
	}

	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}
