SELECT *
from accounts;
SELECT *
from balances;
SELECT *
from transactions;
-- Get Current Balance
SELECT last_balance + balance AS balance
FROM (
		(
			SELECT value AS last_balance
			FROM balances
			WHERE account_id = 1
			ORDER BY timestamp DESC
			LIMIT 1
		), (
			SELECT IFNULL(
					SUM(
						CASE
							WHEN to_id = 1 THEN amount
							WHEN from_id = 1 THEN amount
						END
					),
					0
				) AS balance
			FROM transactions
			WHERE (
					to_id = 1
					OR from_id = 1
				)
				AND timestamp > (
					SELECT timestamp
					FROM balances
					WHERE account_id = 1
					ORDER BY timestamp DESC
					LIMIT 1
				)
		)
	);
-- Get Account with current balance
SELECT *,
	last_balance + income - expense AS balance
FROM (
		(
			SELECT *
			FROM accounts
			WHERE id = 1
		),
		(
			SELECT value AS last_balance,
				timestamp AS last_balance_date
			FROM balances
			WHERE account_id = 1
			ORDER BY timestamp DESC
			LIMIT 1
		), (
			SELECT SUM(amount) AS income
			FROM transactions
			WHERE to_id = 1
		),
		(
			SELECT SUM(amount) AS expense
			FROM transactions
			WHERE from_id = 1
		)
	);
--
SELECT *
from transactions;
SELECT *
from balances;
-- Update Balance
INSERT INTO balances (account_id, value, computed)
SELECT 1,
	last_balance + balance AS balance,
	TRUE
FROM (
		(
			SELECT value AS last_balance
			FROM balances
			WHERE account_id = 1
			ORDER BY timestamp DESC
			LIMIT 1
		), (
			SELECT IFNULL(
					SUM(
						CASE
							WHEN to_id = 1 THEN amount
							WHEN from_id = 1 THEN amount
						END
					),
					0
				) AS balance
			FROM transactions
			WHERE (
					to_id = 1
					OR from_id = 1
				)
				AND timestamp > (
					SELECT timestamp
					FROM balances
					WHERE account_id = 1
					ORDER BY timestamp DESC
					LIMIT 1
				)
		)
	)
WHERE EXISTS (
		SELECT *
		FROM transactions
		WHERE (
				to_id = 1
				OR from_id = 1
			)
			AND timestamp > (
				SELECT timestamp
				FROM balances
				WHERE account_id = 1
				ORDER BY timestamp DESC
				LIMIT 1
			)
	);
--
	SELECT *
from balances;
-- Delete Balances until date or the first non computed one
DELETE FROM balances
WHERE account_id = 1
	AND timestamp >= "2006-01-02T15:04:05.999Z"
	AND computed = TRUE;
--
SELECT  timestamp,value,computed,notes FROM balances WHERE account_id = 1