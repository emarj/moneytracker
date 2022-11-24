package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	go func() {
		if err := srv.Start("localhost:3245"); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
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
