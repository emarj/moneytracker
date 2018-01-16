package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"html/template"

	"github.com/satori/go.uuid"

	"ronche.se/expensetracker/model"
	"ronche.se/expensetracker/utils"
)

func HTMLHandler(srv model.Service, tmplPath string) (http.Handler, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"formatDecimal": func(n int) string {
			return utils.FormatDecimal(n)
		},
	}).ParseGlob(path.Join(tmplPath, "*"))

	if err != nil {
		return nil, err
	}

	h := htmlHandler{srv, t}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.htmlResponseWriter(h.listExpenses))
	mux.HandleFunc("/add/", h.htmlResponseWriter(h.addExpense))
	mux.HandleFunc("/delete/", h.htmlResponseWriter(h.deleteExpense))

	return mux, nil
}

type htmlHandler struct {
	srv  model.Service
	tmpl *template.Template
	//Google Sheets
}

type htmlResponse struct {
	Data       interface{}
	TmplName   string
	Status     int
	Err        error
	IsRedirect bool
	RedirectTo string
}

func htmlResRedir(url string, status int) *htmlResponse {
	return &htmlResponse{IsRedirect: true, RedirectTo: url, Status: status}
}

func htmlResOK(data interface{}, tmpl string) *htmlResponse {
	return &htmlResponse{Data: data, TmplName: tmpl, Status: http.StatusOK}
}

func htmlResErr(err error, status int) *htmlResponse {
	return &htmlResponse{Err: err, Status: status}
}
func (h *htmlHandler) htmlResponseWriter(f func(r *http.Request) *htmlResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		buf := new(bytes.Buffer)
		res := f(r)
		if res.IsRedirect {
			http.Redirect(w, r, res.RedirectTo, res.Status)
			return
		}

		if res.Err != nil {
			log.Printf("Error (%d) %v", res.Status, res.Err)
			res.Data = res.Err
			res.TmplName = "error"
		}
		err := h.tmpl.ExecuteTemplate(buf, res.TmplName, res.Data)
		if err != nil {
			log.Printf("could not execute template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(res.Status)

		_, err = buf.WriteTo(w)
		if err != nil {
			log.Printf("could not write response: %v", err)
		}

	}
}

func (h *htmlHandler) listExpenses(r *http.Request) *htmlResponse {
	es, err := h.srv.ExpensesGetN(20)
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	cat, err := h.srv.CategoriesGetAll()
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	u, err := h.srv.UsersGetAll()
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	pm, err := h.srv.PaymentMethodsGetAll()
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	type result struct {
		Expenses       []*model.Expense
		Categories     []*model.Category
		Users          []*model.User
		PaymentMethods []*model.PaymentMethod
	}

	return htmlResOK(result{es, cat, u, pm}, "index")
}

func (h *htmlHandler) getExpense(r *http.Request) *htmlResponse {

	idstr := r.URL.Query().Get("uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}

	e, err := h.srv.ExpenseGet(id)
	if err != nil {
		return htmlResErr(err, http.StatusNotFound)
	}
	return htmlResOK(e, "view")
}

func (h *htmlHandler) addExpense(r *http.Request) *htmlResponse {

	if r.Method != http.MethodPost {
		return htmlResErr(fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}

	r.ParseForm()

	e := model.Expense{Description: r.FormValue("Description")}

	if r.FormValue("Date") != "" {
		date, err := time.Parse("2006-01-02", r.FormValue("Date"))
		if err != nil {
			return htmlResErr(err, http.StatusBadRequest)
		}
		e.Date = date
	} else {
		e.Date = time.Now().Local()
	}

	am, err := utils.ParseDecimal(r.FormValue("Amount"))
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}
	e.Amount = am

	uid, err := strconv.Atoi(r.FormValue("WhoID"))
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}
	e.Who = &model.User{ID: uid}

	catid, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}
	e.Category = &model.Category{ID: catid}

	pmid, err := strconv.Atoi(r.FormValue("MethodID"))
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}
	e.Method = &model.PaymentMethod{ID: pmid}

	if r.FormValue("Shared") == "on" {
		e.Shared = true

		quota, err := strconv.Atoi(r.FormValue("ShareQuota"))
		if err != nil {
			return htmlResErr(err, http.StatusBadRequest)
		}
		e.ShareQuota = quota
	}

	_, err = h.srv.ExpenseInsert(&e)
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}

	return htmlResRedir("/", http.StatusTemporaryRedirect)

}

/*func (h *htmlHandler) updateExpense(r *http.Request) *htmlResponse {

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

}*/

func (h *htmlHandler) deleteExpense(r *http.Request) *htmlResponse {
	idstr := r.URL.Query().Get("uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}
	err = h.srv.ExpenseDelete(id)
	if err != nil {
		return htmlResErr(err, http.StatusNotFound)
	}
	return htmlResRedir("/", http.StatusTemporaryRedirect)
}
