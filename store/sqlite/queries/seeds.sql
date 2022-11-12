-- Add "system" entity
INSERT INTO "entities" (
		"id",
		"name",
		"is_system"
	)
VALUES (0, "_system", TRUE);
-- Add "world" account
INSERT INTO "accounts" (
		"id",
		"owner_id",
		"name",
		"is_money",
		"is_external",
		"is_system"
	)
VALUES (0, 0, "_world", TRUE, FALSE, TRUE);
INSERT INTO "accounts" (
		"id",
		"owner_id",
		"name",
		"is_money",
		"is_external",
		"is_system"
	)
VALUES (1, 0, "acc1", TRUE, FALSE, FALSE);
-- Add balance
INSERT INTO "balances" ("timestamp","account_id", "value", "computed")
VALUES ("2022-11-12T16:59:35.498Z",1, 1000, 0);
-- Add some transactions
INSERT INTO "transactions" ("timestamp","from_id", "to_id", "amount")
VALUES ("2022-11-12T17:05:30.498Z",1, 2, 300);
INSERT INTO "transactions" ("timestamp","from_id", "to_id", "amount")
VALUES 	("2022-11-12T17:05:32.498Z",2, 1, 230);
INSERT INTO "transactions" ("timestamp","from_id", "to_id", "amount")
VALUES 	("2022-11-12T17:05:34.498Z",3, 1, 500);
-- Add balance
INSERT INTO "balances" ("account_id", "value", "computed")
VALUES (1, 2000, 0);
-- Add some more transactions
INSERT INTO "transactions" ("from_id", "to_id", "amount")
VALUES 	(3, 1, 500);