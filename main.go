package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"ronche.se/moneytracker/handler"
	"ronche.se/moneytracker/model/sqlite"
)

func main() {

	//Commandline args
	dbPath := flag.String("dbpath", "./moneytracker.sqlite", "sqlite3 db path relative to the executable")
	dbCreate := flag.Bool("dbcreate", false, "if true a sqlite3 db is created")
	prefix := flag.String("prefix", "", "prefix to use behind a reverse proxy (e.g. /prefix)")
	port := flag.Int("port", 34567, "port number")
	address := flag.String("address", "", "bind address")
	tmplPath := flag.String("tmplpath", "../handler/templates", "template directory path")

	flag.Parse()

	//Get executable path
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)

	dbSrv, err := sqlite.New(filepath.Join(exPath, *dbPath), *dbCreate)
	if err != nil {
		log.Fatalf("impossible to connect to db: %v", err)
	}
	defer func() {
		if err := dbSrv.Close(); err != nil {
			log.Fatalf("impossible to close db connection: %v", err)
		}
	}()

	mux, err := handler.HTMLHandler(dbSrv, filepath.Join(exPath, *tmplPath), *prefix)
	if err != nil {
		log.Fatalf("impossible to create HTMLHandler: %v", err)
	}
	fullAddr := *address + ":" + strconv.Itoa(*port)
	fmt.Printf("Listening and serving %s...\n", fullAddr+*prefix)
	log.Fatal(http.ListenAndServe(fullAddr, mux))

}
