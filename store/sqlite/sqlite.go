package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed queries/schema.sql
var schema string

type TXDB interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type SQLiteStore struct {
	dsn     string
	migrate bool
	db      *sql.DB
}

func New(dsn string, migrate bool) *SQLiteStore {
	return &SQLiteStore{dsn: dsn, migrate: migrate}
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
