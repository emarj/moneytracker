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
    INSERT INTO account (name,display_name,owner_id,is_system,is_crebit,is_world) VALUES("prova","ssd",9,TRUE,FALSE,FALSE)
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
---
SELECT * FROM category;
---
SELECT category.id AS "category.id",
     category.parent_id AS "category.parent_id",
     category.name AS "category.name",
     parent.id AS "parent.id",
     parent.parent_id AS "parent.parent_id",
     parent.name AS "parent.name"
FROM category
     LEFT JOIN category AS parent ON (parent.id = category.parent_id);
     -----

SELECT *
FROM operation
     LEFT JOIN `transaction` ON (`transaction`.operation_id = operation.id)
     LEFT JOIN account AS `from` ON (`from`.id = `transaction`.from_id)
     LEFT JOIN account AS `to` ON (`to`.id = `transaction`.to_id)
     LEFT JOIN entity AS `from.owner` ON (`from.owner`.id = `from`.owner_id)
     LEFT JOIN entity AS `to.owner` ON (`to.owner`.id = `to`.owner_id)
     LEFT JOIN balance ON (balance.operation_id = operation.id)
     LEFT JOIN entity AS `bal_owner` ON (bal_owner.id = balance.account_id)
WHERE ((`to`.owner_id = 1) OR (`from`.owner_id = 1)) OR (bal_owner.id = 1)
ORDER BY operation.modified_on DESC
LIMIT 1000;
---
SELECT 
     operation.description AS "operation.description",
     `transaction`.id AS "transaction.id",
     `transaction`.timestamp AS "transaction.timestamp",
     `transaction`.from_id AS "transaction.from_id",
     `transaction`.to_id AS "transaction.to_id",
     `transaction`.amount AS "transaction.amount",
      (CASE WHEN (`from.owner`.id IN (
          SELECT entity_share.entity_id AS "entity_share.entity_id"
          FROM entity_share
          WHERE entity_share.user_id = 1
     )) AND (`to.owner`.id IN (
          SELECT entity_share.entity_id AS "entity_share.entity_id"
          FROM entity_share
          WHERE entity_share.user_id = 1
     )) THEN 0 ELSE -1 END) AS "transaction.sign",
     `transaction`.operation_id AS "transaction.operation_id",
     `from`.id AS "from.id",
     `from`.owner_id AS "from.owner_id",
     `from`.name AS "from.name",
     `to`.id AS "to.id",
     `to`.owner_id AS "to.owner_id",
     `to`.name AS "to.name",
     `from.owner`.name AS "from.owner.name",
     `to.owner`.name AS "to.owner.name"
FROM (
          SELECT *
          FROM operation
          ORDER BY operation.modified_on DESC
          LIMIT 1000
     ) AS operation
     LEFT JOIN `transaction` ON (`transaction`.operation_id = operation.id)
     LEFT JOIN account AS `from` ON (`from`.id = `transaction`.from_id)
     LEFT JOIN account AS `to` ON (`to`.id = `transaction`.to_id)
     LEFT JOIN entity AS `from.owner` ON (`from.owner`.id = `from`.owner_id)
     LEFT JOIN entity AS `to.owner` ON (`to.owner`.id = `to`.owner_id)
     LEFT JOIN balance ON (balance.operation_id = operation.id)
     LEFT JOIN account AS account ON (account.id = balance.account_id)
     LEFT JOIN entity AS `account.owner` ON (`account.owner`.id = account.owner_id)
WHERE ((`to`.owner_id IN (
           SELECT entity_share.entity_id AS "entity_share.entity_id"
           FROM entity_share
           WHERE entity_share.user_id = 1
      )) OR (`from`.owner_id IN (
           SELECT entity_share.entity_id AS "entity_share.entity_id"
           FROM entity_share
           WHERE entity_share.user_id = 1
      ))) OR (`account.owner`.id IN (
           SELECT entity_share.entity_id AS "entity_share.entity_id"
           FROM entity_share
           WHERE entity_share.user_id = 1
      ))
ORDER BY operation.modified_on DESC, `transaction`.timestamp DESC, balance.timestamp DESC;
----
SELECT operation.id AS "operation.id",
     operation.created_on AS "operation.created_on",
     operation.modified_on AS "operation.modified_on",
     operation.created_by_id AS "operation.created_by_id",
     operation.description AS "operation.description",
     operation.type_id AS "operation.type_id",
     operation.details AS "operation.details",
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
     `from`.is_group AS "from.is_group",
     `from`.type_id AS "from.type_id",
     `from`.group_id AS "from.group_id",
     `to`.id AS "to.id",
     `to`.owner_id AS "to.owner_id",
     `to`.name AS "to.name",
     `to`.display_name AS "to.display_name",
     `to`.is_system AS "to.is_system",
     `to`.is_group AS "to.is_group",
     `to`.type_id AS "to.type_id",
     `to`.group_id AS "to.group_id",
     `from.owner`.id AS "from.owner.id",
     `from.owner`.name AS "from.owner.name",
     `from.owner`.display_name AS "from.owner.display_name",
     `from.owner`.is_system AS "from.owner.is_system",
     `from.owner`.is_external AS "from.owner.is_external",
     `to.owner`.id AS "to.owner.id",
     `to.owner`.name AS "to.owner.name",
     `to.owner`.display_name AS "to.owner.display_name",
     `to.owner`.is_system AS "to.owner.is_system",
     `to.owner`.is_external AS "to.owner.is_external",
     balance.timestamp AS "balance.timestamp",
     balance.account_id AS "balance.account_id",
     balance.value AS "balance.value",
     balance.delta AS "balance.delta",
     balance.comment AS "balance.comment",
     balance.is_computed AS "balance.is_computed",
     balance.operation_id AS "balance.operation_id",
     account.id AS "account.id",
     account.owner_id AS "account.owner_id",
     account.name AS "account.name",
     account.display_name AS "account.display_name",
     account.is_system AS "account.is_system",
     account.is_group AS "account.is_group",
     account.type_id AS "account.type_id",
     account.group_id AS "account.group_id",
     `account.owner`.id AS "account.owner.id",
     `account.owner`.name AS "account.owner.name",
     `account.owner`.display_name AS "account.owner.display_name",
     `account.owner`.is_system AS "account.owner.is_system",
     `account.owner`.is_external AS "account.owner.is_external",
     `from.owner.user`.user_id AS "from.owner.user.user_id",
     `to.owner.user`.user_id AS "to.owner.user.user_id"
FROM (
          SELECT *
          FROM operation
          ORDER BY operation.modified_on DESC
          LIMIT 1000
     ) AS operation
     LEFT JOIN `transaction` ON (`transaction`.operation_id = operation.id)
     LEFT JOIN account AS `from` ON (`from`.id = `transaction`.from_id)
     LEFT JOIN account AS `to` ON (`to`.id = `transaction`.to_id)
     LEFT JOIN entity AS `from.owner` ON (`from.owner`.id = `from`.owner_id)
     LEFT JOIN entity AS `to.owner` ON (`to.owner`.id = `to`.owner_id)
     LEFT JOIN balance ON (balance.operation_id = operation.id)
     LEFT JOIN account AS account ON (account.id = balance.account_id)
     LEFT JOIN entity AS `account.owner` ON (`account.owner`.id = account.owner_id)
     LEFT JOIN entity_share AS `from.owner.user` ON (`from.owner.user`.entity_id = `from`.owner_id)
     LEFT JOIN entity_share AS `to.owner.user` ON (`to.owner.user`.entity_id = `to`.owner_id)
WHERE ((`to`.owner_id IN (
           SELECT entity_share.entity_id AS "entity_share.entity_id"
           FROM entity_share
           WHERE entity_share.user_id = 1
      )) OR (`from`.owner_id IN (
           SELECT entity_share.entity_id AS "entity_share.entity_id"
           FROM entity_share
           WHERE entity_share.user_id = 1
      ))) OR (`account.owner`.id IN (
           SELECT entity_share.entity_id AS "entity_share.entity_id"
           FROM entity_share
           WHERE entity_share.user_id = 1
      ))
ORDER BY operation.modified_on DESC, `transaction`.timestamp DESC, balance.timestamp DESC;
----
SELECT operation.id AS "operation.id",
     `transaction`.amount AS "transaction.amount",
    1 AS "transaction.direction",
     balance.timestamp AS "balance.timestamp",
     balance.account_id AS "balance.account_id",
     balance.value AS "balance.value",
     `from.owner.user`.user_id AS "from.owner.user.user_id",
     `to.owner.user`.user_id AS "to.owner.user.user_id",
     `account.owner.user`.user_id AS "account.owner.user.user_id"
FROM (
          SELECT *
          FROM operation
          ORDER BY operation.modified_on DESC
          LIMIT 1000
     ) AS operation
     LEFT JOIN `transaction` ON (`transaction`.operation_id = operation.id)
     LEFT JOIN account AS `from` ON (`from`.id = `transaction`.from_id)
     LEFT JOIN account AS `to` ON (`to`.id = `transaction`.to_id)
     LEFT JOIN entity AS `from.owner` ON (`from.owner`.id = `from`.owner_id)
     LEFT JOIN entity AS `to.owner` ON (`to.owner`.id = `to`.owner_id)
     LEFT JOIN balance ON (balance.operation_id = operation.id)
     LEFT JOIN account AS account ON (account.id = balance.account_id)
     LEFT JOIN entity AS `account.owner` ON (`account.owner`.id = account.owner_id)
     LEFT JOIN entity_share AS `from.owner.user` ON (`from.owner.user`.entity_id = `from`.owner_id)
     LEFT JOIN entity_share AS `to.owner.user` ON (`to.owner.user`.entity_id = `to`.owner_id)
     LEFT JOIN entity_share AS `account.owner.user` ON (`account.owner.user`.entity_id = account.owner_id)
WHERE ((`from.owner.user`.user_id = 1) OR (`to.owner.user`.user_id = 1)) OR (`account.owner.user`.user_id = 1)
ORDER BY operation.modified_on DESC, `transaction`.timestamp DESC, balance.timestamp DESC;