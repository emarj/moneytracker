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
		FROM users,types,paymentmethods,categories,transactions t INNER JOIN shares s ON t.uuid = s.tx_uuid
		WHERE 
						t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id AND
						t.uuid="bdbd93db-c5fe-4165-b237-99b27883fac3"
		UNION
		SELECT *, NULL AS tx_uuid, 0 AS with_id,0 AS quota
		FROM users,types,paymentmethods,categories,transactions t 
		WHERE 
						t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id AND
						t.uuid= "bdbd93db-c5fe-4165-b237-99b27883fac3"
		ORDER BY date DESC

/*Get total shared amount*/
SELECT  t.date,t.uuid,t.amount,t.user_id,SUM(s.shared_quota)
FROM shares s INNER JOIN transactions t ON t.uuid = s.transaction_UUID
GROUP BY t.uuid