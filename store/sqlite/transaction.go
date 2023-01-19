package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetTransactionsByAccount(aID int64, limit int64) ([]mt.Transaction, error) {

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		From.AllColumns,
		To.AllColumns,
	).FROM(
		jt.Transaction.INNER_JOIN(
			jt.Operation,
			jt.Operation.ID.EQ(jt.Transaction.OperationID),
		).INNER_JOIN(From, From.ID.EQ(jt.Transaction.FromID)).INNER_JOIN(To, To.ID.EQ(jt.Transaction.ToID)),
	).WHERE(
		jt.Transaction.FromID.EQ(jet.Int(aID)).OR(jt.Transaction.ToID.EQ(jet.Int(aID))),
	).ORDER_BY(jt.Transaction.Timestamp.DESC()).LIMIT(limit)

	transactions := []mt.Transaction{}

	err := stmt.Query(s.db, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func insertTransaction(txdb TXDB, tx *mt.Transaction) error {

	if tx.FromID == tx.ToID {
		return fmt.Errorf("a transaction cannot be from and to the same account: %d to %d", tx.FromID, tx.ToID)
	}

	stmt := jt.Transaction.INSERT(jt.Transaction.AllColumns).MODEL(tx).RETURNING(jt.Transaction.AllColumns)

	err := stmt.Query(txdb, tx)
	if err != nil {
		return fmt.Errorf("insert transactions: %w", err)
	}

	//TODO: Here we should fix future balances and deltas until the first user inserted one
	/* err = updateBalances(txdb, tx.Timestamp, tx.From.ID.Int64, tx.To.ID.Int64)
	if err != nil {
		return fmt.Errorf("insert transactions: %w", err)
	} */

	return nil
}

func updateTransaction(txdb TXDB, tx *mt.Transaction) error {

	if tx.FromID == tx.ToID {
		return fmt.Errorf("a transaction cannot be from and to the same account: %d to %d", tx.FromID, tx.ToID)
	}

	stmt := jt.Transaction.UPDATE(jt.Transaction.AllColumns).
		WHERE(jt.Transaction.ID.EQ(jet.Int(tx.ID.Int64))).
		MODEL(tx) //.RETURNING(jt.Transaction.AllColumns)

	err := stmt.Query(txdb, tx)
	if err != nil {
		return err
	}

	//TODO: Here we should fix future balances and deltas until the first user inserted one
	/* err = updateBalances(txdb, tx.Timestamp, tx.From.ID.Int64, tx.To.ID.Int64)
	if err != nil {
		return fmt.Errorf("insert transactions: %w", err)
	} */

	return nil
}
