

BEGIN TRANSACTION;

/*Users*/

CREATE TABLE users (
	user_id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	user_name	TEXT NOT NULL
);

INSERT INTO users(user_id,user_name) VALUES (1,'Arianna');
INSERT INTO users(user_id,user_name) VALUES (2,'Marco');

UPDATE expenses SET who=1 WHERE expenses.who = 'A';
UPDATE expenses SET who=2 WHERE expenses.who = 'M';


/*Types*/

CREATE TABLE types (
	type_id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	type_name	TEXT NOT NULL
);

INSERT INTO types(type_id,type_name) VALUES (0,'Expense');
INSERT INTO types(type_id,type_name) VALUES (1,'Transfer');
INSERT INTO types(type_id,type_name) VALUES (2,'Income');
/* See https://github.com/emarj/moneytracker/issues/4
UPDATE transactions SET type = 2 WHERE transactions.amount < 0*/

/*PaymentMethods*/
DROP TABLE IF EXISTS paymentmethods;

CREATE TABLE "paymentmethods" (
	`pm_id`	INTEGER NOT NULL,
	`pm_name`	TEXT NOT NULL,
	PRIMARY KEY(pm_id)
);
INSERT INTO `paymentmethods` VALUES
 (1,'Cash'),
 (2,'Debit Card'),
 (3,'Credit Card'),
 (4,'Paypal'),
 (5,'Bank Account');


/*Categories*/

DROP TABLE IF EXISTS categories;
CREATE TABLE "categories" (
	`cat_id`	INTEGER NOT NULL,
	`cat_name`	TEXT NOT NULL,
	PRIMARY KEY(cat_id)
);
INSERT INTO `categories`
VALUES (0,'Uncategorized'),
 (1,'Trasferimento'),
 (2,'Spesa'),
 (3,'Ristorante'),
 (4,'Bar / Caffè'),
 (5,'Server / Cloud'),
 (6,'Pasti'),
 (7,'Salute'),
 (8,'Sport'),
 (9,'Uni'),
 (10,'Svago'),
 (11,'Trasporti'),
 (12,'Assicurazione'),
 (13,'Tasse'),
 (14,'Bollette'),
 (15,'Telefono'),
 (16,'Varie'),
 (17,'Spese Condominiali'),
 (18,'Regali'),
 (19,'Corso Inglese');


/*Dates*/
UPDATE expenses SET date = date || "T00:00:00";

/*Transactions*/

CREATE TABLE transactions (
	uuid	TEXT NOT NULL,
	date_created	TEXT NOT NULL,
	date_modified TEXT NOT NULL,
	date	TEXT NOT NULL,
	type_id	INTEGER NOT NULL,
	user_id	INTEGER NOT NULL,
	amount	NUMERIC NOT NULL,
	description	TEXT NOT NULL,
	pm_id	INTEGER,
	shared	INTEGER NOT NULL, /*Maybe even this should be removed? Count on join*/
	cat_id	INTEGER NOT NULL,
	geolocation TEXT,
	PRIMARY KEY(uuid)
);

INSERT INTO transactions
SELECT uuid,datecreated,datecreated,date,type,who,amount,description,method,shared,category, "" as geolocation
FROM expenses;

UPDATE transactions
SET shared=1
WHERE type_id=1

/*Shares*/

CREATE TABLE shares (
	tx_uuid	TEXT NOT NULL,
	with_id	INTEGER NOT NULL,
	quota	NUMERIC NOT NULL
);

INSERT INTO shares
SELECT uuid, CASE WHEN who = 1 THEN 2 ELSE 1 END, CAST ((CAST(amount AS REAL))*((CAST(quota AS REAL)) / 100) AS NUMERIC)
FROM expenses
WHERE shared=1;

INSERT INTO shares
SELECT uuid, CASE WHEN who = 1 THEN 2 ELSE 1 END, -amount
FROM expenses
WHERE type=1;

DROP TABLE IF EXISTS expenses;

COMMIT;