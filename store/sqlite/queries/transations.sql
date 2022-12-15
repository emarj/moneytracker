SELECT *
FROM entity;
--
SELECT *
FROM account;
--
SELECT *
FROM 'transaction';
--
SELECT *
FROM operation;
--
SELECT  t.*,
op.*,
fa.display_name AS from_name,
ta.display_name AS to_name,
fa.owner_id AS from_owner_id,
ta.owner_id AS to_owner_id
FROM 'transaction' t
    INNER JOIN operation op ON t.operation_id = op.id
    INNER JOIN account AS fa ON t.from_id = fa.id
    INNER JOIN account AS ta ON t.to_id = ta.id
WHERE from_owner_id = 1
    OR to_owner_id = 1
ORDER BY t.timestamp DESC,
    op.id,
    t.id;
-- Get Transactions by Account
SELECT id,
    timestamp,
    from_id,
    to_id,
    amount,
    operation_id
FROM 'transaction'
WHERE from_id = 1
    OR to_id = 1
ORDER BY timestamp DESC;
-- Get Transactions By Entity
SELECT t.*,
    fa.display_name AS from_name,
    ta.display_name AS to_name,
    fa.owner_id AS from_owner_id,
    ta.owner_id AS to_owner_id
FROM 'transaction' AS t
    INNER JOIN account AS fa ON t.from_id = fa.id
    INNER JOIN account AS ta ON t.to_id = ta.id
WHERE from_owner_id = 1
    OR to_owner_id = 1;
--- Get Transactions by Entity with internal/income/expense indication (NOT USED)
SELECT t.*,
    fa.display_name AS from_name,
    ta.display_name AS to_name,
    fa.owner_id AS from_owner_id,
    ta.owner_id AS to_owner_id,
    CASE
        WHEN (
            fa.owner_id = 1
            AND ta.owner_id = 1
        ) THEN '-'
        WHEN fa.owner_id = 1 THEN '->'
        WHEN ta.owner_id = 1 THEN '<-'
    END AS direction
FROM 'transaction' AS t
    INNER JOIN account AS fa ON t.from_id = fa.id
    INNER JOIN account AS ta ON t.to_id = ta.id
WHERE from_owner_id = 1
    OR to_owner_id = 1;
    --
    INSERT INTO account (name,display_name,owner_id,is_system,is_credit,is_world) VALUES("prova","ssd",9,TRUE,FALSE,FALSE)
     RETURNING *;
     ----
     SELECT  count()
					FROM 'transaction' t
						WHERE from_id = 56
						OR to_id = 99;
