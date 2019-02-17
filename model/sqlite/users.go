package sqlite

import (
	"log"

	"ronche.se/moneytracker/model"
)

func (s *sqlite) UsersGetAll() ([]*model.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var us []*model.User

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		us = append(us, &model.User{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (s *sqlite) UserInsert(name string) (*model.User, error) {
	stmt, err := s.db.Prepare("INSERT INTO users(user_name) VALUES(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &model.User{int(lastID), name}, nil
}
