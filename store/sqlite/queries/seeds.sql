-- Add "system" entity
INSERT INTO "entities" (
		"id",
		"name",
		"is_external",
		"is_system"
	)
VALUES (0, "_system", 0, 1);
-- Add "world" account
INSERT INTO "accounts" (
		"id",
		"owner_id",
		"name",
		"is_money",
		"is_external",
		"is_system"
	)
VALUES (0, 0, "_world", 1, 0, 1);
INSERT INTO "accounts" (
		"id",
		"owner_id",
		"name",
		"is_money",
		"is_external",
		"is_system"
	)
VALUES (1, 0, "acc1", 1, 0, 0);
-- Add balance
INSERT INTO "balances" ("account_id", "value", "computed")
VALUES (1, 1000, 0);
-- Add some transactions
INSERT INTO "transactions" ("from_id", "to_id", "amount")
VALUES (1, 2, 300);
INSERT INTO "transactions" ("from_id", "to_id", "amount")
VALUES 	(2, 1, 230);
INSERT INTO "transactions" ("from_id", "to_id", "amount")
VALUES 	(3, 1, 500);
-- Add balance
INSERT INTO "balances" ("account_id", "value", "computed")
VALUES (1, 2000, 0);
-- Add some more transactions
INSERT INTO "transactions" ("from_id", "to_id", "amount")
VALUES 	(3, 1, 500);