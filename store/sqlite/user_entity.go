package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"
	jt "github.com/emarj/moneytracker/.gen/table"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func insertEntityShares(txdb TXDB, shares ...mt.EntityShare) error {
	_, err := jt.EntityShare.INSERT(jt.EntityShare.AllColumns).
		MODELS(&shares).
		Exec(txdb)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) AddSharesForEntity(shares ...mt.EntityShare) error {
	total := int64(0)
	for _, es := range shares {
		total += es.Quota
		if es.EntityID != shares[0].EntityID {
			return fmt.Errorf("shares must be referred to the same entity")
		}
	}

	if total != 100 {
		return fmt.Errorf("the sum of quotas should be 100, got %d", total)
	}

	return insertEntityShares(s.db, shares...)
}

func (s *SQLiteStore) GetUserWithShares(uID int64) (mt.User, error) {
	stmt := jet.SELECT(
		jt.User.AllColumns.Except(jt.User.Password),
		jt.EntityShare.AllColumns,
	).FROM(
		jt.User.LEFT_JOIN(
			jt.EntityShare, jt.User.ID.EQ(jt.EntityShare.UserID),
		),
	).WHERE(
		jt.User.ID.EQ(jet.Int(uID)),
	)
	var u mt.User
	err := stmt.Query(s.db, &u)
	if err != nil {
		return u, err
	}

	return u, nil
}
