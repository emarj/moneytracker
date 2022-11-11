package moneytracker

import (
	_ "embed"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//go:embed embed/debug.html
var debugPage string

type Server struct {
	store  Store
	router *echo.Echo
}

func NewServer(store Store) *Server {

	srv := &Server{store: store, router: echo.New()}
	srv.setup()

	return srv
}

func (s *Server) setup() {
	// Middlewares
	//s.router.Pre(middleware.AddTrailingSlash())
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())

	// Static Routes
	/*ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Dir(ex)
	s.router.Static("/", filepath.Join(basePath, "../static/"))*/

	// API Routes
	/*s.router.GET("/api/users", s.getUsers)
	s.router.GET("/api/account/:aid", s.getAccount)*/
	//e.GET("/accounts/", s.GetAccounts)
	/*s.router.GET("/api/accounts/:uid", s.GetAccountsByUser)
	s.router.GET("/api/transactions/:uid", s.getTransactionsOfUser)
	s.router.GET("/api/transaction/:tid", s.getTransaction)
	s.router.POST("/api/transaction/", s.insertTransaction)*/
	s.router.GET("/debug/", func(c echo.Context) error { return c.HTML(200, debugPage) })

}

func (s *Server) Start(url string) error {
	return s.router.Start(url)
}

//Handlers
/*
func (s *Server) getUsers(c echo.Context) error {
	ul, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ul)
}

func (s *Server) getAccount(c echo.Context) error {
	aid, err := uuid.FromString(c.Param("aid"))
	if err != nil {
		return err
	}
	a, err := s.store.GetAccount(aid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, a)
}

func (s *Server) GetAccountsByUser(c echo.Context) error {
	al, err := s.store.GetAccountsByUser(c.Param("uid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, al)
}

func (s *Server) getTransaction(c echo.Context) error {

	tID, err := uuid.FromString(c.Param("tid"))
	if err != nil {
		return err
	}

	t, err := s.store.GetTransaction(tID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (s *Server) getTransactions(c echo.Context) error {
	aid, err := uuid.FromString(c.Param("aid"))
	if err != nil {
		return err
	}
	tl, err := s.store.GetTransactionsByAccount(aid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}

func (s *Server) getTransactionsOfUser(c echo.Context) error {
	tl, err := s.store.GetTransactionsByUser(c.Param("uid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}

func (s *Server) insertTransaction(c echo.Context) error {
	t := Transaction{}

	err := json.NewDecoder(c.Request().Body).Decode(&t)
	if err != nil {
		return err
	}

	err = s.store.AddTransaction(&t)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}*/
