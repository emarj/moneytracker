package sqlite

import (
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) GetEntities() ([]mt.Entity, error) {

	rows, err := s.db.Query("SELECT * FROM entities")
	if err != nil {
		return nil, err
	}

	entities := []mt.Entity{}
	var e mt.Entity

	for rows.Next() {
		if err = rows.Scan(&e.ID, &e.Name, &e.System, &e.External); err != nil {
			return nil, err
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (s *SQLiteStore) GetEntity(eID int) (*mt.Entity, error) {

	row := s.db.QueryRow("SELECT * FROM entities WHERE id = ?", eID)

	var e mt.Entity

	if err := row.Scan(&e.ID, &e.Name, &e.System, &e.External); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *SQLiteStore) AddEntity(e mt.Entity) (null.Int, error) {

	id := null.Int{}
	res, err := s.db.Exec("INSERT INTO entities (id,name,is_system,is_external) VALUES(?,?,?,?)", e.ID, e.Name, e.System, e.External)
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

// This does not sense, we should delete also all entity accounts (and transactions)
/*func (s *SQLiteStore) DeleteEntity(eID int) error {
	_, err := s.db.Exec("DELETE FROM entities WHERE id=?", eID)
	if err != nil {
		return err
	}

	return nil
}*/
