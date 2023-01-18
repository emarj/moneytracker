package sqlite

import (
	"errors"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"

	"github.com/go-jet/jet/v2/qrm"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetAccounts() ([]mt.Account, error) {

	Owner := jt.Entity.AS("owner")

	stmt := jet.SELECT(
		jt.Account.AllColumns,
		Owner.AllColumns,
	).FROM(jt.Account.
		INNER_JOIN(Owner, Owner.ID.EQ(jt.Account.OwnerID)),
	)

	accounts := []mt.Account{}

	err := stmt.Query(s.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccountsByEntity(eID int64) ([]mt.Account, error) {
	Owner := jt.Entity.AS("owner")

	stmt := jet.SELECT(
		jt.Account.AllColumns,
		Owner.AllColumns,
	).FROM(
		jt.Account.
			INNER_JOIN(Owner,
				Owner.ID.EQ(jt.Account.OwnerID),
			)).
		WHERE(Owner.ID.EQ(jet.Int(int64(eID))))

	accounts := []mt.Account{}

	err := stmt.Query(s.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SQLiteStore) GetUserAccounts(uID int64) ([]mt.Account, error) {
	Owner := jt.Entity.AS("owner")

	stmt := jet.SELECT(
		jt.Account.AllColumns,
		Owner.AllColumns,
		jt.EntityShare.AllColumns,
	).FROM(
		jt.Account.
			INNER_JOIN(Owner,
				Owner.ID.EQ(jt.Account.OwnerID),
			).
			INNER_JOIN(jt.EntityShare,
				jt.EntityShare.EntityID.EQ(Owner.ID),
			)).
		WHERE(Owner.ID.IN(
			jt.EntityShare.SELECT(jt.EntityShare.EntityID).
				WHERE(jt.EntityShare.UserID.EQ(jet.Int(uID))),
		)).ORDER_BY(jt.EntityShare.Quota.DESC())

	accounts := []mt.Account{}

	err := stmt.Query(s.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccount(aID int64) (*mt.Account, error) {

	Owner := jt.Entity.AS("owner")

	stmt := jet.SELECT(
		jt.Account.AllColumns,
		Owner.AllColumns,
	).FROM(jt.Account.
		INNER_JOIN(Owner, Owner.ID.EQ(jt.Account.OwnerID))).
		WHERE(jt.Account.ID.EQ(jet.Int(int64(aID))))

	var a mt.Account
	err := stmt.Query(s.db, &a)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, mt.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (s *SQLiteStore) AddAccount(a *mt.Account) error {
	err := insertAccount(s.db, a)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) UpdateAccount(a *mt.Account) error {
	err := updateAccount(s.db, a)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) DeleteAccount(aID int64) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	err = deleteAccount(tx, aID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
