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
	dbPath := os.Getenv("DBPATH")
	if dbPath == "" {
		dbPath = "../moneytracker.sqlite"
		fmt.Println("INFO: No DBPATH environment variable detected, defaulting to " + dbPath)
	}

	_, err = os.Open(dbPath)
	if err != nil {
		log.Fatalf("impossible to open the db file: %v", err)
	}
	dbSrv, err := sqlite.New(dbPath, false)
	if err != nil {
		log.Fatalf("impossible to connect to db: %v", err)
	}
	defer func() {
		if err := dbSrv.Close(); err != nil {
			log.Fatalf("impossible to close db connection: %v", err)
		}
	}()

	prefix := os.Getenv("PREFIX")

	mux, err := handler.HTMLHandler(dbSrv, sheetsSrv, "handler/templates", prefix)
	if err != nil {
		log.Fatalf("impossible to create HTMLHandler: %v", err)
	}

	port := os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "34567"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}
