package sqlite

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"ronche.se/moneytracker/model"
)

func (s *sqlite) TransactionGet(uid uuid.UUID) (*model.Transaction, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Preparex(
		`SELECT
			uuid,
			date_created,
			date_modified,
			date,
			t.user_id,
			t.user_name,
			amount,
			t.pm_id,
			pm_name,
			description,
			t.cat_id,
			cat_name,
			shared,
			geolocation,
			t.type_id,
			type_name,
			IFNULL(tx_uuid,"` + uuid.Nil.String() + `") AS tx_uuid,
			IFNULL(with_id,0) AS with_id,
			IFNULL(u.user_name,0) as with_name,
			IFNULL(quota,0) AS quota
		FROM
			(SELECT * from transactions t,users,types,paymentmethods,categories
				WHERE	t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id AND
						t.uuid=?
			) t
			LEFT OUTER JOIN shares s ON t.uuid = s.tx_uuid
			LEFT OUTER JOIN users u ON s.with_id = u.user_id`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(uid.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close() //should not be needed if we iterate over all rows

	type Result struct {
		model.Transaction
		model.Share
	}

	var t *model.Transaction
	flagTx := false

	for rows.Next() {
		var result Result
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}

		if !flagTx {
			t = &result.Transaction
			flagTx = true
		}
		if result.Share.TxID != uuid.Nil {
			if !t.Shared {
				return t, errors.New("transaction is not shared, but it has shares!")
			}
			t.Shares = append(t.Shares, &result.Share)
		}

	}

	if t.Shared && len(t.Shares) == 0 {
		return t, errors.New("transaction is shared, but it has no shares!")
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return t, err
}

func (s *sqlite) TransactionInsert(t *model.Transaction) error {
	id := uuid.NewV4()

	t.UUID = id
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return err
	}
	t.DateCreated.Time = time.Now().In(loc)
	t.DateModified = t.DateCreated

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.PrepareNamed(
		`INSERT INTO transactions(
			uuid,
			date_created,
			date_modified,
			date,
			user_id,
			amount,
			pm_id,
			description,
			cat_id,
			shared,
			geolocation,
			type_id
		) VALUES(
			:uuid,
			:date_created,
			:date_modified,
			:date,
			:user_id,
			:amount,
			:pm_id,
			:description,
			:cat_id,
			:shared,
			:geolocation,
			:type_id)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t)

	if err != nil {
		return err
	}

	if t.Shared && (t.Shares != nil) {
		query := `INSERT INTO shares(
			tx_uuid,
			with_id,
			quota) VALUES`

		vals := []interface{}{}
		for _, shr := range t.Shares {
			query += "(?,?,?),"
			vals = append(vals, t.UUID.String(), shr.WithID, shr.Quota)
		}

		query = query[0 : len(query)-1] //Remove last comma

		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(vals...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *sqlite) TransactionUpdate(t *model.Transaction) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return err
	}
	t.DateModified.Time = time.Now().In(loc)

	stmt, err := tx.PrepareNamed(
		`UPDATE transactions
		SET
			date_modified = :date_modified,
			date = :date,
			user_id = :user_id,
			amount = :amount,
			pm_id = :pm_id,
			description = :description,
			cat_id = :cat_id,
			shared = :shared,
			geolocation = :geolocation,
			type_id = :type_id
		WHERE uuid = :uuid`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	//If t was shared and now shares are remove we should account for that!
	//Probably get all shares by id, delete them and add them again is the quick and dirty way!
	//OR WITH ADD A SHARE UUID

	_, err = tx.Exec(`DELETE FROM shares WHERE tx_uuid =?`, t.UUID.String())
	if err != nil {
		return err
	}
	if t.Shared {
		if t.Shares == nil {
			return errors.New("Transaction is shared but it has no shares")
		}
		query := `INSERT INTO shares(
				tx_uuid,
				with_id,
				quota) VALUES`

		vals := []interface{}{}
		for _, shr := range t.Shares {
			query += "(?,?,?),"
			vals = append(vals, t.UUID.String(), shr.WithID, shr.Quota)
		}

		query = query[0 : len(query)-1] //Remove last comma

		stmt2, err := tx.Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt2.Exec(vals...)

		if err != nil {
			return err
		}
	} else if t.Shares != nil {
		return errors.New("Transaction is not shared but it has shares")

	}

	return nil
}

func (s *sqlite) TransactionDelete(id uuid.UUID) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec("DELETE FROM transactions WHERE uuid=?", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM shares WHERE tx_uuid=?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) TransactionsGetNOrderBy(limit int, orderBy string) ([]*model.Transaction, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Preparex(
		`SELECT
				uuid,
				date_created,
				date_modified,
				date,
				t.user_id,
				t.user_name,
				amount,
				t.pm_id,
				pm_name,
				description,
				t.cat_id,
				cat_name,
				shared,
				geolocation,
				t.type_id,
				type_name,
				IFNULL(tx_uuid,"` + uuid.Nil.String() + `") AS tx_uuid,
				IFNULL(with_id,0) AS with_id,
				IFNULL(u.user_name,0) as with_name,
				IFNULL(quota,0) AS quota
			FROM
				(SELECT * from transactions t,users,types,paymentmethods,categories
					WHERE	t.user_id=users.user_id AND
							t.type_id=types.type_id AND
							t.pm_id=paymentmethods.pm_id AND
							t.cat_id=categories.cat_id
					ORDER BY ` + orderBy + `
					LIMIT ?) t
					LEFT OUTER JOIN shares s ON t.uuid = s.tx_uuid
					LEFT OUTER JOIN users u ON s.with_id = u.user_id
			ORDER BY ` + orderBy)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //should not be needed if we iterate over all rows

	type Result struct {
		model.Transaction
		model.Share
	}

	var ts []*model.Transaction
	var prevUUID uuid.UUID

	for rows.Next() {
		var result Result

		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		curUUID := result.Transaction.UUID

		if prevUUID != curUUID {
			prevUUID = curUUID
			ts = append(ts, &result.Transaction)
		}
		if result.Share.TxID != uuid.Nil {
			i := len(ts) - 1
			ts[i].Shares = append(ts[i].Shares, &result.Share)
		}

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ts, err
}

func (s *sqlite) TransactionsGetNOrderByDate(limit int) ([]*model.Transaction, error) {
	return s.TransactionsGetNOrderBy(limit, "date DESC, date_created DESC")
}

func (s *sqlite) TransactionsGetNOrderByInserted(limit int) ([]*model.Transaction, error) {
	return s.TransactionsGetNOrderBy(limit, "date_created DESC, date DESC")
}

func (s *sqlite) TransactionsGetNOrderByModified(limit int) ([]*model.Transaction, error) {
	return s.TransactionsGetNOrderBy(limit, "date_modified DESC, date DESC")
}
