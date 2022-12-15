SELECT *
from entity;
SELECT *
from account;
SELECT *
from balance;
SELECT *
from 'transaction';
SELECT *
from operation;
-- Get Current Balance
SELECT last_balance + balance AS balance
FROM (
		(
			SELECT COUNT(),
				IFNULL(value, 0) AS last_balance
			FROM balance
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
			FROM 'transaction'
				INNER JOIN operation op ON operation_id = op.id
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
							FROM balance
							WHERE account_id = 1001
							ORDER BY timestamp DESC
							LIMIT 1
						)
				)
		)
	);
--
SELECT *
from 'transaction';
SELECT *
from balance;
-- Update Balance
INSERT INTO balance (account_id, value)
SELECT 1,
	last_balance + balance AS balance
FROM (
		(
			SELECT value AS last_balance
			FROM balance
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
			FROM 'transaction'
			WHERE (
					to_id = 1
					OR from_id = 1
				)
				AND timestamp > (
					SELECT timestamp
					FROM balance
					WHERE account_id = 1
					ORDER BY timestamp DESC
					LIMIT 1
				)
		)
	)
WHERE EXISTS (
		SELECT *
		FROM 'transaction'
		WHERE (
				to_id = 1
				OR from_id = 1
			)
			AND timestamp > (
				SELECT timestamp
				FROM balance
				WHERE account_id = 1
				ORDER BY timestamp DESC
				LIMIT 1
			)
	);
--
	SELECT *
from balance;
-- Delete Balances until date or the first non computed one
DELETE FROM balance
WHERE account_id = 1
	AND timestamp >= "2006-01-02T15:04:05.999Z"
	AND computed = TRUE;
--
SELECT  timestamp,value,computed,notes FROM balance WHERE account_id = 1
--
SELECT * from balance b INNER JOIN account a ON b.account_id = a.id ;
--
INSERT INTO balance (account_id, value)
SELECT 1,
	last_balance + balance AS balance
FROM (
		(
			SELECT value AS last_balance
			FROM balance
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
			FROM 'transaction'
			WHERE (
					to_id = 1
					OR from_id = 1
				)
				AND timestamp > (
					SELECT timestamp
					FROM balance
					WHERE account_id = 1
					ORDER BY timestamp DESC
					LIMIT 1
				)
		)
	)
WHERE EXISTS (
		SELECT *
		FROM 'transaction'
		WHERE (
				to_id = 1
				OR from_id = 1
			)
			AND timestamp > (
				SELECT timestamp
				FROM balance
				WHERE account_id = 1
				ORDER BY timestamp DESC
				LIMIT 1
			)
	);
--
INSERT INTO 'transaction' (from_id,to_id,operation_id,amount)
SELECT 0,
	1,
	-1,
	(5000 - value)
FROM (
		SELECT value
		FROM balance
		WHERE account_id = 1
		ORDER BY timestamp DESC
		LIMIT 1
	)
WHERE (
		SELECT value
		FROM balance
		WHERE account_id = 1
		ORDER BY timestamp DESC
		LIMIT 1
	) != 5000;
--
SELECT * from 'transaction' ORDER BY timestamp DESC;
SELECT *
from balance
ORDER BY timestamp DESC;