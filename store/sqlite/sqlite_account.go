package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	mt "ronche.se/moneytracker"

	jt "ronche.se/moneytracker/.gen/table"

	"github.com/go-jet/jet/v2/qrm"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetAccounts() ([]mt.Account, error) {

	Owner := jt.Entity.AS("owner")

	stmt := jet.SELECT(jt.Account.AllColumns,
		Owner.AllColumns,
	).FROM(jt.Account.INNER_JOIN(Owner, Owner.ID.EQ(jt.Account.OwnerID)))

	accounts := []mt.Account{}

	err := stmt.Query(s.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccountsByEntity(eID int64) ([]mt.Account, error) {

	stmt := jet.SELECT(jt.Account.AllColumns,
		jt.Entity.AllColumns,
	).FROM(jt.Account.INNER_JOIN(jt.Entity, jt.Entity.ID.EQ(jt.Account.OwnerID))).WHERE(jt.Entity.ID.EQ(jet.Int(int64(eID))))

	accounts := []mt.Account{}

	err := stmt.Query(s.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccount(aID int64) (*mt.Account, error) {

	stmt := jet.SELECT(jt.Account.AllColumns,
		jt.Entity.AllColumns,
	).FROM(jt.Account.INNER_JOIN(jt.Entity, jt.Entity.ID.EQ(jt.Account.OwnerID))).WHERE(jt.Account.ID.EQ(jet.Int(int64(aID))))

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

	var err error

	/* tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}() */

	err = insertAccount(s.db, a)
	if err != nil {
		return err
	}

	/* b := mt.Balance{
		AccountID: a.ID,
		ValueAt:   initialBalance,
	}

	err = insertBalances(tx, []mt.Balance{b})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	} */

	return nil

}

func insertAccount(db TXDB, a *mt.Account) error {
	stmt := jt.Account.INSERT(jt.Account.AllColumns).MODEL(a).RETURNING(jt.Account.AllColumns)

	err := stmt.Query(db, a)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) DeleteAccount(aID int64, onlyIfEmpty bool) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	if onlyIfEmpty {
		row := tx.QueryRow(`SELECT  count()
						FROM 'transaction' t
						WHERE from_id = :aID
						OR to_id = :aID`, sql.Named("aID", aID))
		var n int
		err = row.Scan(&n)
		if err != nil {
			return err
		}

		if n > 0 {
			return fmt.Errorf("impossible to delete account id=%d since there are %d transaction associated to it", aID, n)
		}

	}

	_, err = tx.Exec(`DELETE FROM account WHERE id=?`, aID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM balance WHERE account_id=?`, aID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
