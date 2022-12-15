package sqlite

import (
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"

	jt "ronche.se/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetEntities() ([]mt.Entity, error) {

	stmt := jet.SELECT(
		jt.Entity.AllColumns,
	).FROM(jt.Entity)

	entities := []mt.Entity{}
	err := stmt.Query(s.db, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (s *SQLiteStore) GetEntity(eID int) (*mt.Entity, error) {

	stmt := jet.SELECT(
		jt.Entity.AllColumns,
	).FROM(jt.Entity).WHERE(
		jt.Entity.ID.EQ(jet.Int(int64(eID))),
	)

	dest := &mt.Entity{}
	err := stmt.Query(s.db, dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (s *SQLiteStore) AddEntity(e mt.Entity) (null.Int, error) {

	id := null.Int{}
	stmt := jt.Entity.INSERT(jt.Entity.AllColumns).RETURNING(jt.Entity.AllColumns).MODEL(e)

	res, err := stmt.Exec(s.db)
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
	_, err := s.db.Exec("DELETE FROM entity WHERE id=?", eID)
	if err != nil {
		return err
	}

	return nil
}*/
