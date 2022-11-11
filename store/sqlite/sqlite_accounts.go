package sqlite

import mt "ronche.se/moneytracker"

func (s *SQLiteStore) GetAllAccounts() ([]mt.Account, error) {

	rows, err := s.db.Query("SELECT id,name,balance FROM accounts")
	if err != nil {
		return nil, err
	}

	accounts := []mt.Account{}
	var a mt.Account

	for rows.Next() {
		if err = rows.Scan(&a.ID, &a.Name, &a.Balance); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *SQLiteStore) AddAccount(a mt.Account) (int, error) {

	res, err := s.db.Exec("INSERT INTO accounts (name) VALUES(?)", a.Name)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}

func (s *SQLiteStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
