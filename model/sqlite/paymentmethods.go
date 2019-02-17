package sqlite

import (
	"ronche.se/moneytracker/model"
)

func (s *sqlite) PaymentMethodsGetAll() ([]*model.PaymentMethod, error) {
	rows, err := s.db.Query("SELECT * FROM paymentmethods")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pms []*model.PaymentMethod

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		pms = append(pms, &model.PaymentMethod{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return pms, nil
}

func (s *sqlite) PaymentMethodInsert(name string) (*model.PaymentMethod, error) {
	stmt, err := s.db.Prepare("INSERT INTO paymentmethods(pm_name) VALUES(?)")
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
	return &model.PaymentMethod{int(lastID), name}, nil
}
