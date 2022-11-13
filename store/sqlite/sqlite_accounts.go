package sqlite

import (
	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) GetAccounts() ([]mt.Account, error) {

	rows, err := s.db.Query(`SELECT  id,name FROM accounts`)
	if err != nil {
		return nil, err
	}

	accounts := []mt.Account{}
	var a mt.Account

	for rows.Next() {
		if err = rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccountsByEntity(eID int) ([]mt.Account, error) {

	rows, err := s.db.Query(`SELECT  id,name,display_name FROM accounts WHERE entity_id = ?`, eID)
	if err != nil {
		return nil, err
	}

	accounts := []mt.Account{}
	var a mt.Account

	for rows.Next() {
		if err = rows.Scan(&a.ID, &a.Name, &a.DisplayName); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccount(aID int) (*mt.Account, error) {

	var a mt.Account

	err := s.db.QueryRow(`SELECT id,name FROM accounts WHERE id = ?`, aID).Scan(
		&a.ID,
		&a.Name,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *SQLiteStore) AddAccount(a mt.Account) (int, error) {

	res, err := s.db.Exec(`INSERT INTO accounts (name) VALUES(?)`, a.Name)
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

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM accounts WHERE id=?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM balances WHERE account_id=?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM transactions WHERE account_id=?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
