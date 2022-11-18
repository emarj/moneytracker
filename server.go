package moneytracker

import (
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

	// API Routes
	s.router.GET("/api/entity/:eid", s.getEntity)
	s.router.GET("/api/entities/", s.getEntities)
	s.router.GET("/api/account/:aid", s.getAccount)
	s.router.GET("/api/accounts/:eid", s.getAccountsByEntity)
	s.router.GET("/api/balances/:aid", s.getBalances)
	s.router.POST("/api/balance/", s.addBalance)
	s.router.GET("/api/balance/:aid", s.getBalance)
	//s.router.GET("/api/transactions", s.getTransactions)
	s.router.GET("/api/operations/entity/:eid", s.getOperationsByEntity)
	s.router.GET("/api/transactions/account/:aid", s.getTransactionsByAccount)
	s.router.GET("/api/operation/:opid", s.getOperation)
	s.router.POST("/api/operation/", s.addOperation)
	//s.router.DELETE("/api/transaction/", s.deleteTransaction)

	return s
}

func (s *Server) Start(url string) error {
	return s.router.Start(url)
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
