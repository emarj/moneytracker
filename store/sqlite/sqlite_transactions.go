package sqlite

import mt "ronche.se/moneytracker"

func (s *SQLiteStore) GetTransactions() ([]mt.Transaction, error) {

	rows, err := s.db.Query("SELECT id,from_id,to_id,amount FROM transactions")
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var t mt.Transaction

	for rows.Next() {
		if err = rows.Scan(&t.ID, &t.From, &t.To, &t.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (s *SQLiteStore) GetTransactionsByAccount(aID int) ([]mt.Transaction, error) {

	rows, err := s.db.Query("SELECT id,from_id,to_id,amount FROM transactions WHERE from_id = ? OR to_id = ?", aID, aID)
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var t mt.Transaction

	for rows.Next() {
		if err = rows.Scan(&t.ID, &t.From, &t.To, &t.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (s *SQLiteStore) GetTransaction(id int) (mt.Transaction, error) {

	row := s.db.QueryRow("SELECT id,from_id,to_id,amount FROM transactions WHERE id = ?", id)

	var t mt.Transaction

	if err := row.Scan(&t.ID, &t.From, &t.To, &t.Amount); err != nil {
		return t, err
	}

	return t, nil
}

func (s *SQLiteStore) AddTransaction(t mt.Transaction) (int, error) {

	res, err := s.db.Exec("INSERT INTO transactions (from_id,to_id,amount) VALUES(?,?,?)", t.From, t.To, t.Amount)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}

/*func (s *SQLiteStore) AddTransaction(t mt.Transaction) (int, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}
	res, err := tx.Exec("INSERT INTO transactions (from_id,to_id,amount) VALUES(?,?,?)", t.From, t.To, t.Amount)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE  id = ?", t.Amount, t.To)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE  id = ?", t.Amount, t.From)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}*/

func (s *SQLiteStore) DeleteTransaction(id int) error {
	_, err := s.db.Exec("DELETE FROM transactions WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
