package sqlite

import (
	"database/sql"
	"fmt"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

const GetTransactionsByAccountQuery string = `SELECT  t.id,
													t.from_id,
													t.to_id,
													t.amount,
													t.operation_id,
													op.id,
													op.timestamp,
													op.created_by_id,
													op.description
											FROM transactions t
											INNER JOIN operations op
											ON t.operation_id = op.id
											WHERE from_id = :aID
											OR to_id = :aID
											ORDER BY op.timestamp DESC
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
												FROM transactions t
													INNER JOIN operations op ON t.operation_id = op.id
													INNER JOIN accounts AS fa ON t.from_id = fa.id
													INNER JOIN accounts AS ta ON t.to_id = ta.id
													INNER JOIN entities AS fe ON fa.owner_id = fe.id
													INNER JOIN entities AS te ON ta.owner_id = te.id
												WHERE fa.owner_id = :eID
													OR ta.owner_id = :eID
												ORDER BY op.timestamp DESC,op.id,t.id
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
											FROM operations op
											INNER JOIN transactions t
											ON t.operation_id = op.id
											WHERE op.id = :oID;`

const InsertTransactionQuery string = `INSERT INTO  transactions (
													from_id,
													to_id,
													amount,
													operation_id)
											VALUES (?,?,?,?);`

const InsertOperationQuery string = `INSERT INTO  operations (
													timestamp,
													created_by_id,
													description,
													category_id)
												VALUES (?,?,?,?);`

func (s *SQLiteStore) GetTransactionsByAccount(aID int, limit int) ([]mt.Transaction, error) {

	rows, err := s.db.Query(GetTransactionsByAccountQuery, sql.Named("aID", aID), limit)
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var t mt.Transaction

	for rows.Next() {
		t.Operation = mt.Operation{}

		if err = rows.Scan(&t.ID, &t.From.ID, &t.To.ID, &t.Amount, &t.Operation.ID, &t.Operation.ID, &t.Operation.Timestamp, &t.Operation.CreatedByID, &t.Operation.Description); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (s *SQLiteStore) GetOperationsByEntity(eID int, limit int) ([]mt.Operation, error) {

	rows, err := s.db.Query(GetOperationByEntityQuery, sql.Named("eID", eID), limit)
	if err != nil {
		return nil, err
	}

	operations := orderedmap.New[int64, mt.Operation]()

	op := mt.Operation{}

	for rows.Next() {
		t := mt.Transaction{}

		if err = rows.Scan(
			&t.ID, &t.From.ID, &t.To.ID, &t.Amount, &t.Operation.ID,
			&op.ID, &op.Timestamp, &op.CreatedByID, &op.Description, &op.CategoryID,
			&t.From.Name, &t.From.DisplayName, &t.To.Name, &t.To.DisplayName,
			&t.From.Owner.ID, &t.From.Owner.Name, &t.To.Owner.ID, &t.To.Owner.Name,
		); err != nil {
			return nil, err
		}

		//op.ID can't be null
		op2, ok := operations.Get(op.ID.Int64)
		if !ok {
			op2 = op
		}

		op2.Transactions = append(op2.Transactions, t)
		operations.Set(op.ID.Int64, op2)
	}

	list := make([]mt.Operation, operations.Len())

	i := 0
	for pair := operations.Oldest(); pair != nil; pair = pair.Next() {
		list[i] = pair.Value
		i++
	}

	return list, nil
}

func (s *SQLiteStore) GetOperation(opID int) (*mt.Operation, error) {

	rows, err := s.db.Query(GetOperationQuery, sql.Named("opID", opID))
	if err != nil {
		return nil, err
	}
	var op mt.Operation
	op.Transactions = []mt.Transaction{}

	for rows.Next() {
		t := mt.Transaction{}
		if err = rows.Scan(
			&op.ID, &op.Timestamp, &op.CreatedByID, &op.Description, &op.CategoryID,
			&t.ID, &t.From.ID, &t.To.ID, &t.Amount,
		); err != nil {
			return nil, err
		}

		op.Transactions = append(op.Transactions, t)
	}

	return &op, nil
}

func (s *SQLiteStore) AddOperation(op mt.Operation) (null.Int, error) {

	id := null.Int{}
	if op.Transactions == nil || len(op.Transactions) == 0 {
		return id, fmt.Errorf("an operation must have at least one transaction")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return id, err
	}
	defer func() {
		tx.Rollback()
	}()

	res, err := tx.Exec(InsertOperationQuery, op.Timestamp, op.CreatedByID, op.Description, op.CategoryID)
	if err != nil {
		return id, err
	}

	id.Int64, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	q, err := tx.Prepare(InsertTransactionQuery)
	if err != nil {
		return id, err
	}

	for _, t := range op.Transactions {
		_, err = q.Exec(t.From.ID, t.To.ID, t.Amount, id.Int64)
		if err != nil {
			return id, err
		}

	}

	err = tx.Commit()
	if err != nil {
		return id, err
	}

	id.Valid = true
	return id, nil

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

	res, err := tx.Exec("DELETE FROM operations WHERE id=?", opID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("error: no operation with id: %d", opID)
	}

	_, err = tx.Exec("DELETE FROM transactions WHERE operation_id=?", opID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
