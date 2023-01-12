package sqlite

import (
	"errors"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"

	"github.com/go-jet/jet/v2/qrm"
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

func (s *SQLiteStore) GetEntity(eID int64) (*mt.Entity, error) {

	stmt := jet.SELECT(
		jt.Entity.AllColumns,
	).FROM(jt.Entity).WHERE(
		jt.Entity.ID.EQ(jet.Int(eID)),
	)

	e := &mt.Entity{}
	err := stmt.Query(s.db, e)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, mt.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (s *SQLiteStore) AddEntity(e *mt.Entity) error {

	stmt := jt.Entity.INSERT(jt.Entity.AllColumns).RETURNING(jt.Entity.AllColumns).MODEL(e)

	err := stmt.Query(s.db, e)
	if err != nil {
		return err
	}

	return nil
}

// This does not sense, we should delete also all entity accounts (and transactions)
/*func (s *SQLiteStore) DeleteEntity(eID int) error {
	_, err := s.db.Exec("DELETE FROM entity WHERE id=?", eID)
	if err != nil {
		return err
	}

	return nil
}*/
