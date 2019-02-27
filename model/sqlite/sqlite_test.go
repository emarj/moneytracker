package sqlite

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"ronche.se/moneytracker/model"
)

func TestNew(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

}

func TestTransactionCRUD(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

	u1, err := s.UserInsert("Marco")
	if err != nil {
		t.Error(err)
		return
	}
	u2, err := s.UserInsert("Arianna")
	if err != nil {
		t.Error(err)
		return
	}
	tp1, err := s.TypeInsert("Expense")
	if err != nil {
		t.Error(err)
		return
	}
	tp2, err := s.TypeInsert("Transfer")
	if err != nil {
		t.Error(err)
		return
	}
	cat, err := s.CategoryInsert("Uncategorized")
	if err != nil {
		t.Error(err)
		return
	}
	pm, err := s.PaymentMethodInsert("CC")
	if err != nil {
		t.Error(err)
		return
	}

	tr1 := model.Transaction{Date: model.Date(time.Now().Local()), Type: *tp1, User: *u1, Method: *pm, Category: *cat}
	tr2 := model.Transaction{Date: model.Date(time.Now().Local()), Type: *tp2, User: *u2, Method: *pm, Category: *cat}
	err = s.TransactionInsert(&tr1)
	if err != nil {
		t.Error(err)
		return
	}
	err = s.TransactionInsert(&tr2)
	if err != nil {
		t.Error(err)
		return
	}

	tr1.Amount, _ = decimal.NewFromString("123456")
	err = s.TransactionUpdate(&tr1)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := s.TransactionGet(tr1.UUID)
	if err != nil {
		t.Error(err)
		return
	}

	if !res.Amount.Equal(tr1.Amount) {
		t.Errorf("%s != %s", res.Amount.String(), tr1.Amount.String())
		return
	}

	trs, err := s.TransactionsGetNOrderByDate(100)
	if err != nil {
		t.Error(err)
		return
	}

	for _, tr := range trs {
		t.Log(tr.User)
	}

	if len(trs) != 2 {
		t.Errorf("2 != %d", len(trs))
		return
	}

	err = s.TransactionDelete(tr1.UUID)
	if err != nil {
		t.Error(err)
		return
	}

	trs, err = s.TransactionsGetNOrderByDate(100)
	if err != nil {
		t.Error(err)
		return
	}

	for _, tr := range trs {
		t.Log(tr.User)
	}

	if len(trs) != 1 {
		t.Errorf("1 != %d", len(trs))
		return
	}

}
