CREATE TABLE info (
	schema_version TEXT NOT NULL
);
INSERT INTO info (schema_version) VALUES ("0.0.1");

/*
CREATE TABLE record (
	id INTEGER PRIMARY KEY,
	created_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	modified_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
);
*/

CREATE TABLE user (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	display_name TEXT NOT NULL,
	is_admin INTEGER NOT NULL DEFAULT FALSE CHECK (is_admin IN (0, 1)),
	password TEXT NOT NULL
);
---
CREATE TABLE entity (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	is_system INTEGER NOT NULL DEFAULT FALSE CHECK (is_system IN (0, 1)),
	is_external INTEGER NOT NULL DEFAULT FALSE CHECK (is_external IN (0, 1))
);
INSERT INTO entity (id,name,is_system,is_external) VALUES (0,"_system",TRUE,FALSE) RETURNING *;
---
CREATE TABLE account_type (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

INSERT INTO account_type (id,name) VALUES (0,"money"),(1,"credit"),(2,"investment")
RETURNING *;

CREATE TABLE account (
	id INTEGER PRIMARY KEY,
	owner_id INTEGER NOT NULL REFERENCES entity(id),
	name TEXT NOT NULL,
	display_name TEXT NOT NULL,
	is_system INTEGER NOT NULL CHECK (is_system IN (0, 1)),
	is_world INTEGER NOT NULL CHECK (is_world IN (0, 1)),
	is_group INTEGER NOT NULL DEFAULT FALSE CHECK (is_group IN (0, 1)),
	type INTEGER NOT NULL DEFAULT 0 REFERENCES account_type(id),
	parent_id INTEGER REFERENCES account(id)
);

INSERT INTO account (id,owner_id,name,display_name,is_system,is_world) VALUES (0,0,"_world","World",TRUE,TRUE);
---
CREATE TABLE balance (
	timestamp TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	account_id INTEGER REFERENCES account(id),
	value TEXT NOT NULL,
	is_computed INTEGER NOT NULL CHECK (is_computed IN (0, 1)),
	operation_id INTEGER REFERENCES operation(id),
	PRIMARY KEY(account_id, timestamp)
);
---
CREATE TABLE 'transaction' (
	id INTEGER PRIMARY KEY,
	timestamp TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	from_id INTEGER NOT NULL,
	to_id INTEGER NOT NULL,
	amount TEXT NOT NULL,
	operation_id INTEGER NOT NULL REFERENCES operation(id)
);
---
CREATE TABLE operation_type (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

INSERT INTO operation_type (id,name) VALUES (0,"other"),(1,"expense"),(2,"income"),(3,"transfer"),(4,"balance")
RETURNING *;

CREATE TABLE operation (
	id INTEGER PRIMARY KEY,
	created_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	modified_on TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'now')),
	created_by_id INTEGER NOT NULL,
	---
	description TEXT NOT NULL,
	type_id INTEGER DEFAULT 0 REFERENCES operation_type(id),
	---
	category_id INTEGER DEFAULT 0 REFERENCES category(id)
);
---
CREATE TABLE category (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);