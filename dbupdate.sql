
/*Users*/

DROP TABLE IF EXISTS users;

CREATE TABLE users (
	id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name	TEXT NOT NULL
);

INSERT INTO users(id,name) VALUES (1,'Arianna');
INSERT INTO users(id,name) VALUES (2,'Marco');

UPDATE expenses SET who=1 WHERE expenses.who = 'A';
UPDATE expenses SET who=2 WHERE expenses.who = 'M';

/*Types*/

CREATE TABLE types (
	id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name	TEXT NOT NULL
);

INSERT INTO types(id,name) VALUES (0,'Expense');
INSERT INTO types(id,name) VALUES (1,'Transfer');
INSERT INTO types(id,name) VALUES (2,'Income');
/* See https://github.com/emarj/moneytracker/issues/4
UPDATE transactions SET type = 2 WHERE transactions.amount < 0*/

/*Transactions*/

CREATE TABLE transactions (
	uuid	TEXT NOT NULL,
	date_created	TEXT NOT NULL,
	date	TEXT NOT NULL,
	type_id	INTEGER NOT NULL,
	user_id	INTEGER NOT NULL,
	amount	NUMERIC NOT NULL,
	description	TEXT NOT NULL,
	method_id	INTEGER,
	shared	INTEGER NOT NULL, /*Maybe even this should be removed? Count on join*/
	category_id	INTEGER NOT NULL,
	geolocation TEXT,
	PRIMARY KEY(uuid)
);

INSERT INTO transactions
SELECT uuid,datecreated,date,type,who,amount,description,method,shared,category, "" as geolocation
FROM expenses;

/*Sharings*/

CREATE TABLE sharings (
	transaction_UUID	TEXT NOT NULL,
	user_ID	INTEGER NOT NULL,
	shared_quota	NUMERIC NOT NULL
);

INSERT INTO sharings
SELECT uuid, CASE WHEN who = 1 THEN 2 ELSE 1 END, CAST ((CAST(amount AS REAL))*((CAST(quota AS REAL)) / 100) AS NUMERIC)
FROM expenses
WHERE shared=1;



DROP TABLE IF EXISTS expenses;