package moneytracker

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//go:embed frontend/dist/*
var content embed.FS

type Server struct {
	store  Store
	router *echo.Echo
}

func NewServer(store Store) *Server {

	s := &Server{store: store, router: echo.New()}

	s.router.HideBanner = true

	// Middlewares
	//s.router.Pre(middleware.AddTrailingSlash())
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())

	s.router.Use(middleware.CORS())

	// Static Routes

	var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))
	// The embedded files will all be in the '/frontend/dist/' folder so need to rewrite the request (could also do this with fs.Sub)
	var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/frontend/dist/$1"})

	s.router.GET("/*", contentHandler, contentRewrite)

	apiGroup := s.router.Group("api")

	// API Routes
	apiGroup.GET("/entity/:eid", s.getEntity)
	apiGroup.GET("/entities", s.getEntities)
	apiGroup.GET("/account/:aid", s.getAccount)
	apiGroup.GET("/accounts", s.getAccounts)
	apiGroup.GET("/accounts/:eid", s.getAccountsByEntity)
	apiGroup.GET("/balances/:aid", s.getBalances)
	apiGroup.POST("/balance/", s.addBalance)
	apiGroup.GET("/balance/:aid", s.getBalance)
	//apiGroup.GET("/transactions", s.getTransactions)
	apiGroup.GET("/operations/entity/:eid", s.getOperationsByEntity)
	apiGroup.GET("/transactions/account/:aid", s.getTransactionsByAccount)
	apiGroup.GET("/operation/:opid", s.getOperation)
	apiGroup.POST("/operation/", s.addOperation)
	//apiGroup.DELETE("/transaction/", s.deleteTransaction)

	return s
}

func (s *Server) Start(url string) error {
	return s.router.Start(url)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}

// Handlers
func (s *Server) getEntity(c echo.Context) error {

	eID, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		return err
	}

	e, err := s.store.GetEntity(eID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, e)
}

func (s *Server) getEntities(c echo.Context) error {
	el, err := s.store.GetEntities()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, el)
}

func (s *Server) getAccount(c echo.Context) error {

	aID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		return err
	}

	a, err := s.store.GetAccount(aID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, a)
}

func (s *Server) getAccounts(c echo.Context) error {

	el, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, el)
}

func (s *Server) getAccountsByEntity(c echo.Context) error {
	eID, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		return err
	}
	el, err := s.store.GetAccountsByEntity(eID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, el)
}

func (s *Server) getBalances(c echo.Context) error {
	aID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		return err
	}
	bl, err := s.store.GetBalances(aID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bl)
}

func (s *Server) getBalance(c echo.Context) error {

	aID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		return err
	}

	b, err := s.store.GetBalance(aID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, b)
}

func (s *Server) addBalance(c echo.Context) error {
	b := Balance{}

	err := json.NewDecoder(c.Request().Body).Decode(&b)
	if err != nil {
		return err
	}

	if b.Value != nil {
		err = s.store.ComputeBalance(b.AccountID)
	} else {
		err = s.store.AddBalance(b)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

/*func (s *Server) getTransactions(c echo.Context) error {
	tl, err := s.store.GetTransactions()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}*/

func (s *Server) getTransactionsByAccount(c echo.Context) error {
	aID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		return err
	}

	tl, err := s.store.GetTransactionsByAccount(aID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}

func (s *Server) getOperationsByEntity(c echo.Context) error {
	aID, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		return err
	}

	list, err := s.store.GetOperationsByEntity(aID)
	if err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, list, "\t")
}

func (s *Server) getOperation(c echo.Context) error {

	opID, err := strconv.Atoi(c.Param("opid"))
	if err != nil {
		return err
	}

	op, err := s.store.GetOperation(opID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, op)
}

func (s *Server) addOperation(c echo.Context) error {
	op := Operation{}

	err := json.NewDecoder(c.Request().Body).Decode(&op)
	if err != nil {
		return err
	}

	id, err := s.store.AddOperation(op)
	if err != nil {
		return err
	}
	//do not return t since it might be incomplete
	return c.JSON(http.StatusOK, id)
}
