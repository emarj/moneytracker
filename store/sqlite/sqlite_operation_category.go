package sqlite

import (
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) GetCategories() ([]mt.Category, error) {

	rows, err := s.db.Query("SELECT * FROM category")
	if err != nil {
		return nil, err
	}

	categories := []mt.Category{}
	var c mt.Category

	for rows.Next() {
		if err = rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func (s *SQLiteStore) AddCategory(c mt.Category) (null.Int, error) {

	id := null.Int{}
	res, err := s.db.Exec("INSERT INTO category (id,name) VALUES(?,?)", c.ID, c.Name)
	if err != nil {
		return id, err
	}

	id.Int64, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	id.Valid = true
	return id, nil

}
