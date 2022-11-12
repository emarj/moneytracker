package sqlite

import (
	"fmt"

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

func (s *SQLiteStore) GetAccount(id int) (*mt.Account, error) {

	var a mt.Account

	err := s.db.QueryRow(`SELECT id,name,income - expense AS balance
	FROM
	((
		SELECT *
		FROM accounts
		WHERE id = ?), (
		SELECT  SUM(amount) AS income
		FROM transactions
		WHERE to_id = ?), (
		SELECT  SUM(amount) AS expense
		FROM transactions
		WHERE from_id = ?)
	)`, id, id, id).Scan(
		&a.ID,
		&a.Name,
		&a.Balance,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *SQLiteStore) GetAccountBalance(id int) (*mt.Account, error) {

	var a mt.Account

	err := s.db.QueryRow(`SELECT "value" + income - expense  FROM balances WHERE account_id = ? ORDER BY computed_at DESC LIMIT 1
	FROM
	(	
		(SELECT "value" FROM balances WHERE account_id = ? ORDER BY computed_at DESC LIMIT 1),
		(SELECT  SUM(amount) AS income
		FROM transactions
		WHERE to_id = ?),
		(SELECT  SUM(amount) AS expense
		FROM transactions
		WHERE from_id = ?)
	)`, id, id, id).Scan(
		&a.ID,
		&a.Name,
		&a.Balance,
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

func (s *SQLiteStore) AddBalance(b mt.Balance) (*mt.Balance, error) {

	if b.Computed {
		return nil, fmt.Errorf("b.Computed cannot be true")
	}

	b.Computed = false // This is not needed, but it's here for safety

	_, err := s.db.Exec(`INSERT INTO balances ("account_id","computed_at","value","computed","notes") VALUES(?,?,?,?,?)`,
		b.AccountID,
		b.Timestamp,
		b.Value,
		b.Computed, // b.Computed is ignored since this is clearly not computed
		b.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &b, nil

}

/*func (s *SQLiteStore) ComputeBalance(id int) (*mt.Balance, error) {

	_, err := s.db.Exec(`INSERT INTO balances ("account_id","computed_at","value","computed","notes") VALUES(?,?,?,?,?)`,
		b.AccountID,
		b.ComputedAt,
		b.Value,
		b.Computed, // b.Computed is ignored since this is clearly not computed
		b.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &b, nil

}*/

func (s *SQLiteStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
