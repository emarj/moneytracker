package sqlite

import mt "ronche.se/moneytracker"

func (s *SQLiteStore) GetEntities() ([]mt.Entity, error) {

	rows, err := s.db.Query("SELECT * FROM entities")
	if err != nil {
		return nil, err
	}

	entities := []mt.Entity{}
	var e mt.Entity

	for rows.Next() {
		if err = rows.Scan(&e.ID, &e.Name, &e.System); err != nil {
			return nil, err
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (s *SQLiteStore) GetEntity(eID int) (*mt.Entity, error) {

	row := s.db.QueryRow("SELECT * FROM entities WHERE id = ?", eID)

	var e mt.Entity

	if err := row.Scan(&e.ID, &e.Name, &e.System); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *SQLiteStore) AddEntity(e mt.Entity) (int, error) {

	res, err := s.db.Exec("INSERT INTO entities (name,system) VALUES(?,?)", e.Name, e.System)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}

// This does not sense, we should delete also all entity accounts (and transactions)
/*func (s *SQLiteStore) DeleteEntity(eID int) error {
	_, err := s.db.Exec("DELETE FROM entities WHERE id=?", eID)
	if err != nil {
		return err
	}

	return nil
}*/
