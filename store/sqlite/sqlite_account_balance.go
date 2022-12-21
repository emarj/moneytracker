package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"

	jt "ronche.se/moneytracker/.gen/table"
	"ronche.se/moneytracker/datetime"
)

func (s *SQLiteStore) GetHistory(aID int) ([]mt.Balance, error) {

	rows, err := s.db.Query(`SELECT account_id,timestamp,value,delta FROM balance WHERE account_id = ? ORDER BY timestamp DESC`, aID)
	if err != nil {
		return nil, err
	}

	balances := []mt.Balance{}
	var b mt.Balance

	for rows.Next() {
		if err = rows.Scan(&b.AccountID, &b.Timestamp, &b.Value, &b.Delta); err != nil {
			return nil, err
		}

		balances = append(balances, b)
	}

	return balances, nil
}
func (s *SQLiteStore) GetBalanceNow(aID int) (*mt.Balance, error) {
	return s.GetBalanceAt(aID, datetime.FromTime(time.Now()))
}

func (s *SQLiteStore) GetBalanceAt(aID int, timestamp datetime.DateTime) (*mt.Balance, error) {

	var b mt.Balance

	err := s.db.QueryRow(`SELECT last_balance + balance AS balance
	FROM (
			(
				SELECT
					value AS last_balance
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
					) AS balance
				FROM 'transaction'
				WHERE (to_id = :aID
					OR from_id = :aID)
					AND timestamp BETWEEN (
									SELECT timestamp
									FROM balance
									WHERE account_id = :aID AND timestamp <= :timestamp
									ORDER BY timestamp DESC
									LIMIT 1) AND :timestamp
			)
		);`, sql.Named("aID", aID), sql.Named("timestamp", timestamp.String())).Scan(&b.Value)
	if err != nil {
		return nil, err
	}

	b.AccountID = null.IntFrom(int64(aID))
	b.Timestamp = datetime.DateTime{}

	return &b, nil
}

func (s *SQLiteStore) SetBalance(b mt.Balance) error {

	err := s.AddOperation(&mt.Operation{
		Title:       "Balance Adjust",
		Description: "",
		Balances:    []mt.Balance{b},
		TypeID:      mt.OpTypeBalance,
	})
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) SnapshotBalance(aID int) error {

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

func insertBalances(db DB, balances []mt.Balance) error {

	stmt := jt.Balance.INSERT(jt.Balance.AllColumns).MODELS(balances)

	_, err := stmt.Exec(db)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) DeleteComputedBalancesAfter(aID int, date datetime.DateTime) error {
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
