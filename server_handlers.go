package moneytracker

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ************* Handlers *****************
func (s *Server) getTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Operation []OperationType `json:"operation"`
		Account   []AccountType   `json:"account"`
	}{
		s.store.GetOperationTypes(),
		s.store.GetAccountTypes(),
	})
}

// Entities

func (s *Server) getEntity(c echo.Context) error {

	eID, err := Atoi64(c.Param("eid"))
	if err != nil {
		return err
	}

	e, err := s.store.GetEntity(eID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.ErrNotFound
		}
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

// Accounts

func (s *Server) getAccount(c echo.Context) error {

	aID, err := Atoi64(c.Param("aid"))
	if err != nil {
		return err
	}

	a, err := s.store.GetAccount(aID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.ErrNotFound
		}
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
	eID, err := Atoi64(c.Param("eid"))
	if err != nil {
		return err
	}
	el, err := s.store.GetAccountsByEntity(int64(eID))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, el)
}

func (s *Server) getBalanceHistory(c echo.Context) error {
	aID, err := Atoi64(c.Param("aid"))
	if err != nil {
		return err
	}
	bl, err := s.store.GetBalanceHistory(aID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bl)
}

func (s *Server) getBalance(c echo.Context) error {

	aID, err := Atoi64(c.Param("aid"))
	if err != nil {
		return err
	}

	b, err := s.store.GetBalanceNow(aID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	return c.JSON(http.StatusOK, b)
}

func (s *Server) setBalance(c echo.Context) error {
	b := Balance{}

	err := json.NewDecoder(c.Request().Body).Decode(&b)
	if err != nil {
		return err
	}

	/* claims, err := extractClaims(c)
	if err != nil {
		return err
	} */

	err = s.store.SetBalance(&b)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (s *Server) addAccount(c echo.Context) error {

	a := Account{}

	err := json.NewDecoder(c.Request().Body).Decode(&a)
	if err != nil {
		return err
	}

	/* claims, err := extractClaims(c)
	if err != nil {
		return err
	} */

	err = s.store.AddAccount(&a)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}

func (s *Server) deleteAccount(c echo.Context) error {
	aID, err := Atoi64(c.Param("aid"))
	if err != nil {
		return err
	}

	err = s.store.DeleteAccount(aID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

// Transactions and Operations

/*func (s *Server) getTransactions(c echo.Context) error {
	tl, err := s.store.GetTransactions()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}*/

func (s *Server) getTransactionsByAccount(c echo.Context) error {
	aID, err := Atoi64(c.Param("aid"))
	if err != nil {
		return err
	}

	limit := int64(5)
	limitStr := c.QueryParam("limit")
	if limitStr != "" {
		limit, err = Atoi64(limitStr)
		if err != nil {
			return err
		}
	}

	tl, err := s.store.GetTransactionsByAccount(aID, limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tl)
}

func (s *Server) getOperationsByEntity(c echo.Context) error {
	aID, err := Atoi64(c.Param("eid"))
	if err != nil {
		return err
	}

	limit := int64(5)
	limitStr := c.QueryParam("limit")
	if limitStr != "" {
		limit, err = Atoi64(limitStr)
		if err != nil {
			return err
		}
	}

	list, err := s.store.GetOperationsByEntity(aID, limit)
	if err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, list, "\t")
}

func (s *Server) getOperation(c echo.Context) error {

	opID, err := Atoi64(c.Param("opid"))
	if err != nil {
		return err
	}

	op, err := s.store.GetOperation(opID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.ErrNotFound
		}
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

	claims, err := extractClaims(c)
	if err != nil {
		return err
	}

	op.CreatedByID = claims.User.ID.Int64

	err = s.store.AddOperation(&op)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, op)
}

func (s *Server) getCategories(c echo.Context) error {
	cl, err := s.store.GetCategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cl)
}

func (s *Server) deleteOperation(c echo.Context) error {

	opID, err := Atoi64(c.Param("opid"))
	if err != nil {
		return err
	}

	err = s.store.DeleteOperation(opID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
