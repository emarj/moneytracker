SELECT *
FROM entities;
--
SELECT *
FROM transactions;
-- Get Transitions By Entity
SELECT t.*,
    fa.display_name AS from_name,
    ta.display_name AS to_name,
    fa.entity_id AS from_entity_id,
    ta.entity_id AS to_entity_id
FROM transactions AS t
    INNER JOIN accounts AS fa ON t.from_id = fa.id
    INNER JOIN accounts AS ta ON t.to_id = ta.id
WHERE from_entity_id = 1
    OR to_entity_id = 1;
---
SELECT t.*,
    fa.display_name AS from_name,
    ta.display_name AS to_name,
    fa.entity_id AS from_entity_id,
    ta.entity_id AS to_entity_id,
    CASE
        WHEN (
            fa.entity_id = 1
            AND ta.entity_id = 1
        ) THEN '-'
        WHEN fa.entity_id = 1 THEN '->'
        WHEN ta.entity_id = 1 THEN '<-'
    END AS direction
FROM transactions AS t
    INNER JOIN accounts AS fa ON t.from_id = fa.id
    INNER JOIN accounts AS ta ON t.to_id = ta.id
WHERE from_entity_id = 1
    OR to_entity_id = 1;