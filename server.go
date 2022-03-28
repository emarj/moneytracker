package moneytracker

import (
	_ "embed"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"ronche.se/moneytracker/db"
)

//go:embed embed/debug.html
var debugPage string

type Server struct {
	store  db.Store
	router *echo.Echo
}

func NewServer(store db.Store) *Server {

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
	s.router.GET("/api/users", s.getUsers)
	s.router.GET("/api/account/:aid", s.getAccount)
	//e.GET("/accounts/", s.GetAccounts)
	s.router.GET("/api/accounts/:uid", s.getAccountsOfUser)
	s.router.GET("/api/transactions/:uid", s.getTransactionsOfUser)
	s.router.GET("/api/transaction/:tid", s.getTransaction)
	s.router.POST("/api/transaction/", s.insertTransaction)
	s.router.GET("/debug/", func(c echo.Context) error { return c.HTML(200, debugPage) })

}

func (s *Server) Start(url string) error {
	return s.router.Start(url)
}
