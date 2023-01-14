package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/emarj/moneytracker/datetime"

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
	BalanceAcc := jt.Account.AS("balance_account")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")
	OwnerBalanceAcc := jt.Entity.AS("balance_account.owner")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
		jt.Balance.AllColumns,
	).FROM(
		(jt.Operation.SELECT(jet.STAR).ORDER_BY(jt.Operation.ModifiedOn.DESC()).LIMIT(limit)).AsTable("operation").LEFT_JOIN(
			jt.Transaction,
			jt.Transaction.OperationID.EQ(jt.Operation.ID),
		).LEFT_JOIN(
			From,
			From.ID.EQ(jt.Transaction.FromID),
		).LEFT_JOIN(
			To,
			To.ID.EQ(jt.Transaction.ToID),
		).LEFT_JOIN(
			OwnerFrom,
			OwnerFrom.ID.EQ(From.OwnerID),
		).LEFT_JOIN(
			OwnerTo,
			OwnerTo.ID.EQ(To.OwnerID),
		).LEFT_JOIN(
			jt.Balance,
			jt.Balance.OperationID.EQ(jt.Operation.ID),
		).LEFT_JOIN(
			BalanceAcc,
			BalanceAcc.ID.EQ(jt.Balance.AccountID),
		).LEFT_JOIN(
			OwnerBalanceAcc,
			OwnerBalanceAcc.ID.EQ(BalanceAcc.OwnerID),
		),
	).WHERE(
		(To.OwnerID.EQ(jet.Int(eID)).OR(From.OwnerID.EQ(jet.Int(eID)))).OR(OwnerBalanceAcc.ID.EQ(jet.Int(eID))),
	)

	//println(stmt.DebugSql())

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
	BalanceAcc := jt.Account.AS("balance_account")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")
	OwnerBalanceAcc := jt.Entity.AS("balance_account.owner")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		jt.Balance.AllColumns,
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
	).FROM(
		jt.Operation.LEFT_JOIN(
			jt.Transaction,
			jt.Transaction.OperationID.EQ(jt.Operation.ID),
		).LEFT_JOIN(
			From,
			From.ID.EQ(jt.Transaction.FromID),
		).LEFT_JOIN(
			To,
			To.ID.EQ(jt.Transaction.ToID),
		).LEFT_JOIN(
			OwnerFrom,
			OwnerFrom.ID.EQ(From.OwnerID),
		).LEFT_JOIN(
			OwnerTo,
			OwnerTo.ID.EQ(To.OwnerID),
		).LEFT_JOIN(
			jt.Balance,
			jt.Balance.OperationID.EQ(jt.Operation.ID),
		).LEFT_JOIN(
			BalanceAcc,
			BalanceAcc.ID.EQ(jt.Balance.AccountID),
		).LEFT_JOIN(
			OwnerBalanceAcc,
			OwnerBalanceAcc.ID.EQ(BalanceAcc.OwnerID),
		),
	).WHERE(
		jt.Operation.ID.EQ(jet.Int(opID)),
	).ORDER_BY(jt.Transaction.Timestamp.DESC())

	op := &mt.Operation{}
	err := stmt.Query(s.db, op)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func insertOperation(txdb TXDB, op *mt.Operation) error {
	now := datetime.Now()

	op.CreatedOn = now
	op.ModifiedOn = now

	stmt := jt.Operation.INSERT(jt.Operation.AllColumns).
		MODEL(op).
		RETURNING(jt.Operation.ID)

	// We do not use Jet here since it only works for structs
	q, args := stmt.Sql()
	err := txdb.QueryRow(q, args...).Scan(&op.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) AddOperation(op *mt.Operation) error {

	if len(op.Transactions) == 0 && op.TypeID != mt.OpTypeBalanceAdjust {
		return fmt.Errorf("an operation must have at least one transaction or must be of type balance")
	}

	if op.TypeID == mt.OpTypeBalanceAdjust {
		if len(op.Balances) == 0 {
			return fmt.Errorf("a balance adjust operation must have at least one balance")
		}
		if len(op.Transactions) > 0 {
			return fmt.Errorf("a balance adjust operation must not have transactions")
		}
	}

	// TODO: More checks

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	//We need to do this in order to update the external operation only in case of success
	newOp := *op

	err = insertOperation(tx, &newOp)
	if err != nil {
		return err
	}

	for i := range newOp.Balances {
		newOp.Balances[i].OperationID = newOp.ID

		err = insertBalance(tx, &newOp.Balances[i])
		if err != nil {
			return err
		}

		newOp.Balances[i].Operation = &newOp
	}

	if len(newOp.Transactions) > 0 {
		for i := range newOp.Transactions {
			newOp.Transactions[i].OperationID = newOp.ID.Int64

			err = insertTransaction(tx, &newOp.Transactions[i])
			if err != nil {
				return err
			}

			newOp.Transactions[i].Operation = &newOp
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// the insert is successful, update external operation
	op = &newOp

	return nil

}

func insertTransaction(txdb TXDB, tx *mt.Transaction) error {

	if tx.FromID == tx.ToID {
		return fmt.Errorf("a transaction cannot be from and to the same account")
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

func (s *SQLiteStore) DeleteOperation(opID int64) error {

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
