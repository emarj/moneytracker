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
---
SELECT operation.id AS "operation.id",
     operation.created_on AS "operation.created_on",
     operation.modified_on AS "operation.modified_on",
     operation.created_by_id AS "operation.created_by_id",
     operation.description AS "operation.description",
     operation.type_id AS "operation.type_id",
     operation.category_id AS "operation.category_id",
     `transaction`.id AS "transaction.id",
     `transaction`.timestamp AS "transaction.timestamp",
     `transaction`.from_id AS "transaction.from_id",
     `transaction`.to_id AS "transaction.to_id",
     `transaction`.amount AS "transaction.amount",
     `transaction`.operation_id AS "transaction.operation_id",
     `from`.id AS "from.id",
     `from`.owner_id AS "from.owner_id",
     `from`.name AS "from.name",
     `from`.display_name AS "from.display_name",
     `from`.is_system AS "from.is_system",
     `from`.is_world AS "from.is_world",
     `from`.is_group AS "from.is_group",
     `from`.type AS "from.type",
     `from`.parent_id AS "from.parent_id",
     `to`.id AS "to.id",
     `to`.owner_id AS "to.owner_id",
     `to`.name AS "to.name",
     `to`.display_name AS "to.display_name",
     `to`.is_system AS "to.is_system",
     `to`.is_world AS "to.is_world",
     `to`.is_group AS "to.is_group",
     `to`.type AS "to.type",
     `to`.parent_id AS "to.parent_id"
FROM `transaction`
     INNER JOIN operation ON (operation.id = `transaction`.operation_id)
     INNER JOIN account AS `from` ON (`from`.id = `transaction`.from_id)
     INNER JOIN account AS `to` ON (`to`.id = `transaction`.to_id)
WHERE (`to`.owner_id = 1) OR (`from`.owner_id = 1);