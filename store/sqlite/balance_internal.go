package sqlite

import (
	"database/sql"

	mt "github.com/emarj/moneytracker"
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"

	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/emarj/moneytracker/timestamp"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func getBalanceHistory(txdb TXDB, aID int64, limit int64, isComputed null.Bool) ([]mt.Balance, error) {

	where := jt.Balance.AccountID.EQ(jet.Int(aID))
	if isComputed.Valid {
		where.AND(jt.Balance.IsComputed.EQ(jet.Int(Btoi(isComputed.Bool))))
	}

	stmt := jet.SELECT(jt.Balance.AllColumns, jt.Operation.AllColumns).
		FROM(jt.Balance.LEFT_JOIN(jt.Operation, jt.Balance.OperationID.EQ(jt.Operation.ID))).
		WHERE(where).
		ORDER_BY(jt.Balance.Timestamp.DESC())

	if limit > 0 {
		stmt = stmt.LIMIT(limit)
	}

	balances := []mt.Balance{}
	err := stmt.Query(txdb, &balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func getBalanceAt(db TXDB, aID int64, timestamp timestamp.Timestamp) (mt.Balance, error) {

	var cb mt.Balance
	var err error

	row := db.QueryRow(`SELECT last_balance + delta AS current_balance, delta
		FROM (
			(
				SELECT IFNULL((
					SELECT value
					FROM balance
					WHERE account_id = :aID AND timestamp <= :timestamp
					ORDER BY timestamp DESC
					LIMIT 1
				),0) AS last_balance
			), (
				SELECT IFNULL((
					SELECT	IFNULL(SUM(
							CASE
								WHEN to_id = :aID THEN amount
								WHEN from_id = :aID THEN - amount
							END
						),0)
					FROM 'transaction'
					WHERE (to_id = :aID
						OR from_id = :aID)
						AND timestamp BETWEEN IFNULL((
										SELECT timestamp
										FROM balance
										WHERE account_id = :aID AND timestamp <= :timestamp
										ORDER BY timestamp DESC
										LIMIT 1),STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now','-1000 years')) AND :timestamp
			),0) AS delta
		)
		);`,
		sql.Named("aID", aID),
		sql.Named("timestamp", timestamp.String()),
	)
	if err != nil {
		return cb, err
	}

	//TODO: Delta is nice to have but not really used in the application.
	// We *could* return (last_balance,delta) and let the user compute the current_balance.
	// This would make more sense but it is not really useful

	err = row.Scan(&cb.Value, &cb.Delta)
	if err != nil {
		return cb, err
	}

	cb.Timestamp = timestamp
	cb.AccountID = null.IntFrom(int64(aID))
	cb.IsComputed = true

	return cb, nil
}

func insertBalance(tx TXDB, balance *mt.Balance) error {

	var err error
	b, err := getBalanceAt(tx, balance.AccountID.Int64, balance.Timestamp)
	if err != nil {
		return err
	}

	balance.Delta = decimal.NewNullDecimal(balance.Value.Sub(b.Value))

	//TODO: We could update delta directly in the insert query. To do this
	// Here we do not return anything since there are no generated fields
	stmt := jt.Balance.INSERT(jt.Balance.AllColumns).MODEL(balance)

	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	return nil

}

func updateBalances(txdb TXDB, timestamp timestamp.Timestamp, aIDs ...int64) error {
	var err error
	for _, aID := range aIDs {
		err = deleteComputedBalances(txdb, aID, timestamp)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteComputedBalances(txdb TXDB, aID int64, timestamp timestamp.Timestamp) error {
	_, err := txdb.Exec(`
						DELETE FROM balance
						WHERE account_id = :aID
						AND is_computed = TRUE
						AND timestamp BETWEEN :timestamp AND (
											-- Find the first user non-computed balance
											IFNULL(SELECT timestamp FROM balance
											WHERE
												account_id = :aID
												AND is_computed = FALSE
												AND timestamp > :timestamp
											ORDER BY timestamp ASC
											LIMIT 1,STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '+1000 year')))
	`,
		sql.Named("aID", aID),
		sql.Named("timestamp", timestamp))
	if err != nil {
		return err
	}

	return nil

}
