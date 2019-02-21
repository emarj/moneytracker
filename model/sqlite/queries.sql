/*Expenses View for user*/

SELECT  t.date,t.uuid,t.amount,t.user_id,s.shared_quota,s.user_id
FROM shares s INNER JOIN transactions t ON t.uuid = s.transaction_uuid
WHERE s.user_id=2
UNION
SELECT date, uuid,amount,user_id, 0 AS shared_quota, NULL as user_ID
FROM transactions
WHERE user_id=2
ORDER BY date DESC

/* Get one expense*/
SELECT  *
		FROM (SELECT * FROM transactions t,users,types,paymentmethods,categories
		WHERE 
						t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id
						ORDER BY t.date DESC
						LIMIT 10) t, users u LEFT OUTER JOIN shares s ON t.uuid = s.tx_uuid
						WHERE u.user_id = s.with_id
		ORDER BY t.date DESC
		
/*Get one TX and username*/
	SELECT
		uuid,
		date_created,
		date,
		ut.user_id,
		ut.user_name,
		amount,
		t.pm_id,
		pm_name,
		description,
		t.cat_id,
		cat_name,
		shared,
		geolocation,
		t.type_id,
		type_name,
		tx_uuid,
		with_id,
		us.user_name AS with_name,
		quota
		FROM users ut,types,paymentmethods,categories,transactions t INNER JOIN shares s ON t.uuid = s.tx_uuid,users us
		WHERE 
						us.user_id = s.with_id AND
						t.user_id=ut.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id
		UNION
		SELECT 		uuid,
		date_created,
		date,
		t.user_id,
		user_name,
		amount,
		t.pm_id,
		pm_name,
		description,
		t.cat_id,
		cat_name,
		shared,
		geolocation,
		t.type_id,
		type_name,
		NULL AS tx_uuid, 0 AS with_id,"" AS with_name, 0 AS quota
		FROM users,types,paymentmethods,categories,transactions t 
		WHERE 
						t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id
		

/*Get total shared amount*/
SELECT  t.date,t.uuid,t.amount,t.user_id,SUM(s.shared_quota)
FROM shares s INNER JOIN transactions t ON t.uuid = s.transaction_UUID
GROUP BY t.uuid


/*Balance*/
SELECT SUM(	
					CASE WHEN type_id = 0 THEN
								CASE WHEN user_id=1 THEN 
									CASE WHEN shared=1 THEN amount-quota
									ELSE amount END
								ELSE quota END
					ELSE
					(CASE WHEN type_id = 1
							THEN
								CASE WHEN user_id = 1
									THEN amount
									ELSE -amount
								END
							ELSE (CASE WHEN type_id = 2 THEN
								-amount
							END)
					END)
					END
			) AS balance

		FROM transactions t LEFT OUTER JOIN shares s
		ON t.uuid = s.tx_uuid
		WHERE	t.user_id=1 OR s.with_id=1

/*Credit user 1 towards 2*/
SELECT SUM(	
					CASE WHEN type_id = 0 THEN
								CASE WHEN user_id=1 THEN 
								quota
								ELSE -quota END
					ELSE
					(CASE WHEN type_id = 1
							THEN
								CASE WHEN user_id = 1
									THEN amount
									ELSE -amount
								END
					END)
					END
			) AS Credit

		FROM transactions t LEFT OUTER JOIN shares s
		ON t.uuid = s.tx_uuid
		WHERE	t.shared = 1 AND ((t.user_id=1 AND s.with_id=2) OR (t.user_id=2 AND s.with_id=1))