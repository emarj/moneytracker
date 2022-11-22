package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "embed"

	"github.com/shopspring/decimal"
	_ "modernc.org/sqlite"
	mt "ronche.se/moneytracker"
	"ronche.se/moneytracker/store/sqlite/queries"
)

//go:embed queries/schema.sql
var schema string

//go:embed queries/seeds.sql
var seeds string

type SQLiteStore struct {
	url    string
	create bool
	db     *sql.DB
}

func New(url string, create bool) *SQLiteStore {
	return &SQLiteStore{url: url, create: create}
}

func (s *SQLiteStore) Open() error {
	_, err := os.Stat(s.url)
	newDB := errors.Is(err, os.ErrNotExist)

	if newDB && !s.create {
		return err
	}

	db, err := sql.Open("sqlite", s.url)
	if err != nil {
		return err
	}
	s.db = db

	if newDB {
		err = s.Migrate()
		if err != nil {
			return err
		}

		err = s.Seed()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStore) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) Migrate() error {

	fmt.Println("Migrating...")
	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}

	return nil

}

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
		EntityID:    1,
		IsWorld:     true,
		IsSystem:    true,
	})
	query += queries.InsertBalance(mt.Balance{
		AccountID: 0,
		Timestamp: mt.DateTime{time.Unix(0, 0)},
		Value:     &decimal.Zero,
	})
	query += queries.InsertEntity(mt.Entity{
		ID:   1,
		Name: "'user1'",
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
	query += queries.InsertBalance(mt.Balance{
		AccountID: 2,
		Timestamp: mt.DateTime{time.Now().AddDate(0, 0, -1)},
		Value:     &decimal.Zero,
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

	query += queries.InsertOperation(mt.Operation{
		Timestamp:   &mt.DateTime{time.Now()},
		CreatedByID: 1,
		Description: "'Cena da Spalto 10'",
		CategoryID:  0,
	})

	query += queries.InsertTransaction(mt.Transaction{
		From:      mt.Account{ID: 2},
		To:        mt.Account{ID: 3},
		Amount:    decimal.New(99, 0),
		Operation: mt.Operation{ID: 1},
	})
	query += queries.InsertTransaction(mt.Transaction{
		From:      mt.Account{ID: 0},
		To:        mt.Account{ID: 0},
		Amount:    decimal.New(120, 0),
		Operation: mt.Operation{ID: 1},
	})

	query += queries.InsertOperation(mt.Operation{
		Timestamp:   &mt.DateTime{time.Now()},
		CreatedByID: 1,
		Description: "'Op2'",
		CategoryID:  0,
	})
	query += queries.InsertTransaction(mt.Transaction{
		From:      mt.Account{ID: 1},
		To:        mt.Account{ID: 2},
		Amount:    decimal.New(345, 0),
		Operation: mt.Operation{ID: 2},
	})
	query += queries.InsertTransaction(mt.Transaction{
		From:      mt.Account{ID: 1},
		To:        mt.Account{ID: 0},
		Amount:    decimal.New(43, 0),
		Operation: mt.Operation{ID: 2},
	})

	fmt.Println(query)

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	s.AddCredit(mt.Credit{
		Debtor:      mt.Entity{ID: 1},
		Creditor:    mt.Entity{ID: 2},
		Account:     mt.Account{ID: 3},
		Amount:      decimal.New(43, 0),
		Description: "first credit",
		Operation:   mt.Operation{ID: 1},
	})

	return nil
}
