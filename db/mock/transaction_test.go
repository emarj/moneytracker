package mock

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/domain"
)

func populateTransactions(ts *mockTransactionStore) {
	ts.AddTransaction(&domain.Transaction{
		Date:        time.Now(),
		Description: "Transfer",
		Notes:       "",
		Amount:      decimal.New(30, 0),
		FromID:      uuid.Must(uuid.NewV4()),
		ToID:        uuid.Must(uuid.NewV4()),
	})

	ts.AddTransaction(&domain.Transaction{
		Date:        time.Now(),
		Description: "TX2",
		Notes:       "",
		Amount:      decimal.New(120, 0),
		FromID:      uuid.Must(uuid.NewV4()),
		ToID:        uuid.Must(uuid.NewV4()),
	})
}

func TestGetInsertTransaction(t *testing.T) {

	ms := newMockTransactionStore()

	populateTransactions(ms)

	tx := &domain.Transaction{
		Date:        time.Now(),
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.New(50, 0),
		FromID:      uuid.Must(uuid.NewV4()),
		ToID:        uuid.Must(uuid.NewV4()),
	}

	err := ms.AddTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}

	tx1, _ := ms.GetTransaction(tx.ID)

	if tx1.Description != "sds" {
		t.Fatal(tx)
	}

}
