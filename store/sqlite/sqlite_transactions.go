package sqlite

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
	"ronche.se/moneytracker/store/sqlite/queries"

	j "ronche.se/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetTransactionsByAccount(aID int) ([]mt.Transaction, error) {

	From := j.Account.AS("From")
	To := j.Account.AS("To")

	stmt := jet.SELECT(j.Transaction.AllColumns,
		j.Operation.AllColumns,
		From.AllColumns,
		To.AllColumns,
	).FROM(
		j.Transaction.INNER_JOIN(
			j.Operation,
			j.Operation.ID.EQ(j.Transaction.OperationID),
		).INNER_JOIN(From, From.ID.EQ(j.Transaction.FromID)).INNER_JOIN(To, To.ID.EQ(j.Transaction.ToID)),
	).WHERE(
		j.Transaction.FromID.EQ(jet.Int(int64(aID))).OR(j.Transaction.ToID.EQ(jet.Int(int64(aID)))),
	)

	transactions := []mt.Transaction{}

	err := stmt.Query(s.db, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *SQLiteStore) GetOperationsByEntity(eID int) ([]mt.Operation, error) {

	From := j.Account.AS("from")
	To := j.Account.AS("to")

	FromEntity := j.Entity.AS("from.entity")
	ToEntity := j.Account.AS("to.entity")

	stmt := jet.SELECT(
		j.Operation.AllColumns,
		j.Transaction.AllColumns,
		From.AllColumns,
		To.AllColumns,
		FromEntity.AllColumns,
		ToEntity.AllColumns,
	).FROM(
		j.Transaction.INNER_JOIN(
			j.Operation,
			j.Operation.ID.EQ(j.Transaction.OperationID),
		).INNER_JOIN(
			From,
			From.ID.EQ(j.Transaction.FromID),
		).INNER_JOIN(
			To,
			To.ID.EQ(j.Transaction.ToID),
		).INNER_JOIN(
			FromEntity,
			FromEntity.ID.EQ(From.OwnerID),
		).INNER_JOIN(
			ToEntity,
			ToEntity.ID.EQ(To.OwnerID),
		),
	).WHERE(
		FromEntity.ID.EQ(jet.Int(int64(eID))).OR(ToEntity.ID.EQ(jet.Int(int64(eID)))),
	)

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
		j.Operation.AllColumns,
		j.Transaction.AllColumns,
	).FROM(j.Operation.INNER_JOIN(j.Transaction, j.Transaction.OperationID.EQ(j.Operation.ID))).WHERE(j.Operation.ID.EQ(jet.Int(int64(opID))))

	op := &mt.Operation{}
	err := stmt.Query(s.db, op)
	if err != nil {
		return nil, err
	}

	return op, nil
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

	res, err := tx.Exec(queries.InsertOperationQuery, op.Timestamp, op.CreatedByID, op.Description)
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
	_, err := s.db.Exec("DELETE FROM operation WHERE id=?", opID)
	if err != nil {
		return err
	}

	// We should delete transactions with a trigger in the database

	return nil
}
