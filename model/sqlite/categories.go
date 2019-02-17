package sqlite

import (
	"ronche.se/moneytracker/model"
)

func (s *sqlite) CategoriesGetAll() ([]*model.Category, error) {
	rows, err := s.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*model.Category

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		cats = append(cats, &model.Category{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *sqlite) CategoryInsert(name string) (*model.Category, error) {
	stmt, err := s.db.Prepare("INSERT INTO categories(cat_name) VALUES(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &model.Category{int(lastID), name}, nil
}
