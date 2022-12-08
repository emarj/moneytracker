package sqlite

import (
	"database/sql"
	"fmt"

	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) GetAccounts() ([]mt.Account, error) {

	rows, err := s.db.Query(`SELECT  a.id,a.name,a.display_name,a.type,e.* FROM accounts a INNER JOIN entities e WHERE a.owner_id = e.id`)
	if err != nil {
		return nil, err
	}

	accounts := []mt.Account{}
	var a mt.Account

	for rows.Next() {
		if err = rows.Scan(&a.ID, &a.Name, &a.DisplayName, &a.Type, &a.Owner.ID, &a.Owner.Name, &a.Owner.System, &a.Owner.External); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccountsByEntity(eID int) ([]mt.Account, error) {

	rows, err := s.db.Query(`SELECT  id,name,display_name,type FROM accounts WHERE owner_id = ? AND is_system == FALSE`, eID)
	if err != nil {
		return nil, err
	}

	accounts := []mt.Account{}
	var a mt.Account

	for rows.Next() {
		if err = rows.Scan(&a.ID, &a.Name, &a.DisplayName, &a.Type); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccount(aID int) (*mt.Account, error) {

	var a mt.Account

	err := s.db.QueryRow(`SELECT id,name,display_name,type FROM accounts WHERE id = ?`, aID).Scan(
		&a.ID,
		&a.Name,
		&a.DisplayName,
		&a.Type,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *SQLiteStore) AddAccount(a mt.Account) (null.Int, error) {

	id := null.Int{}
	res, err := s.db.Exec(`INSERT INTO accounts (id,name,display_name,owner_id,is_system,is_world,type) VALUES(?,?,?,?,?,?,?)`,
		a.ID, a.Name, a.DisplayName, a.Owner.ID, a.IsWorld, a.IsSystem, a.Type)
	if err != nil {
		return id, err
	}

	id.Int64, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	id.Valid = true

	return id, nil

}

func (s *SQLiteStore) DeleteAccount(aID int, onlyIfEmpty bool) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	if onlyIfEmpty {
		row := tx.QueryRow(`SELECT  count()
						FROM transactions t
						WHERE from_id = :aID
						OR to_id = :aID`, sql.Named("aID", aID))
		var n int
		err = row.Scan(&n)
		if err != nil {
			return err
		}

		if n > 0 {
			return fmt.Errorf("impossible to delete account id=%d since there are %d transaction associated to it", aID, n)
		}

	}

	_, err = tx.Exec(`DELETE FROM accounts WHERE id=?`, aID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM balances WHERE account_id=?`, aID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
