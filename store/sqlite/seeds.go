package sqlite

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	mt "ronche.se/moneytracker"
	"ronche.se/moneytracker/store/sqlite/queries"
)

func (s *SQLiteStore) Seed() error {
	fmt.Println("Seeding...")

	var query string
	query += queries.InsertEntity(mt.Entity{
		ID:     0,
		Name:   "'_system'",
		System: true,
	})
	query += queries.InsertAccount(mt.Account{
		ID:          0,
		Name:        "'world'",
		DisplayName: "'World'",
		EntityID:    0, //this is needed for now in order to correctly compute expenses/income in fronted
		IsWorld:     true,
		IsSystem:    true,
	})

	query += queries.InsertEntity(mt.Entity{
		ID:   1,
		Name: "'user1'",
	})
	query += queries.InsertEntity(mt.Entity{
		ID:   2,
		Name: "'user2'",
	})

	query += queries.InsertAccount(mt.Account{
		ID:          1,
		Name:        "'acc1'",
		DisplayName: "'Account 1'",
		EntityID:    1,
	})
	value2 := decimal.New(2000, 0)
	query += queries.InsertBalance(mt.Balance{
		AccountID: 1,
		Timestamp: mt.DateTime{time.Now().AddDate(0, 0, -1)},
		Value:     &value2,
	})
	query += queries.InsertAccount(mt.Account{
		ID:          2,
		Name:        "'acc2'",
		DisplayName: "'Account 2'",
		EntityID:    1,
	})
	query += queries.InsertAccount(mt.Account{
		ID:          3,
		Name:        "'acc3'",
		DisplayName: "'Account 3'",
		EntityID:    1,
	})
	value1 := decimal.New(1000, 0)
	query += queries.InsertBalance(mt.Balance{
		AccountID: 3,
		Timestamp: mt.DateTime{time.Now().AddDate(0, 0, -1)},
		Value:     &value1,
	})

	query += queries.InsertAccount(mt.Account{
		ID:          1001,
		Name:        "'credits_user2'",
		DisplayName: "'Credits User 2'",
		EntityID:    1,
		IsCredit:    true,
	})

	query += queries.InsertAccount(mt.Account{
		ID:          1002,
		Name:        "'credits_user1'",
		DisplayName: "'Credits User 1'",
		EntityID:    2,
		IsCredit:    true,
	})

	query += queries.InsertAccount(mt.Account{
		ID:          4,
		Name:        "'acc4'",
		DisplayName: "'Account 4'",
		EntityID:    2,
	})

	fmt.Println(query)

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	_, err = s.AddOperation(mt.Operation{
		CreatedByID: 1,
		Timestamp:   &mt.DateTime{time.Now()},
		Description: "Cena Fuori",
		Transactions: []mt.Transaction{
			{
				From:   mt.Account{ID: 2},
				To:     mt.Account{ID: 0},
				Amount: decimal.New(80, 0),
			},
			{
				From:   mt.Account{ID: 1002},
				To:     mt.Account{ID: 1001},
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
		Description: "Operation 2",
		Transactions: []mt.Transaction{
			{
				From:   mt.Account{ID: 1},
				To:     mt.Account{ID: 2},
				Amount: decimal.New(345, 0),
			},
			{
				From:   mt.Account{ID: 1},
				To:     mt.Account{ID: 0},
				Amount: decimal.New(43, 0),
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
				From:   mt.Account{ID: 1},
				To:     mt.Account{ID: 4},
				Amount: decimal.New(100, 0),
			},
			{
				From:   mt.Account{ID: 1002},
				To:     mt.Account{ID: 1001},
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
