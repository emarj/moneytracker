package sqlite

import (
	"fmt"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
	"ronche.se/moneytracker/store/sqlite/queries"
)

func (s *SQLiteStore) GetTransactionsByAccount(aID int) ([]mt.Transaction, error) {

	rows, err := s.db.Query(queries.GetTransactionsByAccountQuery, aID, aID)
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

func (s *SQLiteStore) GetOperationsByEntity(eID int) ([]mt.Operation, error) {

	rows, err := s.db.Query(queries.GetOperationByEntityQuery, eID, eID)
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

	rows, err := s.db.Query(queries.GetOperationQuery, opID)
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

	res, err := tx.Exec(queries.InsertOperationQuery, op.Timestamp, op.CreatedByID, op.Description, op.CategoryID)
	if err != nil {
		return id, err
	}

	id.Int64, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	q, err := tx.Prepare(queries.InsertTransactionQuery)
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
	_, err := s.db.Exec("DELETE FROM operations WHERE id=?", opID)
	if err != nil {
		return err
	}

	// We should delete transactions with a trigger in the database

	return nil
}
