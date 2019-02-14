package sqlite

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/iancoleman/strcase"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	//sqlite driver
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"ronche.se/moneytracker/model"
)

var schema = [...]string{

	`CREATE TABLE IF NOT EXISTS transactions (
		uuid	TEXT NOT NULL,
		date_created	TEXT NOT NULL,
		date	TEXT NOT NULL,
		type_id	INTEGER NOT NULL,
		user_id	INTEGER NOT NULL,
		amount	NUMERIC NOT NULL,
		description	TEXT NOT NULL,
		pm_id	INTEGER,
		shared	INTEGER NOT NULL,
		shared_quota	NUMERIC NOT NULL,
		geolocation TEXT,
		cat_id	INTEGER NOT NULL,
		PRIMARY KEY(uuid)
)`,

	`CREATE TABLE IF NOT EXISTS types ( 
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL, 
	PRIMARY KEY(id)
)`,

	`CREATE TABLE IF NOT EXISTS users ( 
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL, 
	PRIMARY KEY(id)
)`,

	`CREATE TABLE IF NOT EXISTS categories ( 
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL, 
	PRIMARY KEY(id)
)`,

	`CREATE TABLE IF NOT EXISTS paymentmethods (
	id	INTEGER NOT NULL,
	name	TEXT NOT NULL,
	PRIMARY KEY(id)
)`}

type sqlite struct {
	db *sqlx.DB
}

func New(path string, create bool) (*sqlite, error) {

	if !create {
		_, err := os.Open(path)
		if err != nil {
			log.Fatalf("impossible to open the db file: %v", err)
		}
	}
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("json", strcase.ToSnake)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if create {
		for i := range schema {
			_, err = db.Exec(schema[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return &sqlite{db}, nil

}

func (s *sqlite) TransactionGet(uid uuid.UUID) (*model.Transaction, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Preparex(
		`SELECT  *
		FROM users,types,paymentmethods,categories,transactions t INNER JOIN shares s ON t.uuid = s.tx_uuid
		WHERE 
						t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id AND
						t.uuid=?
		UNION
		SELECT *, ? AS tx_uuid, 0 AS with_id,0 AS quota
		FROM users,types,paymentmethods,categories,transactions t 
		WHERE 
						t.user_id=users.user_id AND
						t.type_id=types.type_id AND
						t.pm_id=paymentmethods.pm_id AND
						t.cat_id=categories.cat_id AND
						t.uuid=?
		ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(uid.String(), uuid.Nil.String(), uid.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close() //should not be needed if we iterate over all rows

	type Result struct {
		model.Transaction
		model.Share
	}

	var t *model.Transaction
	flagTx := false //A safety flag, checks that only 1 rows is actually (just) a transaction

	for rows.Next() {
		var result Result
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		if result.Share.Parent == uuid.Nil {
			if flagTx {
				return nil, errors.New("multiple rows with Nil UUID in insert query")
			}
			flagTx = true
			t = &result.Transaction
		} else {
			t.Shares = append(t.Shares, &result.Share)
		}

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return t, err
}

func (s *sqlite) TransactionInsert(t *model.Transaction) error {
	id := uuid.NewV4()

	t.UUID = id
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return err
	}
	t.DateCreated.Time = time.Now().In(loc)

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.PrepareNamed(
		`INSERT INTO transactions(
			uuid,
			date_created,
			date,
			user_id,
			amount,
			pm_id,
			description,
			cat_id,
			shared,
			geolocation,
			type_id
		) VALUES(
			:uuid,
			:date_created,
			:date,
			:user_id,
			:amount,
			:pm_id,
			:description,
			:cat_id,
			:shared,
			:geolocation,
			:type_id)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t)

	if err != nil {
		return err
	}

	if t.Shared && (t.Shares != nil) {
		query := `INSERT INTO shares(
			tx_uuid,
			with_id,
			quota) VALUES`

		vals := []interface{}{}
		for _, shr := range t.Shares {
			query += "(?,?,?),"
			vals = append(vals, t.UUID.String(), shr.WithID, shr.Quota)
		}

		query = query[0 : len(query)-1] //Remove last comma

		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(vals...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *sqlite) TransactionUpdate(t *model.Transaction) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.PrepareNamed(
		`UPDATE transactions
		SET			
			date = :date,
			user_id = :user_id,
			amount = :amount,
			pm_id = :pm_id,
			description = :description,
			cat_id = :cat_id,
			shared = :shared,
			geolocation = :geolocation,
			type_id = :type_id
		WHERE uuid = :uuid`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	//If t was shared and now shares are remove we should account for that!
	//Probably get all shares by id, delete them and add them again is the quick and dirty way!
	//OR WITH ADD A SHARE UUID

	_, err = tx.Exec(`DELETE FROM shares WHERE tx_uuid =?`, t.UUID.String())
	if err != nil {
		return err
	}
	if t.Shared {
		if t.Shares == nil {
			return errors.New("Transaction is shared but it has no shares")
		}
		query := `INSERT INTO shares(
				tx_uuid,
				with_id,
				quota) VALUES`

		vals := []interface{}{}
		for _, shr := range t.Shares {
			query += "(?,?,?),"
			vals = append(vals, t.UUID.String(), shr.WithID, shr.Quota)
		}

		query = query[0 : len(query)-1] //Remove last comma

		stmt2, err := tx.Prepare(query)
		if err != nil {
			return err
		}

		_, err = stmt2.Exec(vals...)

		if err != nil {
			return err
		}
	} else if t.Shares != nil {
		return errors.New("Transaction is not shared but it has shares")

	}

	return nil
}

func (s *sqlite) TransactionDelete(id uuid.UUID) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec("DELETE FROM transactions WHERE uuid=?", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM shares WHERE tx_uuid=?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) TransactionsGetNOrderByDate(limit int) ([]*model.Transaction, error) {
	return s.TransactionsGetNOrderBy(limit, "transactions.date DESC, transactions.date_created DESC")
}

func (s *sqlite) TransactionsGetNOrderByInserted(limit int) ([]*model.Transaction, error) {
	return s.TransactionsGetNOrderBy(limit, "transactions.date_created DESC, transactions.date DESC")
}

func (s *sqlite) TransactionsGetNOrderBy(limit int, orderBy string) ([]*model.Transaction, error) {

	stmt, err := s.db.Preparex(
		`SELECT transactions.uuid,
		transactions.date_created,
		transactions.date,
		transactions.amount,
		transactions.description,
		transactions.shared,
		transactions.geolocation,
		users.user_id,
		users.user_name,
		types.type_id,
		types.type_name,
		paymentmethods.pm_id,
		paymentmethods.pm_name,
		categories.cat_id,
		categories.cat_name
		FROM transactions,users,types,paymentmethods,categories
		WHERE 
				transactions.user_id=users.user_id AND
				transactions.type_id=types.type_id AND
				transactions.pm_id=paymentmethods.pm_id AND
				transactions.cat_id=categories.cat_id
		ORDER BY ` + orderBy + `
		LIMIT ?`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close() //should not be needed if we iterate over all rows

	var ts []*model.Transaction

	for rows.Next() {
		var t model.Transaction
		err := rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		ts = append(ts, &t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *sqlite) UsersGetAll() ([]*model.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var us []*model.User

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		us = append(us, &model.User{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (s *sqlite) UserInsert(name string) (*model.User, error) {
	stmt, err := s.db.Prepare("INSERT INTO users(user_name) VALUES(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &model.User{int(lastID), name}, nil
}

func (s *sqlite) TypesGetAll() ([]*model.Type, error) {
	rows, err := s.db.Query("SELECT * FROM types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tps []*model.Type

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		tps = append(tps, &model.Type{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tps, nil
}

func (s *sqlite) TypeInsert(name string) (*model.Type, error) {
	stmt, err := s.db.Prepare("INSERT INTO types(type_name) VALUES(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &model.Type{int(lastID), name}, nil
}

func (s *sqlite) CategoriesGetAll() ([]*model.Category, error) {
	rows, err := s.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*model.Category

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		cats = append(cats, &model.Category{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *sqlite) CategoryInsert(name string) (*model.Category, error) {
	stmt, err := s.db.Prepare("INSERT INTO categories(cat_name) VALUES(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &model.Category{int(lastID), name}, nil
}

func (s *sqlite) PaymentMethodsGetAll() ([]*model.PaymentMethod, error) {
	rows, err := s.db.Query("SELECT * FROM paymentmethods")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pms []*model.PaymentMethod

	var id int
	var name string

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		pms = append(pms, &model.PaymentMethod{id, name})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return pms, nil
}

func (s *sqlite) PaymentMethodInsert(name string) (*model.PaymentMethod, error) {
	stmt, err := s.db.Prepare("INSERT INTO paymentmethods(pm_name) VALUES(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &model.PaymentMethod{int(lastID), name}, nil
}

func (s *sqlite) Close() error {
	return s.db.Close()
}
