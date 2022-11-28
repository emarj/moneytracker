package sqlite

import (
	"time"

	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) GetBalances(aID int) ([]mt.Balance, error) {

	rows, err := s.db.Query(`SELECT timestamp,value,operation_id FROM balance WHERE account_id = ? ORDER BY timestamp DESC`, aID)
	if err != nil {
		return nil, err
	}

	balances := []mt.Balance{}
	var b mt.Balance
	b.AccountID = aID

	for rows.Next() {
		if err = rows.Scan(&b.Timestamp, &b.Value); err != nil {
			return nil, err
		}

		balances = append(balances, b)
	}

	return balances, nil
}

func (s *SQLiteStore) GetBalance(aID int) (*mt.Balance, error) {

	var b mt.Balance

	err := s.db.QueryRow(`SELECT  last_balance + balance AS balance
	FROM (
			(
				SELECT COUNT(),IFNULL(value, 0) AS last_balance
				FROM balance
				WHERE account_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			), (
				SELECT IFNULL(
						SUM(
							CASE
								WHEN to_id = ? THEN amount
								WHEN from_id = ? THEN -amount
							END
						),
						0
					) AS balance
				FROM "transaction"
				INNER JOIN operation op
				ON operation_id = op.id
				WHERE (
						to_id = ?
						OR from_id = ?
					)
					AND op.timestamp > (SELECT timestamp
					FROM (
							SELECT COUNT(),
								IFNULL(
									timestamp,
									STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-100 year')
								) AS timestamp
							FROM balance
							WHERE account_id = ?
							ORDER BY timestamp DESC
							LIMIT 1
						)
					)
			)
		)`, aID, aID, aID, aID, aID, aID).Scan(
		&b.Value)
	if err != nil {
		return nil, err
	}

	b.AccountID = aID
	b.Timestamp = mt.DateTime{Time: time.Now()}

	return &b, nil
}

func (s *SQLiteStore) AddBalance(b mt.Balance) error {

	var fromID, toID int

	if b.Value.IsNegative() {
		fromID = b.AccountID
	} else {
		toID = b.AccountID
	}

	_, err := s.db.Exec(`INSERT INTO "transaction" ("timestamp","from_id","to_id","operation_id","amount") SELECT ?,?,?,?,? - amount FROM balances WHERE account_id = ? ORDER BY timestamp DESC LIMIT 1`,
		b.AccountID,
		b.Timestamp,
		fromID,
		toID,
		b.Value,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) ComputeBalance(aID int) error {

	_, err := s.db.Exec(`INSERT INTO balance (account_id, value, computed)
	SELECT ?,
		last_balance + balance AS balance,
		TRUE
	FROM (
			(
				SELECT value AS last_balance
				FROM balance
				WHERE account_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			), (
				SELECT IFNULL(
						SUM(
							CASE
								WHEN to_id = ? THEN amount
								WHEN from_id = ? THEN -amount
							END
						),
						0
					) AS balance
				FROM "transaction"
				INNER JOIN operation op
				ON operation_id = op.id
				WHERE (
						to_id = ?
						OR from_id = ?
					)
					AND t.timestamp > (
						SELECT timestamp
						FROM balance
						WHERE account_id = ?
						ORDER BY timestamp DESC
						LIMIT 1
					)
			)
		)
		WHERE EXISTS (
		SELECT *
		FROM "transaction"
		INNER JOIN operation op
				ON operation_id = op.id
		WHERE (
				to_id = ?
				OR from_id = ?
			)
			AND op.timestamp > (
				SELECT timestamp
				FROM balance
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
			DELETE FROM balance
			WHERE account_id = ?
			AND timestamp >= ?
			AND computed = TRUE
	`, aID, date)
	if err != nil {
		return err
	}

	return nil

}
