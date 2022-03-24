package mock

import (
	"testing"

	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/db"
)

func TestGetAccount(t *testing.T) {

	ms := NewMockStore()

	db.Populate(ms)

	a, err := ms.GetAccount("primary")
	if err != nil {
		t.Fatal(a)
	}

	_, err = ms.GetAccount("shared")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ms.GetAccount("NotExistent")
	if err == nil {
		t.Fatal(err)
	}
}

func TestGetAccounts(t *testing.T) {

	ms := NewMockStore()
	db.Populate(ms)

	al, _ := ms.GetAccountsOfUser("nonexistentuser")
	if len(al) != 0 {
		t.Fatal(al)
	}

	al, _ = ms.GetAccountsOfUser("marco")
	if len(al) != 3 {
		t.Fatal(al)
	}

	al, _ = ms.GetAccountsOfUser("arianna")
	if len(al) != 2 {
		t.Fatal(al)
	}
}

func TestBalance(t *testing.T) {

	ms := NewMockStore()
	db.Populate(ms)

	acc, _ := ms.GetAccount("primary")
	b := acc.Balance

	ms.Balance("primary", decimal.New(30, 0))

	if !acc.Balance.Equals(b.Add(decimal.New(30, 0))) {
		t.Fatal(acc)
	}
}
