package sqlite

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func (s *SQLiteStore) GetHistory(aID int) ([]mt.Balance, error) {

	rows, err := s.db.Query(`SELECT account_id,timestamp,value FROM balances WHERE account_id = ? ORDER BY timestamp DESC`, aID)
	if err != nil {
		return nil, err
	}

	balances := []mt.Balance{}
	var b mt.Balance

	for rows.Next() {
		if err = rows.Scan(&b.AccountID, &b.Timestamp, &b.Value); err != nil {
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
				FROM balances
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
				FROM transactions
				INNER JOIN operations op
				ON operation_id = op.id
				WHERE (
						to_id = :aID
						OR from_id = :aID
					)
					AND op.timestamp > (SELECT timestamp
					FROM (
							SELECT COUNT(),
								IFNULL(
									timestamp,
									STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-100 year')
								) AS timestamp
							FROM balances
							WHERE account_id = :aID
							ORDER BY timestamp DESC
							LIMIT 1
						)
					)
			)
		)`, sql.Named("aID", aID)).Scan(
		&b.Value)
	if err != nil {
		return nil, err
	}

	b.AccountID = null.IntFrom(int64(aID))
	b.Timestamp = mt.DateTime{Time: time.Now()}

	return &b, nil
}

func (s *SQLiteStore) AdjustBalance(b mt.Balance) error {

	if !b.AccountID.Valid {
		return echo.NewHTTPError(http.StatusBadRequest, "account ID must be non null")
	}
	currentBalance, err := s.GetBalance(int(b.AccountID.Int64))
	if err != nil {
		return err
	}

	delta := b.Value.Sub(currentBalance.Value)
	if delta.IsZero() {
		return echo.NewHTTPError(http.StatusBadRequest, "account balance is already at the specified value")
	}

	t := mt.Transaction{
		From:   mt.Account{},
		To:     mt.Account{},
		Amount: delta.Abs(),
	}

	world := null.IntFrom(0) //this should not be hard coded

	if delta.IsNegative() {
		t.From.ID = b.AccountID
		t.To.ID = world
	} else {
		t.From.ID = world
		t.To.ID = b.AccountID
	}

	op := mt.Operation{
		CreatedByID:  0, //?
		Timestamp:    b.Timestamp,
		Description:  b.Comment,
		Transactions: []mt.Transaction{t},
		Type:         mt.OpTypeBalanceAdjust,
		CategoryID:   0,
	}

	_, err = s.AddOperation(op)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) ComputeBalance(aID int) error {

	_, err := s.db.Exec(`INSERT INTO balances (account_id, value)
	SELECT :aID,
		last_balance + balance AS balance
	FROM (
			(
				SELECT value AS last_balance
				FROM balances
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
				FROM transactions
				INNER JOIN operations op
				ON operation_id = op.id
				WHERE (
						to_id = :aID
						OR from_id = :aID
					)
					AND t.timestamp > (
						SELECT timestamp
						FROM balances
						WHERE account_id = :aID
						ORDER BY timestamp DESC
						LIMIT 1
					)
			)
		)
		WHERE EXISTS (
		SELECT *
		FROM transactions
		INNER JOIN operations op
				ON operation_id = op.id
		WHERE (
				to_id = :aID
				OR from_id = :aID
			)
			AND op.timestamp > (
				SELECT timestamp
				FROM balances
				WHERE account_id = :aID
				ORDER BY timestamp DESC
				LIMIT 1
			)
	)`,
		sql.Named("aID", aID),
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *SQLiteStore) DeleteBalancesAfter(aID int, date mt.DateTime) error {
	_, err := s.db.Exec(`
			DELETE FROM balances
			WHERE account_id = ?
			AND timestamp >= ?
	`, aID, date)
	if err != nil {
		return err
	}

	return nil

}
