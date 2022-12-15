package sqlite

import (
	"database/sql"
	"fmt"

	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"

	jt "ronche.se/moneytracker/.gen/table"

	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetAccounts() ([]mt.Account, error) {

	stmt := jet.SELECT(jt.Account.AllColumns,
		jt.Entity.AllColumns,
	).FROM(jt.Account.INNER_JOIN(jt.Entity, jt.Entity.ID.EQ(jt.Account.ID)))

	accounts := []mt.Account{}

	err := stmt.Query(s.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SQLiteStore) GetAccountsByEntity(eID int) ([]mt.Account, error) {

	rows, err := s.db.Query(`SELECT  id,name,display_name,is_credit FROM account WHERE owner_id = ? AND is_system == FALSE`, eID)
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

	err := s.db.QueryRow(`SELECT id,name FROM account WHERE id = ?`, aID).Scan(
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
	res, err := s.db.Exec(`INSERT INTO account (id,name,display_name,owner_id,is_system,is_world,type) VALUES(?,?,?,?,?,?,?)`,
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
						FROM 'transaction' t
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

	_, err = tx.Exec(`DELETE FROM account WHERE id=?`, aID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM balance WHERE account_id=?`, aID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
