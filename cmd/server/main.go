package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"ronche.se/moneytracker"
	"ronche.se/moneytracker/store/sqlite"
)

var Commit string
var Branch string
var Date string

func main() {
	fmt.Printf("MoneyTracker %s-%s (build %s)\n\n", Branch, Commit, Date)

	var local = flag.Bool("local", false, "")
	var port = flag.Int("port", 3245, "")
	var dir = flag.String("dir", "./data", "")
	var tempDB = flag.Bool("tempDB", false, "")
	var populate = flag.Bool("populate", false, "")
	var dbName = flag.String("db", "moneytracker.sqlite", "")
	flag.Parse()

	hostname := ""
	if *local {
		hostname = "localhost"
	}
	url := fmt.Sprintf("%s:%d", hostname, *port)

	dsn := path.Join(*dir, *dbName)

	if *tempDB {

		f, err := os.CreateTemp("", *dbName)
		if err != nil {
			panic(err)
		}
		dsn = f.Name()
		fmt.Println("created temp db: ", dsn)
		f.Close()
	}

	s := sqlite.New(dsn, true)

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

	if *populate {
		err := s.Seed()
		if err != nil {
			log.Fatal(err)
		}
	}

	srv := moneytracker.NewServer(s)

	go func() {
		if err := srv.Start(url); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server shut down with error: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("\nShutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Fatal(err)
	}

}
