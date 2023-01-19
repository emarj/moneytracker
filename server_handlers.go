package moneytracker

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// ************* Handlers *****************
func (s *Server) getTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"operation": s.store.GetOperationTypes(),
		"account":   s.store.GetAccountTypes(),
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

	cl, err := extractClaims(c)
	if err != nil {
		return err
	}

	el, err := s.store.GetEntitiesOfUser(cl.User.ID.Int64)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, el)
}

func (s *Server) getAllEntities(c echo.Context) error {

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

	cl, err := extractClaims(c)
	if err != nil {
		return err
	}

	_ = cl.User.ID.Int64 //TODO pass this to GetAccounts to check if it is admin
	el, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, el)
}

func (s *Server) getUserAccounts(c echo.Context) error {

	cl, err := extractClaims(c)
	if err != nil {
		return err
	}

	el, err := s.store.GetUserAccounts(cl.User.ID.Int64)
	if err != nil {
		return err
	}
	om := orderedmap.New[string, []Account]()
	/* entities := make(map[int64]string)
	shares := []EntityShare{}
	for _, a := range el {
		entities[a.OwnerID] = a.Owner.DisplayName
		for _, s := range a.Owner.Shares {
			if s.UserID == cl.User.ID.Int64 {
				shares = append(shares, s)
			}
		}
	}

	sort.Slice(shares, func(i int, j int) bool {
		return shares[i].Quota <= shares[j].Quota
	})


	for _, s := range shares {
		om.Set(entities[s.EntityID], []Account{})
	}
	*/
	for _, a := range el {
		list, _ := om.Get(a.Owner.DisplayName)
		list = append(list, a)
		om.Set(a.Owner.DisplayName, list)
	}

	return c.JSON(http.StatusOK, om)
}

func (s *Server) getAccountsByEntity(c echo.Context) error {
	eID, err := Atoi64(c.Param("eid"))
	if err != nil {
		return err
	}
	el, err := s.store.GetAccountsByEntity(eID)
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

func (s *Server) getUserBalance(c echo.Context) error {

	/* cl, err := extractClaims(c)
	if err != nil {
		return err
	}

	b, err := s.store.GetUserBalanceNow(cl.User.ID.Int64)
	if err != nil {
		return err
	} */
	return c.JSON(http.StatusOK, "not implemented")
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

func (s *Server) getOperationsOfUser(c echo.Context) error {
	cl, err := extractClaims(c)
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

	list, err := s.store.GetOperationsOfUser(cl.User.ID.Int64, limit)
	if err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, list, "\t")
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

	cl, err := extractClaims(c)
	if err != nil {
		return err
	}

	op.CreatedByID = cl.User.ID.Int64

	err = s.store.AddOperation(&op)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, op)
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

func (s *Server) getCategories(c echo.Context) error {
	cl, err := s.store.GetCategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cl)
}

func (s *Server) addCategory(c echo.Context) error {

	cat := Category{}

	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		return err
	}

	cat, err = s.store.AddCategory(cat.FullName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cat)
}

func (s *Server) addExpense(c echo.Context) error {

	e := Expense{}

	err := json.NewDecoder(c.Request().Body).Decode(&e)
	if err != nil {
		return err
	}

	op := e.ToOperation()

	err = s.store.AddOperation(&op)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, e)
}
