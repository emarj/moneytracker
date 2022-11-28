SELECT *
FROM entity;
--
SELECT *
FROM account;
--
SELECT *
FROM "transaction";
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
FROM "transaction" t
    INNER JOIN operation op ON t.operation_id = op.id
    INNER JOIN account AS fa ON t.from_id = fa.id
    INNER JOIN accounts AS ta ON t.to_id = ta.id
WHERE from_owner_id = 1
    OR to_owner_id = 1
ORDER BY op.timestamp DESC,
    op.id,
    t.id;
-- Get Transactions by Account
SELECT id,
    timestamp,
    from_id,
    to_id,
    amount,
    operation_id
FROM transactions
WHERE from_id = 1
    OR to_id = 1
ORDER BY timestamp DESC;
-- Get Transactions By Entity
SELECT t.*,
    fa.display_name AS from_name,
    ta.display_name AS to_name,
    fa.owner_id AS from_owner_id,
    ta.owner_id AS to_owner_id
FROM transactions AS t
    INNER JOIN accounts AS fa ON t.from_id = fa.id
    INNER JOIN accounts AS ta ON t.to_id = ta.id
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
FROM transactions AS t
    INNER JOIN accounts AS fa ON t.from_id = fa.id
    INNER JOIN accounts AS ta ON t.to_id = ta.id
WHERE from_owner_id = 1
    OR to_owner_id = 1;
    --
    INSERT INTO accounts (name,display_name,owner_id,is_system,is_credit,is_world) VALUES("prova","ssd",9,TRUE,FALSE,FALSE)
     RETURNING *;

--
SELECT *
FROM operation
     INNER JOIN `transaction` ON (`transaction`.id = operation.id)
WHERE operation.id = 1;
--
SELECT * FROM `transaction` WHERE operation_id = 1;