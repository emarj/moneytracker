SELECT *
from accounts;
SELECT *
from balances;
SELECT *
from transactions;
SELECT *
from operations;
-- Get Current Balance
SELECT last_balance + balance AS balance
FROM (
		(
			SELECT COUNT(),
				IFNULL(value, 0) AS last_balance
			FROM balances
			WHERE account_id = 1001
			ORDER BY timestamp DESC
			LIMIT 1
		), (
			SELECT IFNULL(
					SUM(
						CASE
							WHEN to_id = 1001 THEN amount
							WHEN from_id = 1001 THEN - amount
						END
					),
					0
				) AS balance
			FROM transactions
				INNER JOIN operations op ON operation_id = op.id
			WHERE (
					to_id = 1001
					OR from_id = 1001
				)
				AND op.timestamp > (
					SELECT timestamp
					FROM (
							SELECT COUNT(),
								IFNULL(
									timestamp,
									STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-100 year')
								) AS timestamp
							FROM balances
							WHERE account_id = 1001
							ORDER BY timestamp DESC
							LIMIT 1
						)
				)
		)
	);
--
SELECT *
from transactions;
SELECT *
from balances;
-- Update Balance
INSERT INTO balances (account_id, value)
SELECT 1,
	last_balance + balance AS balance
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
							WHEN from_id = 1 THEN - amount
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
--
SELECT * from balances b INNER JOIN accounts a ON b.account_id = a.id ;
--
INSERT INTO balances (account_id, value)
SELECT 1,
	last_balance + balance AS balance
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
							WHEN from_id = 1 THEN - amount
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
INSERT INTO transactions (from_id,to_id,operation_id,amount)
SELECT 0,
	1,
	-1,
	(5000 - value)
FROM (
		SELECT value
		FROM balances
		WHERE account_id = 1
		ORDER BY timestamp DESC
		LIMIT 1
	)
WHERE (
		SELECT value
		FROM balances
		WHERE account_id = 1
		ORDER BY timestamp DESC
		LIMIT 1
	) != 5000;
--
SELECT * from transactions ORDER BY timestamp DESC;
SELECT *
from balances
ORDER BY timestamp DESC;