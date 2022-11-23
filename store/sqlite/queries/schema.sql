CREATE TABLE entities (
	id INTEGER,
	name TEXT NOT NULL,
	is_system INTEGER NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE accounts (
	id INTEGER,
	created TEXT DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	entity_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	display_name TEXT NOT NULL,
	is_money INTEGER NOT NULL,
	is_external INTEGER NOT NULL,
	is_system INTEGER NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE balances (
	timestamp TEXT DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	account_id INTEGER NOT NULL,
	value TEXT NOT NULL,
	PRIMARY KEY(account_id, timestamp)
);
CREATE TABLE transactions (
	id INTEGER,
	from_id INTEGER NOT NULL,
	to_id INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	operation_id INTEGER,
	PRIMARY KEY(id)
);
CREATE TABLE operations (
	id INTEGER,
	timestamp TEXT DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	created_by_id INTEGER NOT NULL,
	description TEXT NOT NULL,
	category_id INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY(id)
);