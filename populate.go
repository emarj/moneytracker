package moneytracker

import (
	"fmt"

	tt "github.com/emarj/moneytracker/datetime/test"
	"github.com/shopspring/decimal"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/guregu/null.v4"
)

func Populate(s Store) error {
	fmt.Print("Populating...")
	var err error

	tt.Init()

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

	user1 := User{
		Name:        "arianna",
		DisplayName: "Arianna",
		IsAdmin:     false,
	}
	err = s.RegisterUser(&user1, "prova")
	if err != nil {
		return err
	}

	user2 := User{
		Name:        "marco",
		DisplayName: "Marco",
		IsAdmin:     false,
	}
	err = s.RegisterUser(&user2, "pippo")
	if err != nil {
		return err
	}

	entUser1 := Entity{
		ID:         null.IntFrom(1),
		Name:       "arianna",
		IsSystem:   false,
		IsExternal: false,
	}

	err = s.AddEntity(&entUser1)
	if err != nil {
		return err
	}

	err = s.AddSharesForEntity(EntityShare{
		UserID:   user1.ID.Int64,
		EntityID: entUser1.ID.Int64,
		Quota:    100,
	})
	if err != nil {
		return err
	}

	entUser2 := Entity{
		ID:         null.IntFrom(2),
		Name:       "marco",
		IsSystem:   false,
		IsExternal: false,
	}

	err = s.AddEntity(&entUser2)
	if err != nil {
		return err
	}
	err = s.AddSharesForEntity(EntityShare{
		UserID:   user2.ID.Int64,
		EntityID: entUser2.ID.Int64,
		Quota:    100,
	})
	if err != nil {
		return err
	}

	entUser3 := Entity{
		ID:         null.IntFrom(3),
		Name:       "am",
		IsSystem:   false,
		IsExternal: false,
	}

	err = s.AddEntity(&entUser3)
	if err != nil {
		return err
	}
	err = s.AddSharesForEntity(EntityShare{
		UserID:   user1.ID.Int64,
		EntityID: entUser3.ID.Int64,
		Quota:    50,
	}, EntityShare{
		UserID:   user2.ID.Int64,
		EntityID: entUser3.ID.Int64,
		Quota:    50,
	})
	if err != nil {
		return err
	}

	accounts := orderedmap.New[string, Account]()

	accounts.Store("user1:cash", Account{
		Name:        "contanti",
		DisplayName: "Contanti",
		OwnerID:     entUser1.ID.Int64,
	})

	accounts.Store("user1:cc1", Account{
		Name:        "conto_corrente",
		DisplayName: "Conto Corrente",
		OwnerID:     entUser1.ID.Int64,
	})
	accounts.Store("user1:cc2", Account{
		Name:        "conto_corrente_posta",
		DisplayName: "Conto Banco Posta",
		OwnerID:     entUser1.ID.Int64,
	})
	accounts.Store("user2:cc", Account{
		Name:        "conto_corrente",
		DisplayName: "Conto Corrente",
		OwnerID:     entUser2.ID.Int64,
	})
	accounts.Store("user2:cash", Account{
		Name:        "contanti",
		DisplayName: "Contanti",
		OwnerID:     entUser2.ID.Int64,
	})
	accounts.Store("user1:credits", Account{
		Name:        "credits",
		DisplayName: "Crediti",
		OwnerID:     entUser1.ID.Int64,
		TypeID:      AccTypeCredit,
	})
	accounts.Store("user1:credits2", Account{
		Name:        "credits2",
		DisplayName: "Crediti2",
		OwnerID:     entUser1.ID.Int64,
		TypeID:      AccTypeCredit,
	})
	accounts.Store("user2:credits", Account{
		Name:        "credits",
		DisplayName: "Crediti",
		OwnerID:     entUser2.ID.Int64,
		TypeID:      AccTypeCredit,
	})
	accounts.Store("user1:investments", Account{
		Name:        "investments",
		DisplayName: "Investments",
		OwnerID:     entUser1.ID.Int64,
		TypeID:      AccTypeInvestment,
	})
	accounts.Store("user2:investments", Account{
		Name:        "investments",
		DisplayName: "Investments",
		OwnerID:     entUser2.ID.Int64,
		TypeID:      AccTypeInvestment,
	})
	accounts.Store("user3:comune", Account{
		Name:        "cassa_comune",
		DisplayName: "Cassa Comune",
		OwnerID:     entUser3.ID.Int64,
	})

	for pair := accounts.Oldest(); pair != nil; pair = pair.Next() {
		err := s.AddAccount(&pair.Value)
		if err != nil {
			return err
		}
	}

	err = s.SetBalance(&Balance{
		AccountID: accounts.Value("user1:cc1").ID,
		Timestamp: tt.BEFORE,
		Value:     decimal.NewFromInt(4000),
		Comment:   "Initial balance",
	})
	if err != nil {
		return err
	}

	err = s.SetBalance(&Balance{
		AccountID: accounts.Value("user1:cc1").ID,
		Timestamp: tt.Now,
		Value:     decimal.NewFromInt(5000),
	})
	if err != nil {
		return err
	}

	operations := []Operation{
		{
			Description: "Cena Fuori in 2",
			Transactions: []Transaction{
				{
					Timestamp: tt.Before,
					FromID:    accounts.Value("user1:cc1").ID.Int64,
					ToID:      0,
					Amount:    decimal.New(801, -1),
				},
				{
					Timestamp: tt.Before,
					FromID:    accounts.Value("user2:credits").ID.Int64,
					ToID:      accounts.Value("user1:credits").ID.Int64,
					Amount:    decimal.New(40, 0),
				},
			},
			TypeID:     OpTypeExpense,
			CategoryID: 0,
		},
		{
			Description: "Giroconto",
			Transactions: []Transaction{
				{
					Timestamp: tt.Before,
					FromID:    accounts.Value("user1:cc1").ID.Int64,
					ToID:      accounts.Value("user1:cc2").ID.Int64,
					Amount:    decimal.New(345, 0),
				},
			},
			CategoryID: 2,
		},
		{
			Description: "Prestito per acquisto",
			Transactions: []Transaction{
				{
					Timestamp: tt.Before,
					FromID:    accounts.Value("user1:cash").ID.Int64,
					ToID:      accounts.Value("user2:cash").ID.Int64,
					Amount:    decimal.New(100, 0),
				},
				{
					Timestamp: tt.Before,
					FromID:    accounts.Value("user2:credits").ID.Int64,
					ToID:      accounts.Value("user1:credits").ID.Int64,
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
