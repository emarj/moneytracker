SELECT *
from accounts;
SELECT *
from balances;
-- Get Current Balance
SELECT last_balance + income - expense AS balance
FROM (
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
-- Update Balance
INSERT INTO balances (account_id, value, computed)
SELECT 1,
	last_balance + income - expense AS balance,
	TRUE
FROM (
		(
			SELECT value AS last_balance,
				timestamp AS last_computed
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
-- Delete Balances until date or the first non computed one
DELETE FROM balances
WHERE timestamp >= "2006-01-02T15:04:05.999Z"
	AND computed = TRUE;