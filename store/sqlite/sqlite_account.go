package sqlite

import (
	"errors"
	"fmt"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"

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
	err := insertAccount(s.db, a)
	if err != nil {
		return err
	}

	return nil

}

func insertAccount(txdb TXDB, a *mt.Account) error {
	stmt := jt.Account.INSERT(jt.Account.AllColumns).
		MODEL(a).
		RETURNING(jt.Account.AllColumns)

	err := stmt.Query(txdb, a)
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

func updateAccount(txdb TXDB, a *mt.Account) error {
	stmt := jt.Account.UPDATE(jt.Account.AllColumns).
		WHERE(jt.Account.ID.EQ(jet.Int(a.ID.Int64))).
		MODEL(a) //.RETURNING(jt.Account.AllColumns)

	//println(stmt.DebugSql())
	_, err := stmt.Exec(txdb)
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

func deleteAccount(txdb TXDB, aID int64) error {

	selectStmt := jt.Transaction.SELECT(jet.COUNT(jet.STAR)).
		WHERE(
			jt.Transaction.FromID.EQ(jet.Int(aID)).
				OR(jt.Transaction.ToID.EQ(jet.Int(aID))),
		)

	q, args := selectStmt.Sql()
	var n int
	err := txdb.QueryRow(q, args...).Scan(&n)
	if err != nil {
		return err
	}

	if n > 0 {
		return fmt.Errorf("impossible to delete account id=%d: there are %d transaction associated with it", aID, n)
	}

	//println(stmt.DebugSql())

	_, err = jt.Account.DELETE().
		WHERE(jt.Account.ID.EQ(jet.Int(aID))).
		Exec(txdb)
	if err != nil {
		return err
	}

	_, err = jt.Balance.DELETE().
		WHERE(jt.Balance.AccountID.EQ(jet.Int(aID))).
		Exec(txdb)
	if err != nil {
		return err
	}

	return nil
}
