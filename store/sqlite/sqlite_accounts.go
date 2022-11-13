package sqlite

import (
	"fmt"
	"time"

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

func (s *SQLiteStore) GetAccountsOfEntity(eID int) ([]mt.Account, error) {

	rows, err := s.db.Query(`SELECT  id,name,display_name FROM accounts WHERE owner_id = ?`, eID)
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

func (s *SQLiteStore) GetBalances(aID int) ([]mt.Balance, error) {

	rows, err := s.db.Query(`SELECT  timestamp,value,computed,notes FROM balances WHERE account_id = ?`, aID)
	if err != nil {
		return nil, err
	}

	balances := []mt.Balance{}
	var b mt.Balance
	b.AccountID = aID

	for rows.Next() {
		if err = rows.Scan(&b.Timestamp, &b.Value, &b.IsComputed, &b.Notes); err != nil {
			return nil, err
		}

		balances = append(balances, b)
	}

	return balances, nil
}

func (s *SQLiteStore) GetBalance(aID int) (*mt.Balance, error) {

	var b mt.Balance

	err := s.db.QueryRow(`SELECT last_balance + balance AS balance
	FROM (
			(
				SELECT value AS last_balance
				FROM balances
				WHERE account_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			), (
				SELECT IFNULL(
						SUM(
							CASE
								WHEN to_id = ? THEN amount
								WHEN from_id = ? THEN amount
							END
						),
						0
					) AS balance
				FROM transactions
				WHERE (
						to_id = ?
						OR from_id = ?
					)
					AND timestamp > (
						SELECT timestamp
						FROM balances
						WHERE account_id = ?
						ORDER BY timestamp DESC
						LIMIT 1
					)
			)
		)`, aID, aID, aID, aID, aID, aID).Scan(
		&b.Value)
	if err != nil {
		return nil, err
	}

	b.AccountID = aID
	b.Timestamp = mt.DateTime{Time: time.Now()}
	b.IsComputed = true

	return &b, nil
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

func (s *SQLiteStore) AddBalance(b mt.Balance) error {

	if b.Value == nil {
		return fmt.Errorf("Balance.Value cannot be nil")
	}

	_, err := s.db.Exec(`INSERT INTO balances ("account_id","timestamp","value","computed","notes") VALUES(?,?,?,?,?)`,
		b.AccountID,
		b.Timestamp,
		b.Value,
		false, // this is not computed
		b.Notes,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) ComputeBalance(aID int) error {

	_, err := s.db.Exec(`INSERT INTO balances (account_id, value, computed)
	SELECT ?,
		last_balance + balance AS balance,
		TRUE
	FROM (
			(
				SELECT value AS last_balance
				FROM balances
				WHERE account_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			), (
				SELECT IFNULL(
						SUM(
							CASE
								WHEN to_id = ? THEN amount
								WHEN from_id = ? THEN amount
							END
						),
						0
					) AS balance
				FROM transactions
				WHERE (
						to_id = ?
						OR from_id = ?
					)
					AND timestamp > (
						SELECT timestamp
						FROM balances
						WHERE account_id = ?
						ORDER BY timestamp DESC
						LIMIT 1
					)
			)
		)
		WHERE EXISTS (
		SELECT *
		FROM transactions
		WHERE (
				to_id = ?
				OR from_id = ?
			)
			AND timestamp > (
				SELECT timestamp
				FROM balances
				WHERE account_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			)
	)`,
		aID, aID, aID, aID, aID, aID, aID, aID, aID, aID,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) DeleteBalancesAfter(aID int, date mt.DateTime) error {
	_, err := s.db.Exec(`
			DELETE FROM balances
			WHERE account_id = ?
			AND timestamp >= ?
			AND computed = TRUE
	`, aID, date)
	if err != nil {
		return err
	}

	return nil

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
