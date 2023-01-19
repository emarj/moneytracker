package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"
	"gopkg.in/guregu/null.v4"

	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/emarj/moneytracker/timestamp"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetBalanceHistory(aID int64) ([]mt.Balance, error) {
	return getBalanceHistory(s.db, aID, 0, null.Bool{})
}

func (s *SQLiteStore) GetLastBalance(aID int64) (mt.Balance, error) {

	var b mt.Balance
	history, err := getBalanceHistory(s.db, aID, 1, null.Bool{})
	if err != nil {
		return b, err
	}

	if len(history) == 0 {
		return b, mt.ErrNotFound
	}

	b = history[0]

	return b, nil
}

func (s *SQLiteStore) GetBalanceAt(aID int64, timestamp timestamp.Timestamp) (mt.Balance, error) {
	return getBalanceAt(s.db, aID, timestamp)
}

func (s *SQLiteStore) GetBalanceNow(aID int64) (mt.Balance, error) {
	return s.GetBalanceAt(aID, timestamp.Now())
}

func (s *SQLiteStore) SetBalance(b *mt.Balance) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	desc := "Balance Adjust"
	if b.Comment != "" {
		desc = b.Comment
	}
	op := mt.Operation{
		Description: desc,
		TypeID:      mt.OpTypeBalanceAdjust,
	}

	err = insertOperation(tx, &op)
	if err != nil {
		return err
	}

	b.OperationID = op.ID
	b.Operation = &op

	err = insertBalance(tx, b)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) SnapshotBalance(aID int64) error {

	return fmt.Errorf("not implemented")

	/* _, err := s.db.Exec(`INSERT INTO balance (account_id, value,is_computed)
	SELECT :aID,
		last_balance + balance AS balance,
		TRUE
	FROM (
			(
				SELECT value AS last_balance
				FROM balance
				WHERE account_id = :aID
				ORDER BY timestamp DESC
				LIMIT 1
			), (
				SELECT IFNULL(
						SUM(
							CASE
								WHEN to_id = :aID THEN amount
								WHEN from_id = :aID THEN -amount
							END
						),
						0
					) AS balance
				FROM 'transaction' AS t
				WHERE (
						to_id = :aID
						OR from_id = :aID
					)
					AND t.timestamp > (
						SELECT IFNULL(
										timestamp,
										STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-100 year')
									) AS timestamp
					FROM (
						SELECT timestamp,COUNT()
						FROM (
								SELECT timestamp
								FROM balance
								WHERE account_id = :aID
								ORDER BY timestamp DESC
								LIMIT 1
							)
					)
					)
			)
		)
		WHERE EXISTS (
		SELECT *
		FROM 'transaction'
		WHERE (
				to_id = :aID
				OR from_id = :aID
			)
			AND timestamp > (
				SELECT IFNULL(
										timestamp,
										STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-100 year')
									) AS timestamp
					FROM (
						SELECT timestamp,COUNT()
						FROM (
								SELECT timestamp
								FROM balance
								WHERE account_id = :aID
								ORDER BY timestamp DESC
								LIMIT 1
							)
					)
			)
	)`,
		sql.Named("aID", aID),
	)
	if err != nil {
		return err
	}

	return nil */

}

func (s *SQLiteStore) DeleteBalance(aID int64, timestamp timestamp.Timestamp) error {
	stmt := jt.Balance.DELETE().
		WHERE(
			jt.Balance.AccountID.EQ(jet.Int(aID)).
				AND(jt.Balance.Timestamp.EQ(jet.String(timestamp.String()))),
		)

	_, err := stmt.Exec(s.db)
	if err != nil {
		return err
	}

	return nil
}
