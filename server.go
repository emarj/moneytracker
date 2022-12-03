package moneytracker

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

//go:embed frontend/dist/*
var content embed.FS

// Auth

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	User
	jwt.StandardClaims
}

var secret_key = []byte("super_secret_key")

type Server struct {
	store  Store
	router *echo.Echo
}

func NewServer(store Store) *Server {

	decimal.MarshalJSONWithoutQuotes = true

	s := &Server{store: store, router: echo.New()}

	s.router.HideBanner = true

	// Middlewares
	//s.router.Pre(middleware.AddTrailingSlash())
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())

	s.router.Use(middleware.CORS())

	// Frontend

	proxyURL := os.Getenv("MT_FRONTEND_URL")
	if proxyURL != "" {
		frontendURL, err := url.Parse(proxyURL)
		if err != nil {
			s.router.Logger.Fatal(err)
		}

		s.router.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
			Skipper: func(c echo.Context) bool {
				p := strings.Split(c.Path()[1:], "/")
				return p[0] == "api"
			},
			Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: frontendURL}}),
		}))
	} else {
		var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))
		// The embedded files will all be in the '/frontend/dist/' folder so need to rewrite the request (could also do this with fs.Sub)
		var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/frontend/dist/$1"})

		s.router.GET("/*", contentHandler, contentRewrite)
	}

	apiGroup := s.router.Group("api")

	config := middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/login"
		},
		ErrorHandler: func(err error) error {
			return err
		},
		SigningKey:  secret_key,
		Claims:      &jwtCustomClaims{},
		TokenLookup: "cookie:auth",
		ContextKey:  "auth",
	}
	apiGroup.Use(middleware.JWTWithConfig(config))

	// API Routes
	apiGroup.GET("/echo", s.Echo)
	apiGroup.POST("/login", s.Login)
	apiGroup.POST("/logout", s.Logout)
	apiGroup.GET("/entity/:eid", s.getEntity)
	apiGroup.GET("/entities", s.getEntities)
	apiGroup.GET("/account/:aid", s.getAccount)
	apiGroup.GET("/accounts", s.getAccounts)
	apiGroup.GET("/accounts/:eid", s.getAccountsByEntity)
	apiGroup.GET("/balances/:aid", s.getBalances)
	apiGroup.POST("/balance", s.addBalance)
	apiGroup.GET("/balance/:aid", s.getBalance)
	//apiGroup.GET("/transactions", s.getTransactions)
	apiGroup.GET("/operations/entity/:eid", s.getOperationsByEntity)
	apiGroup.GET("/transactions/account/:aid", s.getTransactionsByAccount)
	apiGroup.GET("/operation/:opid", s.getOperation)
	apiGroup.DELETE("/operation/:opid", s.deleteOperation)
	apiGroup.POST("/operation", s.addOperation)

	apiGroup.GET("/categories", s.getCategories)

	return s
}

func (s *Server) Start(url string) error {
	return s.router.Start(url)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}

// Auth

func newAuthCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "auth",
		Secure:   false, // This should be set to true as soon as we implement HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (s *Server) Login(c echo.Context) error {

	login := struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(c.Request().Body).Decode(&login)
	if err != nil {
		return err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(login.Password), 14)
	if err != nil {
		return err
	}

	ok, err := s.store.Login(login.User, passHash)
	if err != nil {
		return err
	}
	if !ok {
		return echo.ErrUnauthorized
	}

	// Valid login

	expiration := time.Now().Add(time.Hour * 72)

	// Set custom claims
	claims := &jwtCustomClaims{
		User{Name: login.User, Admin: true},
		jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(secret_key)
	if err != nil {
		return err
	}

	cookie := newAuthCookie()
	cookie.Value = t
	cookie.Expires = expiration

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "logged in",
	})
}

func (s *Server) Logout(c echo.Context) error {

	cookie := newAuthCookie()
	cookie.MaxAge = -1

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "logged out",
	})
}

func (s *Server) Echo(c echo.Context) error {

	authRaw := c.Get("auth")
	if authRaw == nil {
		return echo.ErrInternalServerError
	}
	auth, ok := authRaw.(*jwt.Token)
	if !ok {
		return echo.ErrInternalServerError
	}
	claims, ok := auth.Claims.(*jwtCustomClaims)
	if !ok {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Hi " + claims.Name,
	})
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

func (s *Server) getCategories(c echo.Context) error {
	cl, err := s.store.GetCategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cl)
}

func (s *Server) deleteOperation(c echo.Context) error {

	opID, err := strconv.Atoi(c.Param("opid"))
	if err != nil {
		return err
	}

	err = s.store.DeleteOperation(opID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
