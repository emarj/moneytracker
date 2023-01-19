package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

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
