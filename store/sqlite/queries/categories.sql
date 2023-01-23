SELECT category.id AS "category.id",
     category.parent_id AS "category.parent_id",
     category.name AS "category.name",
     parent.id AS "parent.id",
     parent.parent_id AS "parent.parent_id",
     parent.name AS "parent.name",
     COALESCE((parent.name || '/') || category.name, category.name) AS "category.full_name"
FROM category
     LEFT JOIN category AS parent ON (parent.id = category.parent_id)
ORDER BY "category.full_name";