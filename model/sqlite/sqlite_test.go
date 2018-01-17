package sqlite

import (
	"log"
	"os"
	"testing"
	"time"

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

func TestCategoryInsert(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Fail()
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

	cat1, err := s.CategoryInsert("Prova")
	if err != nil {
		t.Error(err)
	}

	if cat1.ID != 1 {
		t.Errorf("Expecting ID %d got %d", 1, cat1.ID)
	}

	cat2, err := s.CategoryInsert("Prova2")
	if err != nil {
		t.Error(err)
	}
	if cat2.ID != 2 {
		t.Errorf("Expecting ID %d got %d", 2, cat2.ID)
	}

}

func TestCategoriesGetAll(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Fail()
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

	s.CategoryInsert("foo")
	s.CategoryInsert("bar")

	cats, err := s.CategoriesGetAll()
	if err != nil {
		t.Error(err)
	}
	if cats[0].Name != "foo" || cats[1].Name != "bar" {
		t.Errorf("%s != foo || %s != bar", cats[0].Name, cats[1].Name)
	}

}

func TestExpenseInsert(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

	u := model.User{1, ""}
	c := model.Category{1, ""}
	pm := model.PaymentMethod{1, ""}
	e := model.Expense{Date: time.Now().Local(), Who: &u, Method: &pm, Category: &c}
	_, err = s.ExpenseInsert(&e)
	if err != nil {
		t.Error(err)
	}

}

func TestExpensesGetAll(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

	s.UserInsert("M")
	s.UserInsert("A")

	s.CategoryInsert("Uncategorized")
	s.CategoryInsert("Spesa")
	s.CategoryInsert("Ristorante")

	s.PaymentMethodInsert("Contanti")
	s.PaymentMethodInsert("Bancomat")

	u := model.User{1, ""}
	c := model.Category{1, ""}
	pm := model.PaymentMethod{1, ""}
	e1 := model.Expense{DateCreated: time.Now().Local(), Date: time.Now().Local(), Who: &u, Method: &pm, Category: &c}
	s.ExpenseInsert(&e1)
	s.ExpenseInsert(&e1)
	s.ExpenseInsert(&e1)

	es, err := s.ExpensesGetN(2)
	if err != nil {
		t.Error(err)
	}

	for _, e := range es {
		log.Println(*e)
	}

}

func TestExpenseDelete(t *testing.T) {
	s, err := New("./test.db", true)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		s.Close()
		os.Remove("./test.db")
	}()

	s.UserInsert("M")
	s.UserInsert("A")

	s.CategoryInsert("Uncategorized")
	s.CategoryInsert("Spesa")
	s.CategoryInsert("Ristorante")

	s.PaymentMethodInsert("Contanti")
	s.PaymentMethodInsert("Bancomat")

	u := model.User{1, ""}
	c := model.Category{1, ""}
	pm := model.PaymentMethod{1, ""}
	e1 := model.Expense{DateCreated: time.Now().Local(), Date: time.Now().Local(), Who: &u, Method: &pm, Category: &c}
	e, _ := s.ExpenseInsert(&e1)

	err = s.ExpenseDelete(e.UUID)
	if err != nil {
		t.Error(err)
	}
	es, err := s.ExpensesGetN(2)
	if err != nil {
		t.Error(err)
	}
	for _, e := range es {
		log.Println(*e)
	}

}
