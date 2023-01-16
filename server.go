package moneytracker

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shopspring/decimal"
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

// FIXME: change this
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

	//s.router.Pre(middleware.AddTrailingSlash()) This causes a mess with nested routes

	s.router.Use(middleware.Recover())

	s.router.Use(middleware.CORS())

	// Frontend

	proxyURL := os.Getenv("MT_FRONTEND_URL")
	if proxyURL != "" {
		fmt.Printf("\nFrontend proxy mode %s\n", proxyURL)
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
	apiGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_unix_micro}: method=${method}, uri=${uri}, status=${status}\t error=${error} \n",
	}))

	config := middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/login" || c.Path() == "/api/logout"
		},
		ErrorHandler: func(err error) error {
			return err
		},
		SigningKey:  secret_key,
		Claims:      &jwtCustomClaims{},
		TokenLookup: "cookie:token",
		ContextKey:  "token",
	}
	apiGroup.Use(middleware.JWTWithConfig(config))

	// API Routes
	apiGroup.GET("/greet", s.Greet)
	apiGroup.POST("/login", s.Login)
	apiGroup.POST("/logout", s.Logout)

	apiGroup.GET("/types", s.getTypes)

	apiGroup.GET("/entities", s.getEntities)
	apiGroup.GET("/entity/:eid", s.getEntity)

	apiGroup.GET("/accounts", s.getAccounts)
	apiGroup.GET("/accounts/:eid", s.getAccountsByEntity)
	apiGroup.GET("/account/:aid", s.getAccount)
	apiGroup.POST("/account", s.addAccount)
	apiGroup.DELETE("/account/:aid", s.deleteAccount)

	apiGroup.GET("/balance/:aid", s.getBalance)
	apiGroup.GET("/balance/history/:aid", s.getBalanceHistory)
	apiGroup.POST("/balance", s.setBalance)

	apiGroup.GET("/transactions/account/:aid", s.getTransactionsByAccount)

	apiGroup.GET("/operations/entity/:eid", s.getOperationsByEntity)
	apiGroup.GET("/operation/:opid", s.getOperation)
	apiGroup.POST("/operation", s.addOperation)
	apiGroup.DELETE("/operation/:opid", s.deleteOperation)

	apiGroup.GET("/categories", s.getCategories)
	apiGroup.POST("/category", s.addCategory)

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
		Name:     "token",
		Secure:   true, // This should be set to true as soon as we implement HTTPS
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

	user, err := s.store.Login(login.User, login.Password)
	if err != nil {
		//FIXME look for error
		return err
	}

	// Valid login

	expiration := time.Now().Add(time.Hour * 72)

	// Set custom claims
	claims := &jwtCustomClaims{
		user,
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
		"user":    user,
		"expires": expiration.Unix(),
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

func extractClaims(c echo.Context) (*jwtCustomClaims, error) {
	authRaw := c.Get("token")
	if authRaw == nil {
		return nil, echo.ErrInternalServerError
	}
	auth, ok := authRaw.(*jwt.Token)
	if !ok {
		return nil, echo.ErrInternalServerError
	}
	claims, ok := auth.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, echo.ErrInternalServerError
	}
	return claims, nil
}

func getUser(c echo.Context) (User, error) {
	claims, err := extractClaims(c)
	if err != nil {
		return User{}, err
	}

	return claims.User, nil
}

func (s *Server) Greet(c echo.Context) error {

	u, err := getUser(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, u)
}
