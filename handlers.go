package moneytracker

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"ronche.se/moneytracker/db"
	"ronche.se/moneytracker/domain"
)

func NewHandler(store db.Store) *Server {
	return &Server{store: store}
}

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
	t := domain.Transaction{}

	err := json.NewDecoder(c.Request().Body).Decode(&t)
	if err != nil {
		return err
	}

	err = s.store.AddTransaction(&t)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}
