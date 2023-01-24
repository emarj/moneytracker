package main

import (
	"flag"
	"log"
	"path"

	"github.com/emarj/moneytracker/store/sqlite"
)

func main() {
	var s *sqlite.SQLiteStore

	var dir = flag.String("dir", "./data", "")
	var dbName = flag.String("db", "moneytracker.sqlite", "")
	flag.Parse()

	dsn := path.Join(*dir, *dbName)

	s = sqlite.New(dsn, true)

	err := s.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = populate(s)
	if err != nil {
		log.Fatal(err)
	}
}
