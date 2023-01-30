package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"
	"github.com/emarj/moneytracker/timestamp"

	jt "github.com/emarj/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetOperationsOfUser(uID int64, limit int64) ([]mt.Operation, error) {

	//userEntities := jt.EntityShare.SELECT(jt.EntityShare.EntityID).WHERE(jt.EntityShare.UserID.EQ(jet.Int(uID)))

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")
	Balance := jt.Account.AS("account")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")
	OwnerBalance := jt.Entity.AS("account.owner")

	UserFrom := jt.EntityShare.AS("from.user")
	UserTo := jt.EntityShare.AS("to.user")
	UserBalance := jt.EntityShare.AS("account.user")

	operationTable := jt.Operation.SELECT(jet.STAR).
		WHERE(jet.EXISTS(
			jet.SELECT(jet.STAR).
				FROM(jt.Transaction.LEFT_JOIN(
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
					Balance,
					Balance.ID.EQ(jt.Balance.AccountID),
				).LEFT_JOIN(
					OwnerBalance,
					OwnerBalance.ID.EQ(Balance.OwnerID),
				).LEFT_JOIN(
					UserFrom,
					UserFrom.EntityID.EQ(From.OwnerID),
				).LEFT_JOIN(
					UserTo,
					UserTo.EntityID.EQ(To.OwnerID),
				).LEFT_JOIN(
					UserBalance,
					UserBalance.EntityID.EQ(Balance.OwnerID),
				),
				).WHERE(
				jt.Transaction.OperationID.EQ(jt.Operation.ID).
					AND(
						(UserFrom.UserID.EQ(jet.Int(uID)).
							OR(UserTo.UserID.EQ(jet.Int(uID)))).
							OR(UserBalance.UserID.EQ(jet.Int(uID))),
					))),
		).
		ORDER_BY(jt.Operation.ModifiedOn.DESC()).
		LIMIT(limit).
		AsTable("operation")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		(jet.CASE().
			WHEN(UserTo.UserID.EQ(UserFrom.UserID)).
			THEN(jet.Int(0)).
			WHEN(UserFrom.UserID.EQ(jet.Int(uID))).
			THEN(jet.Int(-1)).
			WHEN(UserTo.UserID.EQ(jet.Int(uID))).
			THEN(jet.Int(1)).
			ELSE(jet.NULL)).
			AS("transaction.sign"),
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
		jt.Balance.AllColumns,
		Balance.AllColumns,
		OwnerBalance.AllColumns,
		/* UserFrom.UserID,
		UserTo.UserID,
		UserBalance.UserID, */
	).FROM(operationTable.LEFT_JOIN(
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
		Balance,
		Balance.ID.EQ(jt.Balance.AccountID),
	).LEFT_JOIN(
		OwnerBalance,
		OwnerBalance.ID.EQ(Balance.OwnerID),
	).LEFT_JOIN(
		UserFrom,
		UserFrom.EntityID.EQ(From.OwnerID),
	).LEFT_JOIN(
		UserTo,
		UserTo.EntityID.EQ(To.OwnerID),
	).LEFT_JOIN(
		UserBalance,
		UserBalance.EntityID.EQ(Balance.OwnerID),
	),
	).ORDER_BY(
		jt.Operation.ModifiedOn.DESC(),
		jt.Transaction.Timestamp.DESC(),
		jt.Balance.Timestamp.DESC(),
	)

	//fmt.Println(stmt.DebugSql())

	operations := []mt.Operation{}

	err := stmt.Query(s.db, &operations)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func (s *SQLiteStore) GetOperationsByEntity(eID int64, limit int64) ([]mt.Operation, error) {

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")
	Balance := jt.Account.AS("account")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")
	OwnerBalance := jt.Entity.AS("account.owner")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
		jt.Balance.AllColumns,
		Balance.AllColumns,
		OwnerBalance.AllColumns,
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
			Balance,
			Balance.ID.EQ(jt.Balance.AccountID),
		).LEFT_JOIN(
			OwnerBalance,
			OwnerBalance.ID.EQ(Balance.OwnerID),
		),
	).WHERE(
		(To.OwnerID.EQ(jet.Int(eID)).OR(From.OwnerID.EQ(jet.Int(eID)))).OR(OwnerBalance.ID.EQ(jet.Int(eID))),
	).ORDER_BY(
		jt.Operation.ModifiedOn.DESC(),
		jt.Transaction.Timestamp.DESC(),
		jt.Balance.Timestamp.DESC(),
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
	Balance := jt.Account.AS("account")

	OwnerFrom := jt.Entity.AS("from.owner")
	OwnerTo := jt.Entity.AS("to.owner")
	OwnerBalance := jt.Entity.AS("account.owner")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		From.AllColumns,
		To.AllColumns,
		OwnerFrom.AllColumns,
		OwnerTo.AllColumns,
		jt.Balance.AllColumns,
		Balance.AllColumns,
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
			Balance,
			Balance.ID.EQ(jt.Balance.AccountID),
		).LEFT_JOIN(
			OwnerBalance,
			OwnerBalance.ID.EQ(Balance.OwnerID),
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
	now := timestamp.Now()

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

		// Comment this to avoid cycles in JSON marshaling
		//newOp.Balances[i].Operation = &newOp
	}

	if len(newOp.Transactions) > 0 {
		for i := range newOp.Transactions {
			newOp.Transactions[i].OperationID = newOp.ID.Int64

			err = insertTransaction(tx, &newOp.Transactions[i])
			if err != nil {
				return err
			}
			// Comment this to avoid cycles in JSON marshaling
			//newOp.Transactions[i].Operation = &newOp
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// the insert is successful, update external operation
	*op = newOp

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
