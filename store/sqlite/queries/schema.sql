PRAGMA user_version = 1;

/*
CREATE TABLE IF NOT EXISTS record (
	id INTEGER PRIMARY KEY,
	created_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	modified_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
);
*/

CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	display_name TEXT NOT NULL,
	is_admin INTEGER NOT NULL DEFAULT FALSE CHECK (is_admin IN (0, 1)),
	password TEXT NOT NULL
);
---
CREATE TABLE IF NOT EXISTS entity (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	is_system INTEGER NOT NULL DEFAULT FALSE CHECK (is_system IN (0, 1)),
	is_external INTEGER NOT NULL DEFAULT FALSE CHECK (is_external IN (0, 1))
);
---
CREATE TABLE IF NOT EXISTS account_type (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS account (
	id INTEGER PRIMARY KEY,
	owner_id INTEGER NOT NULL REFERENCES entity(id),
	name TEXT NOT NULL,
	display_name TEXT NOT NULL,
	is_system INTEGER NOT NULL DEFAULT FALSE CHECK (is_system IN (0, 1)),
	is_world INTEGER NOT NULL DEFAULT FALSE CHECK (is_world IN (0, 1)),
	is_group INTEGER NOT NULL DEFAULT FALSE CHECK (is_group IN (0, 1)),
	type_id INTEGER NOT NULL ON CONFLICT REPLACE DEFAULT 0 REFERENCES account_type(id),
	group_id INTEGER REFERENCES account(id)
);
---
CREATE TABLE IF NOT EXISTS balance (
	timestamp TEXT NOT NULL ON CONFLICT REPLACE DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	account_id INTEGER REFERENCES account(id),
	value TEXT NOT NULL,
	delta TEXT,
	comment TEXT,
	is_computed INTEGER NOT NULL CHECK (is_computed IN (0, 1)),
	operation_id INTEGER REFERENCES operation(id),
	PRIMARY KEY(account_id, timestamp)
);
---
CREATE TABLE IF NOT EXISTS 'transaction' (
	id INTEGER PRIMARY KEY,
	timestamp TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	from_id INTEGER NOT NULL,
	to_id INTEGER NOT NULL,
	amount TEXT NOT NULL,
	operation_id INTEGER NOT NULL REFERENCES operation(id)
);
---
CREATE TABLE IF NOT EXISTS operation_type (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS operation (
	id INTEGER PRIMARY KEY,
	created_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	modified_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	created_by_id INTEGER NOT NULL,
	---
	description TEXT NOT NULL,
	type_id INTEGER DEFAULT 0 REFERENCES operation_type(id),
	---
	details TEXT,
	category_id INTEGER DEFAULT 0 REFERENCES category(id)
);
---
CREATE TABLE IF NOT EXISTS category (
	id INTEGER PRIMARY KEY,
	parent_id INTEGER,
	name TEXT NOT NULL
);
---

--- Insert Types and Base Categories
INSERT INTO account_type (id,name) VALUES (0,"money"),(1,"credit"),(2,"investment") ON CONFLICT DO NOTHING;
INSERT INTO operation_type (id,name) VALUES (0,"other"),(1,"balance"),(2,"expense"),(3,"income"),(4,"transfer") ON CONFLICT DO NOTHING;
INSERT INTO category (id,name) VALUES (0,"Uncategorized") ON CONFLICT DO NOTHING;

--- Insert System Entity and Account
INSERT INTO entity (id,name,is_system,is_external) VALUES (0,"_system",TRUE,FALSE) ON CONFLICT DO NOTHING;
INSERT INTO account (id,owner_id,name,display_name,is_system,is_world) VALUES (0,0,"_world","World",TRUE,TRUE) ON CONFLICT DO NOTHING;
