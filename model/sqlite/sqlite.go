package sqlite

import (
	"log"
	"os"

	"github.com/iancoleman/strcase"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	//sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

var schema = [...]string{

	`CREATE TABLE IF NOT EXISTS transactions (
		uuid	TEXT NOT NULL,
		date_created	TEXT NOT NULL,
		date	TEXT NOT NULL,
		type_id	INTEGER NOT NULL,
		user_id	INTEGER NOT NULL,
		amount	NUMERIC NOT NULL,
		description	TEXT NOT NULL,
		pm_id	INTEGER,
		shared	INTEGER NOT NULL,
		shared_quota	NUMERIC NOT NULL,
		geolocation TEXT,
		cat_id	INTEGER NOT NULL,
		PRIMARY KEY(uuid)
)`,

	`CREATE TABLE IF NOT EXISTS types ( 
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL, 
	PRIMARY KEY(id)
)`,

	`CREATE TABLE IF NOT EXISTS users ( 
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL, 
	PRIMARY KEY(id)
)`,

	`CREATE TABLE IF NOT EXISTS categories ( 
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL, 
	PRIMARY KEY(id)
)`,

	`CREATE TABLE IF NOT EXISTS paymentmethods (
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL,
	PRIMARY KEY(id)
)`}

type sqlite struct {
	db *sqlx.DB
}

func New(path string, create bool) (*sqlite, error) {

	if !create {
		_, err := os.Open(path)
		if err != nil {
			log.Fatalf("impossible to open the db file: %v", err)
		}
	}
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("json", strcase.ToSnake)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if create {
		for i := range schema {
			_, err = db.Exec(schema[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return &sqlite{db}, nil

}

func (s *sqlite) Close() error {
	return s.db.Close()
}
