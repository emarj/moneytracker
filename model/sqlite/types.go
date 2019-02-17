package sqlite

import (
	"ronche.se/moneytracker/model"
)

func (s *sqlite) TypesGetAll() ([]*model.Type, error) {
	rows, err := s.db.Query("SELECT * FROM types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tps []*model.Type

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		tps = append(tps, &model.Type{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tps, nil
}

func (s *sqlite) TypeInsert(name string) (*model.Type, error) {
	stmt, err := s.db.Prepare("INSERT INTO types(type_name) VALUES(?)")
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
	return &model.Type{int(lastID), name}, nil
}
