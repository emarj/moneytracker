package sqlite

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
	"ronche.se/moneytracker/datetime"
)

func (s *SQLiteStore) Seed() error {
	fmt.Println("Seeding...")
	var err error

	entUser1 := mt.Entity{
		ID:         null.IntFrom(1),
		Name:       "arianna",
		IsSystem:   false,
		IsExternal: false,
	}

	err = s.AddEntity(&entUser1)
	if err != nil {
		return err
	}

	entUser2 := mt.Entity{
		ID:         null.IntFrom(2),
		Name:       "marco",
		IsSystem:   false,
		IsExternal: false,
	}

	err = s.AddEntity(&entUser2)
	if err != nil {
		return err
	}
	entUser3 := mt.Entity{
		ID:         null.IntFrom(3),
		Name:       "am",
		IsSystem:   false,
		IsExternal: false,
	}

	err = s.AddEntity(&entUser3)
	if err != nil {
		return err
	}

	var accounts map[string]mt.Account = map[string]mt.Account{
		"user1:cash": {
			Name:        "contanti",
			DisplayName: "Contanti",
			Owner:       entUser1,
		},
		"user1:cc1": {
			ID:          null.IntFrom(1004),
			Name:        "conto_corrente",
			DisplayName: "Conto Corrente",
			Owner:       entUser1,
		},
		"user1:cc2": {
			Name:        "conto_corrente_posta",
			DisplayName: "Conto Banco Posta",
			Owner:       entUser1,
		},
		"user2:cc": {
			Name:        "conto_corrente",
			DisplayName: "Conto Corrente",
			Owner:       entUser2,
		},
		"user2:cash": {
			Name:        "contanti",
			DisplayName: "Contanti",
			Owner:       entUser2,
		},
		"user1:credits": {
			Name:        "credits",
			DisplayName: "Crediti",
			Owner:       entUser1,
			Type:        mt.AccountCredit,
		},
		"user2:credits": {
			Name:        "credits",
			DisplayName: "Crediti",
			Owner:       entUser2,
			Type:        mt.AccountCredit,
		},
		"user3:comune": {
			Name:        "cassa_comune",
			DisplayName: "Cassa Comune",
			Owner:       entUser3,
		},
	}

	for k, a := range accounts {
		err := s.AddAccount(&a, nil)
		if err != nil {
			return err
		}
		accounts[k] = a
	}

	err = s.SetBalance(mt.Balance{
		AccountID: accounts["user1:cc1"].ID,
		Timestamp: datetime.FromTime(time.Now().AddDate(0, 0, -3)),
		Value:     decimal.NewFromInt(4000),
	})
	if err != nil {
		return err
	}

	operations := []mt.Operation{
		{
			Description: "Cena Fuori in 2",
			Transactions: []mt.Transaction{
				{Timestamp: datetime.FromTime(time.Now().AddDate(0, 0, -1)), From: accounts["user1:cc1"], To: mt.Account{ID: null.IntFrom(0)}, Amount: decimal.New(80, 0)},
				{Timestamp: datetime.FromTime(time.Now().AddDate(0, 0, -1)), From: accounts["user2:credits"], To: accounts["user1:credits"], Amount: decimal.New(40, 0)}},
			TypeID:     mt.OpTypeExpense,
			CategoryID: 0,
		},
		{
			Description: "Giroconto",
			Transactions: []mt.Transaction{
				{
					Timestamp: datetime.FromTime(time.Now()),
					From:      accounts["user1:cc1"],
					To:        accounts["user1:cc2"],
					Amount:    decimal.New(345, 0),
				},
			},
			CategoryID: 2,
		},
		{
			Description: "Prestito 100 Euro a Marco",
			Transactions: []mt.Transaction{
				{
					Timestamp: datetime.FromTime(time.Now().AddDate(0, -1, 0)),
					From:      accounts["user1:cash"],
					To:        accounts["user2:cash"],
					Amount:    decimal.New(100, 0),
				},
				{
					Timestamp: datetime.FromTime(time.Now().AddDate(0, -1, 0)),
					From:      accounts["user2:credits"],
					To:        accounts["user1:credits"],
					Amount:    decimal.New(100, 0),
				},
			},
			CategoryID: 1,
		},
	}

	for k, op := range operations {
		err = s.AddOperation(&op)
		if err != nil {
			return err
		}
		operations[k] = op
	}

	return nil
}
