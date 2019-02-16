package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	//"github.com/julienschmidt/httprouter"
	"github.com/bouk/httprouter"

	"github.com/jpillora/cookieauth"

	"html/template"

	uuid "github.com/satori/go.uuid"

	"ronche.se/moneytracker/model"
)

func HTMLHandler(dbSrv model.Service, tmplPath string, prefix string) (http.Handler, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"IsNeg": func(d decimal.Decimal) bool {
			return d.LessThan(decimal.Zero)
		},
		"FixFormatDec2": func(d decimal.Decimal) string {
			return d.StringFixed(2)
		},
		"AutoFormatDec": func(d decimal.Decimal) string {
			return d.String()
		},
		"Now": func(format string) string {
			loc, _ := time.LoadLocation("Europe/Rome")
			return time.Now().In(loc).Format(format)
		},
		"ToLower": func(str string) string {
			return strings.ToLower(str)
		},
		"Prefix": func() string {
			return prefix
		},
		"SubDec": func(a, b decimal.Decimal) decimal.Decimal {
			return a.Sub(b)
		},
		"PercDec": func(a, b decimal.Decimal) decimal.Decimal {
			return a.Div(b).Mul(decimal.New(1, 2)).Ceil()
		},
	}).ParseGlob(path.Join(tmplPath, "*"))

	if err != nil {
		return nil, err
	}

	h := htmlHandler{dbSrv, t, prefix}

	router := httprouter.New()
	router.GET(prefix+"/", h.render(h.home))
	router.POST(prefix+"/", h.render(h.home))
	router.GET(prefix+"/all/", h.render(h.all))
	router.POST(prefix+"/add/", h.render(h.add))
	router.GET(prefix+"/view/:uuid", h.render(h.view))
	router.POST(prefix+"/update/", h.render(h.update))
	router.GET(prefix+"/delete/:uuid", h.render(h.delete))
	/*router.GET(prefix+"/sheet/add/:uuid", h.render(h.addToSheet))
	router.POST(prefix+"/sheet/add/:uuid", h.render(h.addToSheet))
	router.GET(prefix+"/sheet/reset", h.render(h.resetSheet))*/

	router.ServeFiles(prefix+"/static/*filepath", http.Dir("./static"))

	protected := cookieauth.Wrap(router, "spendi", "schei")

	return protected, nil
}

type htmlHandler struct {
	dbSrv  model.Service
	tmpl   *template.Template
	prefix string
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
			http.Redirect(w, r, h.prefix+res.RedirectTo, res.Status)
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

func (h *htmlHandler) home(r *http.Request) *htmlResponse {
	ts, err := h.dbSrv.TransactionsGetNOrderByInserted(5)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	us, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	tps, err := h.dbSrv.TypesGetAll()
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
		Transactions   []*model.Transaction
		Transaction    *model.Transaction
		Types          []*model.Type
		Categories     []*model.Category
		Users          []*model.User
		PaymentMethods []*model.PaymentMethod
	}

	return resOK(result{ts, nil, tps, cat, us, pm}, "home")
}

func (h *htmlHandler) all(r *http.Request) *htmlResponse {
	ts, err := h.dbSrv.TransactionsGetNOrderByDate(99999) //NEED TO IMPLEMENT NO LIMIT
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	us, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	tps, err := h.dbSrv.TypesGetAll()
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
		Transactions   []*model.Transaction
		Categories     []*model.Category
		Users          []*model.User
		Types          []*model.Type
		PaymentMethods []*model.PaymentMethod
	}

	return resOK(result{ts, cat, us, tps, pm}, "all")
}

func (h *htmlHandler) view(r *http.Request) *htmlResponse {

	idstr := httprouter.GetParam(r, "uuid")

	id, err := uuid.FromString(idstr)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}
	//Get transaction
	t, err := h.dbSrv.TransactionGet(id)
	if err != nil {
		return resError(err, http.StatusNotFound)
	}

	//Get resources for UI
	us, err := h.dbSrv.UsersGetAll()
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}
	tps, err := h.dbSrv.TypesGetAll()
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
		Transactions   []*model.Transaction
		Transaction    *model.Transaction
		Types          []*model.Type
		Categories     []*model.Category
		Users          []*model.User
		PaymentMethods []*model.PaymentMethod
	}

	return resOK(result{nil, t, tps, cat, us, pm}, "view")
}

func (h *htmlHandler) parseForm(r *http.Request) (*model.Transaction, error) {
	r.ParseForm()

	t := model.Transaction{Description: r.FormValue("Description")}

	d, err := time.Parse("2006-01-02", r.FormValue("Date"))
	if err != nil {
		return &t, err
	}
	t.Date.Time = d

	am, err := decimal.NewFromString(r.FormValue("Amount"))
	if err != nil {
		return &t, err
	}
	if am.Equals(decimal.Zero) {
		return &t, fmt.Errorf("amount cannot be zero")
	}
	t.Amount = am

	userID, err := strconv.Atoi(r.FormValue("UserID"))
	if err != nil {
		return &t, err
	}
	t.User.ID = userID

	typeID, err := strconv.Atoi(r.FormValue("TypeID"))
	if err != nil {
		return &t, err
	}
	t.Type.ID = typeID

	catID, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		return &t, err
	}

	t.Category.ID = catID

	pmid, err := strconv.Atoi(r.FormValue("MethodID"))
	if err != nil {
		return &t, err
	}
	t.PaymentMethod.ID = pmid

	if r.FormValue("Shared") == "on" {
		t.Shared = true

		sq, err := decimal.NewFromString(r.FormValue("SharedQuota"))
		if err != nil {
			return &t, err
		}
		if sq.Equals(decimal.Zero) {
			return &t, fmt.Errorf("Shared Quota cannot be zero")
		}
		shareWithID := 1
		if t.User.ID == 1 {
			shareWithID = 2
		}
		t.Shares = append(t.Shares, &model.Share{t.UUID, shareWithID, "", sq})
	}

	t.GeoLocation = r.FormValue("GeoLoc")

	return &t, nil
}

func (h *htmlHandler) add(r *http.Request) *htmlResponse {

	t, err := h.parseForm(r)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	err = h.dbSrv.TransactionInsert(t)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	return resRedirect("/", http.StatusTemporaryRedirect)

}

func (h *htmlHandler) update(r *http.Request) *htmlResponse {

	t, err := h.parseForm(r)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	id, err := uuid.FromString(r.FormValue("UUID"))
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	t.UUID = id

	//More checks

	err = h.dbSrv.TransactionUpdate(t)
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
	err = h.dbSrv.TransactionDelete(id)
	if err != nil {
		return resError(err, http.StatusNotFound)
	}
	return resRedirect("/", http.StatusTemporaryRedirect)
}

/*
func (h *htmlHandler) addToSheet(r *http.Request) *htmlResponse {
	id, err := uuid.FromString(httprouter.GetParam(r, "uuid"))
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	t, err := h.dbSrv.TransactionGet(id)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	err = h.sheetSrv.Insert(*t)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	t.InSheet = true
	err = h.dbSrv.TransactionUpdate(t)
	if err != nil {
		return resError(err, http.StatusInternalServerError)
	}

	return resRedirect("/", http.StatusTemporaryRedirect)
}

func (h *htmlHandler) resetSheet(r *http.Request) *htmlResponse {
	ts, err := h.dbSrv.TransactionsGetNOrderByDate(100)
	if err != nil {
		return resError(err, http.StatusBadRequest)
	}

	for _, t := range ts {
		err = h.sheetSrv.Insert(*t)
		if err != nil {
			return resError(err, http.StatusInternalServerError)
		}

		t.InSheet = true
		err = h.dbSrv.TransactionUpdate(t)
		if err != nil {
			return resError(err, http.StatusInternalServerError)
		}
	}

	return resRedirect("/", http.StatusTemporaryRedirect)
}
*/
