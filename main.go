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

	sheetsSrv, err := sheet.New("client_secret_v1.json", "1ud3T4uUPOv94Atj4Qopy1qhwatLaXsXnLOl_n-Qxya4")

	dbSrv, err := sqlite.New("./db.sqlite", true)
	if err != nil {
		log.Fatalf("impossible to connect to db: %v", err)
	}
	defer func() {
		if err := dbSrv.Close(); err != nil {
			log.Fatalf("impossible to close db connection: %v", err)
		}
	}()

	mux, err := handler.HTMLHandler(dbSrv, sheetsSrv, "handler/templates")
	if err != nil {
		log.Fatalf("impossible to create HTMLHandler: %v", err)
	}

	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "34567"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}
