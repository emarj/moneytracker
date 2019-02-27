package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/jpillora/cookieauth"
	"github.com/shopspring/decimal"

	//"github.com/julienschmidt/httprouter"
	"github.com/bouk/httprouter"

	"ronche.se/moneytracker/model"
)

func NewJSONHandler(dbSrv model.Service, prefix string) (http.Handler, error) {

	h := JSONHandler{dbSrv, prefix}

	router := httprouter.New()

	router.GET(prefix+"/home/", h.render(h.home))

	router.GET(prefix+"/transactions/:limit/:offset", h.render(h.getTxs))

	router.GET(prefix+"/transaction/:uuid", h.render(h.getTx))
	router.POST(prefix+"/transaction/", h.render(h.addTx))
	router.PUT(prefix+"/transaction/", h.render(h.updateTx))
	router.DELETE(prefix+"/transaction/:uuid", h.render(h.deleteTx))

	router.GET(prefix+"/users/", h.render(h.getUsers))
	router.GET(prefix+"/categories/", h.render(h.getCategories))
	router.GET(prefix+"/types/", h.render(h.getTypes))

	router.ServeFiles(prefix+"/static/*filepath", http.Dir("./static"))

	protected := cookieauth.Wrap(router, "spendi", "schei")

	return protected, nil
}

type JSONHandler struct {
	dbSrv  model.Service
	prefix string
}

func (h *JSONHandler) render(f func(r *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		data, status, err := f(r)

		if err != nil {
			log.Printf("Error (%d) %v", status, err)
			data = err.Error()
		}
		w.WriteHeader(status)

		if json.NewEncoder(w).Encode(data) != nil {
			log.Printf("Error in JSON encoding")
			w.Write([]byte("Error in JSON encoding")) //should use JSON error
		}
	}
}

func (h *JSONHandler) home(r *http.Request) (interface{}, int, error) {
	ts, err := h.dbSrv.TransactionsGetNOrderByModified(5, 0)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	us, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tp, err := h.dbSrv.TypesGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	cats, err := h.dbSrv.CategoriesGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	pms, err := h.dbSrv.PaymentMethodsGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	type result struct {
		Transactions []*model.Transaction `json:"transactions"`
		Types        []*model.Type        `json:"types"`
		Categories   []*model.Category    `json:"categories"`
		Users        []*model.User        `json:"users"`
		Methods      []*model.Method      `json:"methods"`
	}

	return result{ts, tp, cats, us, pms}, http.StatusOK, nil

}

func (h *JSONHandler) getTxs(r *http.Request) (interface{}, int, error) {
	limit, err := strconv.Atoi(httprouter.GetParam(r, "limit"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	offset, err := strconv.Atoi(httprouter.GetParam(r, "offset"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	ts, err := h.dbSrv.TransactionsGetNOrderByDate(limit, offset)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return ts, http.StatusOK, nil
}

func (h *JSONHandler) userView(r *http.Request) (interface{}, int, error) {
	id, err := strconv.Atoi(httprouter.GetParam(r, "userid"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	ts, err := h.dbSrv.TransactionsGetNByUser(id, 99999)

	balance, err := h.dbSrv.TransactionsGetBalance(id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	type Credit struct {
		WithName string
		Amount   decimal.Decimal
	}

	us, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	credits := make([]Credit, 0, len(us)-1)
	for _, u := range us {
		if u.ID != id {
			c, err := h.dbSrv.TransactionsGetCredit(id, u.ID)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			if !c.IsZero() {
				credits = append(credits, Credit{u.Name, c})
			}

		}
	}

	type result struct {
		Transactions []*model.Transaction
		UserID       int
		Balance      decimal.Decimal
		Credits      []Credit
	}

	return result{ts, id, balance, credits}, http.StatusOK, nil
}

func (h *JSONHandler) getTx(r *http.Request) (interface{}, int, error) {

	idstr := httprouter.GetParam(r, "uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	//Get transaction
	t, err := h.dbSrv.TransactionGet(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return t, http.StatusOK, nil
}

func (h *JSONHandler) addTx(r *http.Request) (interface{}, int, error) {

	t := &model.Transaction{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(t)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	err = h.dbSrv.TransactionInsert(t)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return t, http.StatusOK, nil

}

func (h *JSONHandler) updateTx(r *http.Request) (interface{}, int, error) {

	t := &model.Transaction{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(t)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	err = h.dbSrv.TransactionUpdate(t)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return t, http.StatusOK, nil

}

func (h *JSONHandler) deleteTx(r *http.Request) (interface{}, int, error) {

	id, err := uuid.FromString(httprouter.GetParam(r, "uuid"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = h.dbSrv.TransactionDelete(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	return nil, http.StatusOK, nil
}

func (h *JSONHandler) getUsers(r *http.Request) (interface{}, int, error) {

	us, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return us, http.StatusOK, nil

}

func (h *JSONHandler) getTypes(r *http.Request) (interface{}, int, error) {

	tp, err := h.dbSrv.TypesGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return tp, http.StatusOK, nil

}

func (h *JSONHandler) getCategories(r *http.Request) (interface{}, int, error) {

	cats, err := h.dbSrv.CategoriesGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return cats, http.StatusOK, nil

}
