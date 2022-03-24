package moneytracker

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"ronche.se/moneytracker/db"
	"ronche.se/moneytracker/domain"
)

type Handler struct {
	store db.Store
}

func NewHandler(store db.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetUsers(c echo.Context) error {
	ul, err := h.store.GetUsers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ul)
}

func (h *Handler) GetAccount(c echo.Context) error {
	a, err := h.store.GetAccount(c.Param("aid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, a)
}

func (h *Handler) GetAccountsOfUser(c echo.Context) error {
	al, err := h.store.GetAccountsOfUser(c.Param("uid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, al)
}

func (h *Handler) GetTransaction(c echo.Context) error {

	tID, err := uuid.FromString(c.Param("tid"))
	if err != nil {
		return err
	}

	t, err := h.store.GetTransaction(tID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *Handler) GetTransactions(c echo.Context) error {
	tl, err := h.store.GetTransactionsByAccount(c.Param("aid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}

func (h *Handler) GetTransactionsOfUser(c echo.Context) error {
	tl, err := h.store.GetTransactionsByUser(c.Param("uid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}

func (h *Handler) InsertTransaction(c echo.Context) error {
	t := domain.Transaction{}

	err := json.NewDecoder(c.Request().Body).Decode(&t)
	if err != nil {
		return err
	}

	_, err = h.store.InsertTransaction(&t)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}
