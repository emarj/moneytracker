package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) Open(url string) error {
	db, err := sql.Open("sqlite", url)
	if err != nil {
		return err
	}

	s.db = db

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

	query += `CREATE TABLE "recipes" (
		"id"	INTEGER,
		"date_created"	TEXT DEFAULT (datetime('now', 'localtime')),
		"date_modified"	TEXT DEFAULT (datetime('now', 'localtime')),
		"body"	TEXT NOT NULL DEFAULT '{}',
		"author_id" INTEGER GENERATED ALWAYS AS (IFNULL(json_extract(body, '$.author_id'),0)) STORED NOT NULL,
		"json" TEXT GENERATED ALWAYS AS (json_patch(body,json_object('id',id,'date_created',date_created,'date_modified',date_modified))) STORED,
		PRIMARY KEY("id")
	);`

	query += `CREATE TRIGGER update_date
			AFTER UPDATE ON recipes
 			BEGIN
				UPDATE recipes SET date_modified = datetime('now', 'localtime')
	 			WHERE id = NEW.id;
 			END;`

	query += `CREATE TRIGGER strip_json_on_update
			 AFTER UPDATE OF body ON recipes
			  BEGIN
				 UPDATE recipes SET body = json_remove(body, '$.id','$.date_created','$.date_modified')
				  WHERE id = NEW.id;
			  END;
			  CREATE TRIGGER strip_json_on_insert
			  AFTER INSERT ON recipes
			   BEGIN
				  UPDATE recipes SET body = json_remove(body, '$.id','$.date_created','$.date_modified')
				   WHERE id = NEW.id;
			   END;`

	query += `CREATE TABLE "authors" (
				"id"	INTEGER,
				"name"	TEXT NOT NULL,
				"type" TEXT,
				"desc" TEXT,
				PRIMARY KEY("id")
			);`

	query += `CREATE TABLE "extra_fields" (
				"id"	INTEGER,
				"name"	TEXT NOT NULL,
				"desc" TEXT,
				PRIMARY KEY("id")
			);`

	fmt.Println("Executing:\n " + query)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}
