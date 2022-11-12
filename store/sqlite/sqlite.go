package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "embed"

	_ "modernc.org/sqlite"
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

	var query string

	query += schema
	query += seeds

	fmt.Println("Executing:\n " + query)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}
