package sqlite

import (
	"database/sql"
	"fmt"

	mt "github.com/emarj/moneytracker"
	"gopkg.in/guregu/null.v4"

	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/emarj/moneytracker/datetime"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) getBalanceHistory(aID int64, limit int64) ([]mt.Balance, error) {

	stmt := jet.SELECT(jt.Balance.AllColumns).FROM(jt.Balance).WHERE(jt.Balance.AccountID.EQ(jet.Int(aID))).ORDER_BY(jt.Balance.Timestamp.DESC())

	if limit > 0 {
		stmt = stmt.LIMIT(limit)
	}

	balances := []mt.Balance{}
	err := stmt.Query(s.db, &balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func (s *SQLiteStore) GetBalanceHistory(aID int64) ([]mt.Balance, error) {
	return s.getBalanceHistory(aID, 0)
}

func (s *SQLiteStore) GetLastBalance(aID int64) (mt.Balance, error) {

	var b mt.Balance
	history, err := s.getBalanceHistory(aID, 1)
	if err != nil {
		return b, err
	}

	if len(history) == 0 {
		return b, mt.ErrNotFound
	}

	b = history[0]

	return b, nil
}

func getBalanceAt(db TXDB, aID int64, timestamp datetime.DateTime) (mt.Balance, error) {

	var cb mt.Balance

	row := db.QueryRow(`
		SELECT last_balance + delta AS current_balance, delta
		FROM (
			(
				SELECT
					timestamp AS last_timestamp,
					IFNULL(value,0) AS last_balance,
					COUNT()
				FROM balance
				WHERE account_id = :aID AND timestamp <= :timestamp
				ORDER BY timestamp DESC
				LIMIT 1
			), (
				SELECT IFNULL(
						SUM(
							CASE
								WHEN to_id = :aID THEN amount
								WHEN from_id = :aID THEN - amount
							END
						),
						0
					) AS delta,
					COUNT()
				FROM 'transaction'
				WHERE (to_id = :aID
					OR from_id = :aID)
					AND timestamp BETWEEN IFNULL((
									SELECT timestamp
									FROM balance
									WHERE account_id = :aID AND timestamp <= :timestamp
									ORDER BY timestamp DESC
									LIMIT 1),STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now','-1000 years')) AND :timestamp
			)
		);`,
		sql.Named("aID", aID),
		sql.Named("timestamp", timestamp.String()),
	)

	err := row.Scan(&cb.ValueAt.Value, &cb.Delta)
	if err != nil {
		return cb, err
	}

	cb.Timestamp = timestamp
	cb.AccountID = null.IntFrom(int64(aID))
	cb.IsComputed = true

	return cb, nil
}

func (s *SQLiteStore) GetValueAt(aID int64, timestamp datetime.DateTime) (mt.Balance, error) {
	return getBalanceAt(s.db, aID, timestamp)
}

func (s *SQLiteStore) GetValueNow(aID int64) (mt.Balance, error) {
	return s.GetValueAt(aID, datetime.Now())
}

func (s *SQLiteStore) SetBalance(b mt.Balance) error {

	err := s.AddOperation(&mt.Operation{
		Description: "Balance Adjust",
		Balances:    []mt.Balance{b},
		TypeID:      mt.OpTypeBalanceAdjust,
	})
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

func insertBalance(tx TXDB, balance *mt.Balance) error {

	b, err := getBalanceAt(tx, balance.AccountID.Int64, balance.Timestamp)
	if err != nil {
		return err
	}

	balance.Delta = b.Delta

	//TODO: We could update delta directly in the insert query. To do this

	stmt := jt.Balance.INSERT(jt.Balance.AllColumns).MODEL(balance)

	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) DeleteComputedBalancesAfter(aID int64, date datetime.DateTime) error {
	_, err := s.db.Exec(`
			DELETE FROM balance
			WHERE account_id = ?
			AND is_computed = TRUE
			AND timestamp >= ?
	`, aID, date)
	if err != nil {
		return err
	}

	return nil

}
