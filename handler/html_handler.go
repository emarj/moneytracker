package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	//"github.com/julienschmidt/httprouter"
	"github.com/bouk/httprouter"

	"github.com/jpillora/cookieauth"

	"html/template"

	uuid "github.com/satori/go.uuid"

	"ronche.se/moneytracker/model"
	"ronche.se/moneytracker/sheet"
)

func HTMLHandler(dbSrv model.Service, sheetSrv *sheet.SheetService, tmplPath string) (http.Handler, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"IsNeg": func(d decimal.Decimal) bool {
			return d.LessThan(decimal.Zero)
		},
		"FormatDecimal": func(d decimal.Decimal) string {
			return d.StringFixed(2)
		},
		"Now": func(format string) string {
			loc, _ := time.LoadLocation("Europe/Rome")
			return time.Now().In(loc).Format(format)
		},
	}).ParseGlob(path.Join(tmplPath, "*"))

	if err != nil {
		return nil, err
	}

	h := htmlHandler{dbSrv, sheetSrv, t}

	router := httprouter.New()
	router.GET("/", h.render(h.index))
	router.POST("/", h.render(h.index))
	router.GET("/all/", h.render(h.all))
	router.POST("/add/", h.render(h.add))
	router.GET("/view/:uuid", h.render(h.view))
	router.POST("/update/", h.render(h.update))
	router.GET("/delete/:uuid", h.render(h.delete))
	router.GET("/sheet/add/:uuid", h.render(h.addToSheet))
	router.POST("/sheet/add/:uuid", h.render(h.addToSheet))
	router.GET("/sheet/reset", h.render(h.resetSheet))

	protected := cookieauth.Wrap(router, "spendi", "schei")

	return protected, nil
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

func resRedirect(url string, status int) *htmlResponse {
	return &htmlResponse{IsRedirect: true, RedirectTo: url, Status: status}
}

func resOK(data interface{}, tmpl string) *htmlResponse {
	return &htmlResponse{Data: data, TmplName: tmpl, Status: http.StatusOK}
}

func resError(err error, status int) *htmlResponse {
	return &htmlResponse{Err: err, Status: status}
}

func (h *htmlHandler) render(f func(r *http.Request) *htmlResponse) http.HandlerFunc {
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

func (h *htmlHandler) index(r *http.Request) *htmlResponse {
	es, err := h.dbSrv.ExpensesGetNOrderByInserted(10)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	cat, err := h.dbSrv.CategoriesGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	pm, err := h.dbSrv.PaymentMethodsGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	type result struct {
		Expenses       []*model.Expense
		Expense        *model.Expense
		Categories     []*model.Category
		Users          []string
		PaymentMethods []*model.PaymentMethod
	}

	return resOK(result{es, nil, cat, model.Users, pm}, "index")
}

func (h *htmlHandler) all(r *http.Request) *htmlResponse {
	es, err := h.dbSrv.ExpensesGetNOrderByDate(1000)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	cat, err := h.dbSrv.CategoriesGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	pm, err := h.dbSrv.PaymentMethodsGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	type result struct {
		Expenses       []*model.Expense
		Categories     []*model.Category
		Users          []string
		PaymentMethods []*model.PaymentMethod
	}

	return resOK(result{es, cat, model.Users, pm}, "all")
}

func (h *htmlHandler) view(r *http.Request) *htmlResponse {

	idstr := httprouter.GetParam(r, "uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	e, err := h.dbSrv.ExpenseGet(id)
	if err != nil {
		return resError(err, http.StatusNotFound)
	}

	cat, err := h.dbSrv.CategoriesGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	pm, err := h.dbSrv.PaymentMethodsGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	type result struct {
		Expenses       []*model.Expense
		Expense        *model.Expense
		Categories     []*model.Category
		Users          []string
		PaymentMethods []*model.PaymentMethod
	}

	return resOK(result{nil, e, cat, model.Users, pm}, "view")
}

func (h *htmlHandler) parseForm(r *http.Request) (*model.Expense, error) {
	r.ParseForm()

	e := model.Expense{Description: r.FormValue("Description")}

	date, err := time.Parse("2006-01-02", r.FormValue("Date"))
	if err != nil {
		return &e, err
	}
	e.Date = date

	am, err := decimal.NewFromString(r.FormValue("Amount"))
	if err != nil {
		return &e, err
	}
	if am.Equals(decimal.Zero) {
		return &e, fmt.Errorf("amount cannot be zero")
	}
	e.Amount = am

	e.Who = r.FormValue("Who")

	typ, err := strconv.Atoi(r.FormValue("Type"))
	if err != nil {
		return &e, err
	}
	e.Type = typ

	catid, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		return &e, err
	}

	e.Category = &model.Category{ID: catid}

	pmid, err := strconv.Atoi(r.FormValue("MethodID"))
	if err != nil {
		return &e, err
	}
	e.Method = &model.PaymentMethod{ID: pmid}

	if r.FormValue("Shared") == "on" {
		e.Shared = true

		quota, err := strconv.Atoi(r.FormValue("ShareQuota"))
		if err != nil {
			return &e, err
		}
		if quota == 0 {
			return &e, fmt.Errorf("quota cannot be zero")
		}
		e.ShareQuota = quota
	}

	if r.FormValue("InSheet") == "on" {
		e.InSheet = true
	}

	return &e, nil
}

func (h *htmlHandler) add(r *http.Request) *htmlResponse {

	e, err := h.parseForm(r)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	err = h.dbSrv.ExpenseInsert(e)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	if !e.InSheet {
		return resRedirect("/sheet/add/"+e.UUID.String(), http.StatusTemporaryRedirect)
	}

	return resRedirect("/", http.StatusTemporaryRedirect)

}

func (h *htmlHandler) update(r *http.Request) *htmlResponse {

	e, err := h.parseForm(r)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	id, err := uuid.FromString(r.FormValue("UUID"))
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	e.UUID = id

	//More checks

	err = h.dbSrv.ExpenseUpdate(e)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	return resRedirect("/", http.StatusTemporaryRedirect)

}

func (h *htmlHandler) delete(r *http.Request) *htmlResponse {

	id, err := uuid.FromString(httprouter.GetParam(r, "uuid"))
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}
	err = h.dbSrv.ExpenseDelete(id)
	if err != nil {
		return resError(err, http.StatusNotFound)
	}
	return resRedirect("/", http.StatusTemporaryRedirect)
}

func (h *htmlHandler) addToSheet(r *http.Request) *htmlResponse {
	id, err := uuid.FromString(httprouter.GetParam(r, "uuid"))
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	e, err := h.dbSrv.ExpenseGet(id)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	err = h.sheetSrv.Insert(*e)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	e.InSheet = true
	err = h.dbSrv.ExpenseUpdate(e)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	return resRedirect("/", http.StatusTemporaryRedirect)
}

func (h *htmlHandler) resetSheet(r *http.Request) *htmlResponse {
	es, err := h.dbSrv.ExpensesGetNOrderByDate(100)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	for _, e := range es {
		err = h.sheetSrv.Insert(*e)
		if err != nil {
			return resError(err, http.StatusInternalServerError)
		}

		e.InSheet = true
		err = h.dbSrv.ExpenseUpdate(e)
		if err != nil {
			return resError(err, http.StatusInternalServerError)
		}
	}

	return resRedirect("/", http.StatusTemporaryRedirect)
}
