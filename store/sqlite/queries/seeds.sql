-- Add system entity
INSERT INTO entities (id, name, is_system)
VALUES (0, "_system", TRUE);
-- Add world account
INSERT INTO accounts (
		id,
		entity_id,
		name,
		display_name,
		is_money,
		is_external,
		is_system
	)
VALUES (0, 0, "_world", "World", TRUE, FALSE, TRUE);
--
INSERT INTO entities (id, name, is_system)
VALUES (1, "user_1", FALSE);
--
INSERT INTO entities (id, name, is_system)
VALUES (2, "user_2", FALSE);
--
INSERT INTO accounts (
id,
entity_id,
name,
display_name,
is_money,
is_external,
is_system
)
VALUES (
		1,
		1,
		"account_1",
		"Account 1",
		TRUE,
		FALSE,
		FALSE
	);
--
	INSERT INTO accounts (
id,
entity_id,
name,
display_name,
is_money,
is_external,
is_system
)
VALUES (
		2,
		1,
		"account_2",
		"Account 2",
		TRUE,
		FALSE,
		FALSE
	);
--
	INSERT INTO accounts (
id,
entity_id,
name,
display_name,
is_money,
is_external,
is_system
)
VALUES (
		3,
		2,
		"account_3",
		"Account 3",
		TRUE,
		FALSE,
		FALSE
	);
-- Add balance
INSERT INTO balances (timestamp, account_id, value, computed)
VALUES ("2022-11-12T16:59:35.498Z", 1, 1000, FALSE);
INSERT INTO balances (timestamp, account_id, value, computed)
VALUES ("2022-11-12T16:59:35.498Z", 2, 80, FALSE);
-- Add some transactions
INSERT INTO transactions (timestamp, from_id, to_id, amount)
VALUES ("2022-11-12T17:05:30.498Z", 1, 2, 300);
--
INSERT INTO transactions (timestamp, from_id, to_id, amount)
VALUES ("2022-11-12T17:05:32.498Z", 2, 1, 230);
--
INSERT INTO transactions (timestamp, from_id, to_id, amount)
VALUES ("2022-11-12T17:05:40.498Z", 3, 0, 999);
-- Add balance
INSERT INTO balances (account_id, value, computed)
VALUES (1, 2000, FALSE);
--
INSERT INTO transactions (from_id, to_id, amount)
VALUES (1, 0, 777);
---
INSERT INTO transactions (from_id, to_id, amount)
VALUES (3, 2, 400);