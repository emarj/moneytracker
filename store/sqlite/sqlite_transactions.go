package sqlite

import mt "ronche.se/moneytracker"

func (s *SQLiteStore) GetTransactions() ([]mt.Transaction, error) {

	rows, err := s.db.Query("SELECT id,timestamp,from_id,to_id,amount,operation_id FROM transactions")
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var t mt.Transaction

	for rows.Next() {
		if err = rows.Scan(&t.ID, &t.Timestamp, &t.FromID, &t.ToID, &t.Amount, &t.OperationID); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (s *SQLiteStore) GetTransactionsByAccount(aID int) ([]mt.Transaction, error) {

	rows, err := s.db.Query("SELECT id,timestamp,from_id,to_id,amount,operation_id FROM transactions WHERE from_id = ? OR to_id = ?", aID, aID)
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var t mt.Transaction

	for rows.Next() {
		if err = rows.Scan(&t.ID, &t.Timestamp, &t.FromID, &t.ToID, &t.Amount, &t.OperationID); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (s *SQLiteStore) GetTransactionsByEntity(eID int) ([]mt.Transaction, error) {

	rows, err := s.db.Query(`
	SELECT t.*,
    fa.display_name AS from_name,
    ta.display_name AS to_name,
    fa.entity_id AS from_entity_id,
    ta.entity_id AS to_entity_id
FROM transactions AS t
    INNER JOIN accounts AS fa ON t.from_id = fa.id
    INNER JOIN accounts AS ta ON t.to_id = ta.id
WHERE from_entity_id = ?
    OR to_entity_id = ?`,
		eID, eID)
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var t mt.Transaction

	for rows.Next() {
		if err = rows.Scan(
			&t.ID, &t.Timestamp, &t.FromID, &t.ToID, &t.Amount, &t.OperationID,
			&t.From.DisplayName, &t.To.DisplayName, &t.From.EntityID, &t.To.EntityID,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (s *SQLiteStore) GetTransaction(aID int) (*mt.Transaction, error) {

	row := s.db.QueryRow("SELECT id,from_id,to_id,amount FROM transactions WHERE id = ?", aID)

	var t mt.Transaction

	if err := row.Scan(&t.ID, &t.Timestamp, &t.FromID, &t.ToID, &t.Amount, &t.OperationID); err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *SQLiteStore) AddTransaction(t mt.Transaction) (int, error) {

	res, err := s.db.Exec("INSERT INTO transactions (timestamp,from_id,to_id,amount) VALUES(?,?,?,?)", t.Timestamp, t.FromID, t.ToID, t.Amount)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}

func (s *SQLiteStore) DeleteTransaction(tID int) error {
	_, err := s.db.Exec("DELETE FROM transactions WHERE id=?", tID)
	if err != nil {
		return err
	}

	// We should reset balances with a trigger in the database

	return nil
}
