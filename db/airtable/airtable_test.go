package airtable

import (
	"testing"
	"time"

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
		FromID:      "sasas",
		ToID:        "asasa",
	})
	if err != nil {
		t.Fatal(err)
	}

}
