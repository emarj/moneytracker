package handler

import (
	"bytes"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"html/template"

	uuid "github.com/satori/go.uuid"

	"ronche.se/moneytracker/model"
	"ronche.se/moneytracker/sheet"
	"ronche.se/moneytracker/utils"
)

func HTMLHandler(dbSrv model.Service, sheetSrv *sheet.SheetService, tmplPath string) (http.Handler, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"formatDecimal": func(n int) string {
			return utils.FormatDecimal(n)
		},
	}).ParseGlob(path.Join(tmplPath, "*"))

	if err != nil {
		return nil, err
	}

	h := htmlHandler{dbSrv, sheetSrv, t}

	router := httprouter.New()
	router.GET("/", h.htmlResponseWriter(h.listExpenses))
	router.POST("/", h.htmlResponseWriter(h.listExpenses))
	router.POST("/add/", h.htmlResponseWriter(h.addExpense))
	router.GET("/delete/:uuid", h.htmlResponseWriter(h.deleteExpense))
	router.GET("/sheet/add/:uuid", h.htmlResponseWriter(h.addExpenseToSheet))

	return router, nil
}

type htmlHandler struct {
	dbSrv    model.Service
	sheetSrv *sheet.SheetService
	tmpl     *template.Template
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
func (h *htmlHandler) htmlResponseWriter(f func(r *http.Request, ps httprouter.Params) *htmlResponse) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		buf := new(bytes.Buffer)
		res := f(r, ps)
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

func (h *htmlHandler) listExpenses(r *http.Request, ps httprouter.Params) *htmlResponse {
	es, err := h.dbSrv.ExpensesGetN(20)
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	cat, err := h.dbSrv.CategoriesGetAll()
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	u, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}
	pm, err := h.dbSrv.PaymentMethodsGetAll()
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

func (h *htmlHandler) getExpense(r *http.Request, ps httprouter.Params) *htmlResponse {

	idstr := r.URL.Query().Get("uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}

	e, err := h.dbSrv.ExpenseGet(id)
	if err != nil {
		return htmlResErr(err, http.StatusNotFound)
	}
	return htmlResOK(e, "view")
}

func (h *htmlHandler) addExpense(r *http.Request, ps httprouter.Params) *htmlResponse {

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

	_, err = h.dbSrv.ExpenseInsert(&e)
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

	result, err := h.dbSrv.ExpenseUpdate(&e)
	if err != nil {
		return nil, http.StatusInternalServerError, nil
	}

	return result, http.StatusOK, nil

}*/

func (h *htmlHandler) deleteExpense(r *http.Request, ps httprouter.Params) *htmlResponse {

	id, err := uuid.FromString(ps.ByName("uuid"))
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}
	err = h.dbSrv.ExpenseDelete(id)
	if err != nil {
		return htmlResErr(err, http.StatusNotFound)
	}
	return htmlResRedir("/", http.StatusTemporaryRedirect)
}

func (h *htmlHandler) addExpenseToSheet(r *http.Request, ps httprouter.Params) *htmlResponse {
	id, err := uuid.FromString(ps.ByName("uuid"))
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}

	e, err := h.dbSrv.ExpenseGet(id)
	if err != nil {
		return htmlResErr(err, http.StatusBadRequest)
	}

	err = h.sheetSrv.Insert(*e)
	if err != nil {
		return htmlResErr(err, http.StatusInternalServerError)
	}

	/*
		e.InSheet = true
		err := h.dbSrv.ExpenseUpdate(e)
		if err != nil {
			return htmlResErr(err, http.StatusInternalServerError)
		}
	*/

	return htmlResRedir("/", http.StatusTemporaryRedirect)
}
