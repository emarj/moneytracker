/*Expenses View for user*/

SELECT  t.date,t.UUID,t.amount,t.user_id,s.shared_quota,s.user_id
FROM shares s INNER JOIN transactions t ON t.uuid = s.transaction_uuid
WHERE s.user_id=2
UNION
SELECT date, UUID,amount,user_id, 0 AS shared_quota, NULL as user_ID
FROM transactions
WHERE user_id=2
ORDER BY date DESC

/* Get one expense*/
SELECT  t.date,t.UUID,t.amount,t.user_id,s.quota,s.user_id
FROM shares s INNER JOIN transactions t ON t.uuid = s.tx_uuid
WHERE tx_uuid="ec9291eb-bdef-44b8-92f6-614b41d6e98c"
UNION
SELECT date, UUID,amount,user_id, 0 AS quota, NULL as user_ID
FROM transactions
WHERE uuid="ec9291eb-bdef-44b8-92f6-614b41d6e98c"
ORDER BY date DESC

/*Get total shared amount*/
SELECT  t.date,t.UUID,t.amount,t.user_id,SUM(s.shared_quota)
FROM shares s INNER JOIN transactions t ON t.UUID = s.transaction_UUID
GROUP BY t.uuid