package mock

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/db"
	"ronche.se/moneytracker/domain"
)

func TestGetInsertTransaction(t *testing.T) {

	ms := NewMockStore()
	db.Populate(ms)

	id, err := ms.InsertTransaction(&domain.Transaction{
		Date:        domain.DateTime{time.Now()},
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.New(50, 0),
		FromID:      "",
		ToID:        "",
		Type:        "expense",
	})
	if err != nil {
		t.Fatal(err)
	}

	tx, _ := ms.GetTransaction(id)

	if tx.Description != "sds" {
		t.Fatal(tx)
	}

}

/*
func TestBalance(t *testing.T) {

	mTS := NewMockStore()
	db.Populate(mTS)

	acc1 := domain.Account{
		ID:      "Acc3",
		Owners:  []*domain.User{},
		Balance: decimal.New(100, 0),
	}

	acc2 := domain.Account{
		ID:      "Acc2",
		Owners:  []*domain.User{},
		Balance: decimal.New(100, 0),
	}

	mTS.InsertTransaction(&domain.Transaction{
		Date:        domain.Date{time.Now()},
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.New(30, 0),
		From:        &acc1,
		To:          &domain.WorldAcc,
		Type:        "expense",
	})
	mTS.InsertTransaction(&domain.Transaction{
		Date:        domain.Date{time.Now()},
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.New(-50, 0),
		From:        &acc2,
		To:          &acc1,
		Type:        "expense",
	})
	mTS.InsertTransaction(&domain.Transaction{
		Date:        domain.Date{time.Now()},
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.New(-50, 0),
		From:        &domain.WorldAcc,
		To:          &domain.WorldAcc,
		Type:        "expense",
	})

	tx := domain.Transaction{
		Date:        domain.Date{time.Now()},
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.New(-50, 0),
		From:        &domain.WorldAcc,
		To:          &acc2,
		Type:        "expense",
	}
	id, _ := mTS.InsertTransaction(&tx)
	if id != tx.ID {
		t.Fatal(tx)
	}

	tl, _ := mTS.GetTransactionsByAccount(acc1.ID)
	if len(tl) != 2 {
		t.Fatal(tl)
	}

	tl, _ = mTS.GetTransactionsByAccount(acc2.ID)
	if len(tl) != 2 {
		t.Fatal(tl)
	}

	tx.To = &acc1
	_, err := mTS.UpdateTransaction(&tx)
	if err != nil {
		t.Fatal(err)
	}

	tl, _ = mTS.GetTransactionsByAccount(acc1.ID)
	if len(tl) != 3 {
		t.Fatal(tl)
	}

	tl, _ = mTS.GetTransactionsByAccount(acc2.ID)
	if len(tl) != 1 {
		t.Fatal(tl)
	}

	mTS.DeleteTransaction(id)

	tl, _ = mTS.GetTransactionsByAccount(acc1.ID)
	if len(tl) != 2 {
		t.Fatal(tl)
	}

}
*/
