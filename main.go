package main

import (
	"fmt"
	"log"
	"net/http"

	"ronche.se/expensetracker/handler"
	"ronche.se/expensetracker/model/sqlite"
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

	srv.UserInsert("M")
	srv.UserInsert("A")

	srv.CategoryInsert("Uncategorized")
	srv.CategoryInsert("Spesa")
	srv.CategoryInsert("Ristorante")

	srv.PaymentMethodInsert("Contanti")
	srv.PaymentMethodInsert("Bancomat")
	srv.PaymentMethodInsert("CC / Paypal")

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe("localhost:3000", mux))

}
