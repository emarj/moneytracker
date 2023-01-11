package sqlite

import (
	"fmt"

	mt "ronche.se/moneytracker"

	jt "ronche.se/moneytracker/.gen/table"
	"ronche.se/moneytracker/datetime"

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

func (s *SQLiteStore) GetOperationsByEntity(eID int64, limit int64) ([]mt.Operation, error) {

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
		jt.Balance.AllColumns,
	).FROM(
		jt.Operation.INNER_JOIN(
			jt.Transaction,
			jt.Transaction.OperationID.EQ(jt.Operation.ID),
		).INNER_JOIN(
			From,
			From.ID.EQ(jt.Transaction.FromID),
		).INNER_JOIN(
			To,
			To.ID.EQ(jt.Transaction.ToID),
		).INNER_JOIN(
			OwnerFrom,
			OwnerFrom.ID.EQ(From.OwnerID),
		).INNER_JOIN(
			OwnerTo,
			OwnerTo.ID.EQ(To.OwnerID),
		).LEFT_JOIN(
			jt.Balance,
			jt.Balance.OperationID.EQ(jt.Operation.ID),
		),
	).WHERE(
		To.OwnerID.EQ(jet.Int(eID)).OR(From.OwnerID.EQ(jet.Int(eID))),
	).ORDER_BY(jt.Operation.ModifiedOn.DESC()).LIMIT(limit)

	operations := []mt.Operation{}

	err := stmt.Query(s.db, &operations)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func (s *SQLiteStore) GetOperation(opID int64) (*mt.Operation, error) {

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		jt.Balance.AllColumns,
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
	).FROM(
		jt.Operation.INNER_JOIN(
			jt.Transaction,
			jt.Transaction.OperationID.EQ(jt.Operation.ID),
		).LEFT_JOIN(
			jt.Balance,
			jt.Balance.OperationID.EQ(jt.Operation.ID),
		).INNER_JOIN(
			From,
			From.ID.EQ(jt.Transaction.FromID),
		).INNER_JOIN(
			To,
			To.ID.EQ(jt.Transaction.ToID),
		).INNER_JOIN(
			OwnerFrom,
			OwnerFrom.ID.EQ(From.OwnerID),
		).INNER_JOIN(
			OwnerTo,
			OwnerTo.ID.EQ(To.OwnerID),
		),
	).WHERE(
		jt.Operation.ID.EQ(jet.Int(opID)),
	).ORDER_BY(jt.Transaction.Timestamp.DESC())

	fmt.Println(stmt.DebugSql())

	op := &mt.Operation{}
	err := stmt.Query(s.db, op)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (s *SQLiteStore) AddOperation(op *mt.Operation) error {

	if len(op.Transactions) == 0 && op.TypeID != mt.OpTypeBalanceAdjust {
		return fmt.Errorf("an operation must have at least one transaction or must be of type balance")
	}

	if len(op.Balances) == 0 && op.TypeID == mt.OpTypeBalanceAdjust {
		return fmt.Errorf("a balance adjust operation must have at least one balance")
	}

	// TODO: More checks

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	// save transaction list since we will overwrite it
	transactions := op.Transactions
	balances := op.Balances

	// we define another operation in order to update the pointer only if the transaction is successful
	var newOp mt.Operation

	stmt := jt.Operation.INSERT(jt.Operation.AllColumns.Except(jt.Operation.CreatedOn, jt.Operation.ModifiedOn)).MODEL(op).RETURNING(jt.Operation.AllColumns)

	err = stmt.Query(tx, &newOp)
	if err != nil {
		return err
	}

	if len(balances) > 0 && op.TypeID == mt.OpTypeBalanceAdjust {

		for i := range balances {

			balances[i].OperationID = newOp.ID

			b, err := getBalanceAt(tx, balances[i].AccountID.Int64, datetime.Now())
			if err != nil {
				return err
			}

			balances[i].Delta = b.Delta
		}

		err = insertBalances(tx, balances)
		if err != nil {
			return err
		}

		newOp.Balances = balances

	}

	if len(transactions) > 0 {
		for i := range transactions {
			transactions[i].Operation.ID = newOp.ID
			fmt.Println(transactions[i].From.ID, transactions[i].To.ID)
		}

		err = insertTransactions(tx, transactions)
		if err != nil {
			return err
		}

		newOp.Transactions = transactions
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// the insert is successful, update external operation
	op = &newOp

	return nil

}

func insertTransactions(db TXDB, txs []mt.Transaction) error {

	stmt := jt.Transaction.INSERT(jt.Transaction.AllColumns).MODELS(txs)
	//.RETURNING(jt.Transaction.AllColumns)

	_, err := stmt.Exec(db)
	if err != nil {
		return fmt.Errorf("insert transactions: %w", err)
	}

	return nil
}

func (s *SQLiteStore) DeleteOperation(opID int64) error {

	// We *could* delete transactions with a trigger in the database

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	var stmt jet.DeleteStatement
	stmt = jt.Transaction.DELETE().WHERE(jt.Transaction.OperationID.EQ(jet.Int(opID)))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}
	stmt = jt.Balance.DELETE().WHERE(jt.Balance.OperationID.EQ(jet.Int(opID)))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}
	stmt = jt.Operation.DELETE().WHERE(jt.Operation.ID.EQ(jet.Int(opID)))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
