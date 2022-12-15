package sqlite

import (
	"database/sql"
	"fmt"

	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"

	jt "ronche.se/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

const GetTransactionsByAccountQuery string = `SELECT  t.id,
													t.timestamp,	
													t.from_id,
													t.to_id,
													t.amount,
													t.operation_id,
													op.id,
													op.created_by_id,
													op.description
											FROM 'transaction' t
											INNER JOIN operation op
											ON t.operation_id = op.id
											WHERE from_id = :aID
											OR to_id = :aID
											ORDER BY t.timestamp DESC
											LIMIT ?;`

const GetOperationByEntityQuery string = `SELECT  t.*,
												op.*,
												fa.name AS from_name,
												fa.display_name AS from_display_name,
												ta.name AS to_name,
												ta.display_name AS to_display_name,
												fe.id,
												fe.name,
												te.id,
												te.name
												FROM 'transaction' t
													INNER JOIN operation op ON t.operation_id = op.id
													INNER JOIN account AS fa ON t.from_id = fa.id
													INNER JOIN account AS ta ON t.to_id = ta.id
													INNER JOIN entity AS fe ON fa.owner_id = fe.id
													INNER JOIN entity AS te ON ta.owner_id = te.id
												WHERE fa.owner_id = :eID
													OR ta.owner_id = :eID
												ORDER BY t.timestamp DESC,op.id,t.id
												LIMIT ?;`

const GetOperationQuery string = `SELECT  			op.id,
													op.timestamp,
													op.created_by_id,
													op.description,
													op.category_id,
													t.id,
													t.from_id,
													t.to_id,
													t.amount
											FROM operation op
											INNER JOIN 'transaction' t
											ON t.operation_id = op.id
											WHERE op.id = :oID;`

const InsertTransactionQuery string = `INSERT INTO  'transaction' (
													timestamp,
													from_id,
													to_id,
													amount,
													operation_id)
											VALUES (?,?,?,?,?);`

const InsertOperationQuery string = `INSERT INTO  operation (
													created_by_id,
													description,
													category_id)
												VALUES (?,?,?);`

func (s *SQLiteStore) GetTransactionsByAccount(aID int, limit int) ([]mt.Transaction, error) {

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")

	stmt := jet.SELECT(jt.Transaction.AllColumns,
		jt.Operation.AllColumns,
		From.AllColumns,
		To.AllColumns,
	).FROM(
		jt.Transaction.INNER_JOIN(
			jt.Operation,
			jt.Operation.ID.EQ(jt.Transaction.OperationID),
		).INNER_JOIN(From, From.ID.EQ(jt.Transaction.FromID)).INNER_JOIN(To, To.ID.EQ(jt.Transaction.ToID)),
	).WHERE(
		jt.Transaction.FromID.EQ(jet.Int(int64(aID))).OR(jt.Transaction.ToID.EQ(jet.Int(int64(aID)))),
	)

	transactions := []mt.Transaction{}

	err := stmt.Query(s.db, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *SQLiteStore) GetOperationsByEntity(eID int, limit int) ([]mt.Operation, error) {

	From := jt.Account.AS("from")
	To := jt.Account.AS("to")

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
		//jt.Balance.AllColumns,
		From.AllColumns,
		To.AllColumns,
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
		),
	).WHERE(
		To.OwnerID.EQ(jet.Int(int64(eID))).OR(From.OwnerID.EQ(jet.Int(int64(eID)))),
	).ORDER_BY(jt.Operation.ModifiedOn.DESC()).LIMIT(int64(limit))

	fmt.Println(stmt.DebugSql())

	operations := []mt.Operation{}

	err := stmt.Query(s.db, &operations)
	if err != nil {
		return nil, err
	}
	return operations, nil
}

func (s *SQLiteStore) GetOperation(opID int) (*mt.Operation, error) {

	stmt := jet.SELECT(
		jt.Operation.AllColumns,
		jt.Transaction.AllColumns,
	).FROM(jt.Operation.INNER_JOIN(jt.Transaction, jt.Transaction.OperationID.EQ(jt.Operation.ID))).WHERE(jt.Operation.ID.EQ(jet.Int(int64(opID))))

	op := &mt.Operation{}
	err := stmt.Query(s.db, op)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (s *SQLiteStore) AddOperation(op mt.Operation) (null.Int, error) {

	if op.Transactions == nil || len(op.Transactions) == 0 {
		return op.ID, fmt.Errorf("an operation must have at least one transaction")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return op.ID, err
	}
	defer func() {
		tx.Rollback()
	}()

	res, err := tx.Exec(InsertOperationQuery, op.CreatedByID, op.Description, op.CategoryID)
	if err != nil {
		return op.ID, err
	}

	op.ID.Int64, err = res.LastInsertId()
	if err != nil {
		return op.ID, err
	}

	err = addTransactions(tx, op.Transactions, op.ID.Int64)
	if err != nil {
		return op.ID, err
	}

	err = tx.Commit()
	if err != nil {
		return op.ID, err
	}

	op.ID.Valid = true
	return op.ID, nil

}

type DB interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
}

func addTransaction(db DB, t *mt.Transaction) error {
	res, err := db.Exec(InsertTransactionQuery, t.Timestamp, t.From.ID, t.To.ID, t.Amount, t.Operation.ID)
	if err != nil {
		return err
	}

	t.ID.Int64, err = res.LastInsertId()
	if err != nil {
		return err
	}

	t.ID.Valid = true
	return nil
}

func addTransactions(db DB, txs []mt.Transaction, opID int64) error {
	q, err := db.Prepare(InsertTransactionQuery)
	if err != nil {
		return err
	}

	var res sql.Result
	for _, t := range txs {
		res, err = q.Exec(t.Timestamp, t.From.ID, t.To.ID, t.Amount, opID)
		if err != nil {
			return fmt.Errorf("sdsdd %w", err)
		}

		t.ID.Int64, err = res.LastInsertId()
		if err != nil {
			return err
		}
		t.ID.Valid = true
	}

	return nil
}

func (s *SQLiteStore) DeleteOperation(opID int) error {

	// We could delete transactions with a trigger in the database

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	res, err := tx.Exec("DELETE FROM operation WHERE id=?", opID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("error: no operation with id: %d", opID)
	}

	_, err = tx.Exec("DELETE FROM 'transaction' WHERE operation_id=?", opID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
