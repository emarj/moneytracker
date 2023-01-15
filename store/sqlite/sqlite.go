package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	_ "modernc.org/sqlite"

	mt "github.com/emarj/moneytracker"
)

//go:embed queries/schema.sql
var schema string

type TXDB interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type SQLiteStore struct {
	dsn     string
	migrate bool
	db      *sql.DB
	//
	accountTypes   []mt.AccountType
	operationTypes []mt.OperationType
}

func New(filename string, migrate bool) *SQLiteStore {
	return &SQLiteStore{dsn: filename + "?cache=shared", migrate: migrate}
}

func NewTemp() *SQLiteStore {
	return New("file::memory:", true)
}

func (s *SQLiteStore) Open() error {
	db, err := sql.Open("sqlite", s.dsn)
	if err != nil {
		return err
	}
	s.db = db

	_, err = s.db.Exec(`PRAGMA journal_mode = WAL;
						PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	if s.migrate {
		err = s.Migrate()
		if err != nil {
			return err
		}

		err = s.Seeding()
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

	fmt.Print("Migrating...")
	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}

	fmt.Println("OK")

	return nil

}

func Btoi(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
