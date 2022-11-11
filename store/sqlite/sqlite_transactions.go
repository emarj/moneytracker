package sqlite

import mt "ronche.se/moneytracker"

func (s *SQLiteStore) GetAllTransactions() ([]mt.Transaction, error) {

	rows, err := s.db.Query("SELECT id,from,to,amount FROM transactions")
	if err != nil {
		return nil, err
	}

	transactions := []mt.Transaction{}
	var tx mt.Transaction

	for rows.Next() {
		if err = rows.Scan(&tx.ID, &tx.From, &tx.To, &tx.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (s *SQLiteStore) AddTransaction(tx mt.Transaction) (int, error) {

	res, err := s.db.Exec("INSERT INTO transactions (from,to,amount) VALUES(?,?,?)", tx.From, tx.To, tx.Amount)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}

func (s *SQLiteStore) DeleteTransaction(id int) error {
	_, err := s.db.Exec("DELETE FROM transactions WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
