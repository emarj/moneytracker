package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"ronche.se/moneytracker"
	"ronche.se/moneytracker/db"
	"ronche.se/moneytracker/db/mock"
)

func main() {

	s := mock.NewMockStore()
	db.Populate(s)

	handler := moneytracker.NewHandler(s)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", handler.GetUsers)

	e.GET("/account/:aid", handler.GetAccount)
	//e.GET("/accounts/", handler.GetAccounts)
	e.GET("/accounts/:uid", handler.GetAccountsOfUser)
	e.GET("/transactions/:uid", handler.GetTransactionsOfUser)
	e.GET("/transaction/:tid", handler.GetTransaction)
	e.POST("/transaction/", handler.InsertTransaction)

	// Start server
	e.Logger.Fatal(e.Start("localhost:1323"))
}
