package sqlite

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) Seed() error {
	fmt.Println("Seeding...")
	var err error

	entSystem := mt.Entity{
		ID:       null.IntFrom(0),
		Name:     "_system",
		System:   true,
		External: false,
	}

	_, err = s.AddEntity(entSystem)
	if err != nil {
		return err
	}

	entUser1 := mt.Entity{
		ID:       null.IntFrom(1),
		Name:     "arianna",
		System:   false,
		External: false,
	}

	_, err = s.AddEntity(entUser1)
	if err != nil {
		return err
	}

	entUser2 := mt.Entity{
		ID:       null.IntFrom(2),
		Name:     "marco",
		System:   false,
		External: false,
	}

	_, err = s.AddEntity(entUser2)
	if err != nil {
		return err
	}

	var accounts map[string]mt.Account = map[string]mt.Account{
		"world": {
			ID:          null.IntFrom(0),
			Name:        "world",
			DisplayName: "World",
			Owner:       entSystem, //this is needed for now in order to correctly compute expenses/income in fronted
			IsWorld:     true,
			IsSystem:    true,
		},
		"user1:cash": {
			Name:        "contanti",
			DisplayName: "Contanti",
			Owner:       entUser1,
		},
		"user1:cc1": {
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
		},
		"user2:credits": {
			Name:        "credits",
			DisplayName: "Crediti",
			Owner:       entUser2,
		},
	}

	for k, a := range accounts {
		id, err := s.AddAccount(a)
		if err != nil {
			return err
		}
		if a.ID.Valid && !id.Equal(a.ID) {
			log.Warnf("seeding: insert account %s: expecting id %d, got %d", a.Name, a.ID.Int64, id.Int64)
		}
		a.ID = id
		accounts[k] = a
	}

	_, err = s.AddOperation(mt.Operation{
		CreatedByID: 1,
		Timestamp:   &mt.DateTime{time.Now()},
		Description: "Cena Fuori in 2",
		Transactions: []mt.Transaction{
			{
				From:   accounts["user1:cc1"],
				To:     accounts["world"],
				Amount: decimal.New(80, 0),
			},
			{
				From:   accounts["user2:credits"],
				To:     accounts["user1:credits"],
				Amount: decimal.New(40, 0),
			},
		},
		CategoryID: 0,
	})
	if err != nil {
		return err
	}

	_, err = s.AddOperation(mt.Operation{
		CreatedByID: 1,
		Timestamp:   &mt.DateTime{time.Now()},
		Description: "Giroconto",
		Transactions: []mt.Transaction{
			{
				From:   accounts["user1:cc1"],
				To:     accounts["user1:cc2"],
				Amount: decimal.New(345, 0),
			},
		},
		CategoryID: 0,
	})
	if err != nil {
		return err
	}

	_, err = s.AddOperation(mt.Operation{
		CreatedByID: 1,
		Timestamp:   &mt.DateTime{time.Now()},
		Description: "Prestito 100 Euro A -> M",
		Transactions: []mt.Transaction{
			{
				From:   accounts["user1:cash"],
				To:     accounts["user2:cash"],
				Amount: decimal.New(100, 0),
			},
			{
				From:   accounts["user2:credits"],
				To:     accounts["user1:credits"],
				Amount: decimal.New(100, 0),
			},
		},
		CategoryID: 0,
	})
	if err != nil {
		return err
	}

	return nil
}
