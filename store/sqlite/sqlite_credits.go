package sqlite

import (
	mt "ronche.se/moneytracker"
)

const selectCreditsQuery string = `SELECT id,debtor_id,creditor_id,account_id,amount,description,operation_id FROM credits`

func (s *SQLiteStore) GetCredits() ([]mt.Credit, error) {

	rows, err := s.db.Query(selectCreditsQuery)
	if err != nil {
		return nil, err
	}

	credits := []mt.Credit{}
	var c mt.Credit

	for rows.Next() {
		if err = rows.Scan(&c.ID, &c.Debtor.ID, &c.Creditor.ID, &c.Account.ID, &c.Amount, &c.Description, &c.Operation.ID); err != nil {
			return nil, err
		}

		credits = append(credits, c)
	}

	return credits, nil
}

/*func (s *SQLiteStore) GetAccount(aID int) (*mt.Account, error) {

	var a mt.Account

	err := s.db.QueryRow(`SELECT id,name FROM accounts WHERE id = ?`, aID).Scan(
		&a.ID,
		&a.Name,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}*/

const insertCreditQuery string = `INSERT INTO credits (debtor_id,creditor_id,account_id,amount,description,operation_id) VALUES(?,?,?,?,?,?)`

func (s *SQLiteStore) AddCredit(c mt.Credit) (int, error) {

	res, err := s.db.Exec(insertCreditQuery,
		c.Debtor.ID,
		c.Creditor.ID,
		c.Account.ID,
		c.Amount,
		c.Description,
		c.Operation.ID,
	)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil

}
