package sqlite

import (
	"database/sql"
	"fmt"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed queries/schema.sql
var schema string

type SQLiteStore struct {
	url    string
	create bool
	db     *sql.DB
}

func New(url string, create bool) *SQLiteStore {
	return &SQLiteStore{url: url, create: create}
}

func (s *SQLiteStore) Open() error {
	db, err := sql.Open("sqlite", s.url)
	if err != nil {
		return err
	}
	s.db = db

	_, err = s.db.Exec(`PRAGMA journal_mode = WAL;
						PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	if s.create {
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
