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
SELECT last_balance,delta, last_balance + delta AS balance
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
			SELECT IFNULL(
					SUM(
						CASE
							WHEN to_id = 1 THEN amount
							WHEN from_id = 1 THEN - amount
						END
					),
					0
				) AS delta,
				COUNT()
			FROM 'transaction'
			WHERE (to_id = 1
				OR from_id = 1)
				AND timestamp BETWEEN (
								SELECT timestamp
								FROM balance
								WHERE account_id = 1 AND timestamp <= STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')
								ORDER BY timestamp DESC
								LIMIT 1
								) AND STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')
		)
	);
--
SELECT *
from 'transaction';
SELECT *
from balance;
---

SELECT timestamp
							FROM balance
							WHERE account_id = 8
							ORDER BY timestamp DESC
							LIMIT 1
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
--
SELECT * FROM balance
WHERE account_id = 2
AND timestamp BETWEEN
					STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-1 day')
					AND (SELECT timestamp FROM balance
										WHERE
											account_id = 2
											AND is_computed = FALSE
											AND timestamp > STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now', '-1 day')
										ORDER BY timestamp ASC
										LIMIT 1)