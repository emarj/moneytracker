package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/satori/go.uuid"

	"ronche.se/moneytracker/model"
)

func JSONHandler(srv model.Service) http.Handler {
	mux := http.NewServeMux()
	h := jsonHandler{srv}
	mux.HandleFunc("/", jsonResponseWriter(h.listExpenses))
	mux.HandleFunc("/add/", jsonResponseWriter(h.addExpense))
	mux.HandleFunc("/get/", jsonResponseWriter(h.getExpense))
	mux.HandleFunc("/update/", jsonResponseWriter(h.getExpense))
	mux.HandleFunc("/delete/", jsonResponseWriter(h.deleteExpense))
	return mux
}

type jsonHandler struct {
	srv model.Service
	//Google Sheets
}

func jsonResponseWriter(f func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, status, err := f(w, r)
		if err != nil {
			log.Printf("Error (%d) %v", status, err)
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}

	}
}

func (h *jsonHandler) listExpenses(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	es, err := h.srv.ExpensesGetN(20)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	cat, err := h.srv.CategoriesGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	u, err := h.srv.UsersGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	pm, err := h.srv.PaymentMethodsGetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	type result struct {
		Expenses       []*model.Expense
		Categories     []*model.Category
		Users          []*model.User
		PaymentMethods []*model.PaymentMethod
	}

	return result{es, cat, u, pm}, http.StatusOK, nil
}

func (h *jsonHandler) getExpense(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	idstr := r.URL.Query().Get("uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	e, err := h.srv.ExpenseGet(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	return e, http.StatusOK, nil
}

func (h *jsonHandler) addExpense(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method)
	}

	var e model.Expense
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("cannot parse json request: %v", err)
	}

	result, err := h.srv.ExpenseInsert(&e)
	if err != nil {
		return nil, http.StatusInternalServerError, nil
	}

	return result, http.StatusOK, nil

}

func (h *jsonHandler) updateExpense(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method)
	}

	var e model.Expense
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("cannot parse json request: %v", err)
	}

	result, err := h.srv.ExpenseUpdate(&e)
	if err != nil {
		return nil, http.StatusInternalServerError, nil
	}

	return result, http.StatusOK, nil

}

func (h *jsonHandler) deleteExpense(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	idstr := r.URL.Query().Get("uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = h.srv.ExpenseDelete(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	return "ok", http.StatusOK, nil
}
