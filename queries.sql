/*Expenses View for user*/

SELECT  t.date,t.UUID,t.amount,t.user_id,s.shared_quota,s.user_ID
FROM sharings s INNER JOIN transactions t ON t.UUID = s.transaction_UUID
WHERE s.user_ID=2
UNION
SELECT date, UUID,amount,user_id, 0 AS shared_quota, NULL as user_ID
FROM transactions
WHERE user_id=2
ORDER BY date DESC

/*Get total shared amount*/
SELECT  t.date,t.UUID,t.amount,t.user_id,SUM(s.shared_quota)
FROM sharings s INNER JOIN transactions t ON t.UUID = s.transaction_UUID
GROUP BY t.uuid