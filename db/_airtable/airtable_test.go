package airtable

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/domain"
)

func TestInsertTransaction(t *testing.T) {

	decimal.MarshalJSONWithoutQuotes = true

	const apiKey = "keyAuR8F3wLAUXZAL"

	at := NewAirtable(apiKey)

	_, err := at.InsertTransaction(&domain.Transaction{
		Date:        domain.DateTime{time.Now()},
		Description: "sds",
		Notes:       "fdf",
		Amount:      decimal.NewFromInt(50),
		FromAccount: uuid.Must(uuid.NewV4()),
		ToAccount:   uuid.Must(uuid.NewV4()),
	})
	if err != nil {
		t.Fatal(err)
	}

}
