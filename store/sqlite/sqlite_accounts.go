package sqlite

import (
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
		if err = rows.Scan(&a.ID, &a.Name, &a.DisplayName, &a.IsCredit); err != nil {
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
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *SQLiteStore) AddAccount(a mt.Account) (null.Int, error) {

	id := null.Int{}
	res, err := s.db.Exec(`INSERT INTO account (id,name,display_name,owner_id,is_system,is_world,is_credit) VALUES(?,?,?,?,?,?,?)`,
		a.ID, a.Name, a.DisplayName, a.Owner.ID, a.IsWorld, a.IsSystem, a.IsCredit)
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

func (s *SQLiteStore) DeleteAccount(id int) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM account WHERE id=?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM balance WHERE account_id=?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM "transaction" WHERE account_id=?`, id)
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
