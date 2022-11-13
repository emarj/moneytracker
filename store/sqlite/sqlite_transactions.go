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
