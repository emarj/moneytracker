package sqlite

import (
	"fmt"
	"time"

	mt "github.com/emarj/moneytracker"
	"github.com/emarj/moneytracker/datetime"
	"github.com/shopspring/decimal"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/guregu/null.v4"
)

func (s *SQLiteStore) Seed() error {
	fmt.Print("Seeding...")
	var err error

	// Add Categories
	categories := []string{
		"Spesa",
		"Utenze",
		"Utenze/Energia",
		"Utenze/Internet",
		"Utenze/Telefonia",
		"Salute",
		"Ristoranti - Bar",
		"Sport",
		"Sport/Tennis",
		"Trasporti",
		"Tasse",
		"Tasse/Rifiuti",
		"Regali",
		"Viaggi",
	}

	for _, n := range categories {
		_, err = s.AddCategory(n)
		if err != nil {
			return err
		}
	}

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

	accounts := orderedmap.New[string, mt.Account]()

	accounts.Store("user1:cash", mt.Account{
		Name:        "contanti",
		DisplayName: "Contanti",
		Owner:       entUser1,
	})

	accounts.Store("user1:cc1", mt.Account{
		Name:        "conto_corrente",
		DisplayName: "Conto Corrente",
		Owner:       entUser1,
	})
	accounts.Store("user1:cc2", mt.Account{
		Name:        "conto_corrente_posta",
		DisplayName: "Conto Banco Posta",
		Owner:       entUser1,
	})
	accounts.Store("user2:cc", mt.Account{
		Name:        "conto_corrente",
		DisplayName: "Conto Corrente",
		Owner:       entUser2,
	})
	accounts.Store("user2:cash", mt.Account{
		Name:        "contanti",
		DisplayName: "Contanti",
		Owner:       entUser2,
	})
	accounts.Store("user1:credits", mt.Account{
		Name:        "credits",
		DisplayName: "Crediti",
		Owner:       entUser1,
		TypeID:      mt.AccountCredit,
	})
	accounts.Store("user2:credits", mt.Account{
		Name:        "credits",
		DisplayName: "Crediti",
		Owner:       entUser2,
		TypeID:      mt.AccountCredit,
	})
	accounts.Store("user3:comune", mt.Account{
		Name:        "cassa_comune",
		DisplayName: "Cassa Comune",
		Owner:       entUser3,
	})

	for pair := accounts.Oldest(); pair != nil; pair = pair.Next() {
		err := s.AddAccount(&pair.Value)
		if err != nil {
			return err
		}
	}

	/* err = s.SetBalance(mt.Balance{
		AccountID: accounts.Value("user1:cc1").ID,
		ValueAt: mt.ValueAt{
			Timestamp: datetime.FromTime(time.Now().AddDate(0, 0, -3)),
			Value:     decimal.NewFromInt(4000),
		},
	})
	if err != nil {
		return err
	} */

	operations := []mt.Operation{
		{
			Description: "Cena Fuori in 2",
			Transactions: []mt.Transaction{
				{Timestamp: datetime.FromTime(time.Now().AddDate(0, 0, -1)), From: accounts.Value("user1:cc1"), To: mt.Account{ID: null.IntFrom(0)}, Amount: decimal.New(80, 0)},
				{Timestamp: datetime.FromTime(time.Now().AddDate(0, 0, -1)), From: accounts.Value("user2:credits"), To: accounts.Value("user1:credits"), Amount: decimal.New(40, 0)}},
			TypeID:     mt.OpTypeExpense,
			CategoryID: 0,
		},
		{
			Description: "Giroconto",
			Transactions: []mt.Transaction{
				{
					Timestamp: datetime.FromTime(time.Now()),
					From:      accounts.Value("user1:cc1"),
					To:        accounts.Value("user1:cc2"),
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
					From:      accounts.Value("user1:cash"),
					To:        accounts.Value("user2:cash"),
					Amount:    decimal.New(100, 0),
				},
				{
					Timestamp: datetime.FromTime(time.Now().AddDate(0, -1, 0)),
					From:      accounts.Value("user2:credits"),
					To:        accounts.Value("user1:credits"),
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

	fmt.Println("OK")

	return nil
}
