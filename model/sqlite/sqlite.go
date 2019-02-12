package sqlite

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	//sqlite driver
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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
		method_id	INTEGER,
		shared	INTEGER NOT NULL,
		shared_quota	NUMERIC NOT NULL,
		geolocation TEXT,
		category_id	INTEGER NOT NULL,
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
	db *sql.DB
}

func New(path string, createSchema bool) (*sqlite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if createSchema {
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

	var (
		id          string
		dateCreated string
		date        string
		amount      decimal.Decimal
		description string
		shared      string
		sharedQuota decimal.Decimal
		geoLoc      string
		userID      int
		userName    string
		typeID      int
		typeName    string
		methodID    int
		methodName  string
		catID       int
		catName     string
	)

	err := s.db.QueryRow(
		`SELECT
		transactions.uuid,
		transactions.date_created,
		transactions.date,
		transactions.amount,
		transactions.description,
		transactions.shared,
		transactions.shared_quota,
		transactions.geolocation,
		users.id,
		users.name,
		types.id,
		types.name,
		paymentmethods.id,
		paymentmethods.name,
		categories.id,
		categories.name
		FROM transactions,users,types,paymentmethods,categories
		WHERE 
				transactions.user_id=users.id AND
				transactions.type_id=types.id AND
				transactions.method_id=paymentmethods.id AND
				transactions.category_id=categories.id AND
				transactions.uuid ='`+uid.String()+"'").Scan(
		&id,
		&dateCreated,
		&date,
		&amount,
		&description,
		&shared,
		&sharedQuota,
		&geoLoc,
		&userID,
		&userName,
		&typeID,
		&typeName,
		&methodID,
		&methodName,
		&catID,
		&catName,
	)
	if err != nil {
		return nil, err
	}

	t, err := model.NewTransaction(
		id,
		dateCreated,
		date,
		amount,
		description,
		shared,
		sharedQuota,
		geoLoc,
		userID,
		userName,
		typeID,
		typeName,
		methodID,
		methodName,
		catID,
		catName)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *sqlite) TransactionInsert(t *model.Transaction) error {
	id := uuid.NewV4()

	t.UUID = id
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return err
	}
	t.DateCreated = time.Now().In(loc)

	stmt, err := s.db.Prepare(
		`INSERT INTO transactions(
			uuid,
			date_created,
			date,
			user_id,
			amount,
			method_id,
			description,
			category_id,
			shared,
			shared_quota,
			geolocation,
			type_id
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		t.UUID.String(),
		t.DateCreated.Format("2006-01-02T15:04:05"),
		t.Date.Format("2006-01-02"),
		t.User.ID,
		t.Amount,
		t.Method.ID,
		t.Description,
		t.Category.ID,
		t.Shared,
		t.SharedQuota,
		t.GeoLocation,
		t.Type.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) TransactionUpdate(t *model.Transaction) error {
	stmt, err := s.db.Prepare(
		`UPDATE transactions
		SET
			date = ?,
			user_id = ?,
			amount = ?,
			method_id = ?,
			description = ?,
			category_id = ?,
			shared = ?,
			shared_quota = ?,
			geolocation = ?,
			type_id = ?
		WHERE uuid = '` + t.UUID.String() + `'`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		t.Date.Format("2006-01-02"),
		t.User.ID,
		t.Amount,
		t.Method.ID,
		t.Description,
		t.Category.ID,
		t.Shared,
		t.SharedQuota,
		t.GeoLocation,
		t.Type.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) TransactionDelete(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM transactions WHERE uuid=?", id)
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

	rows, err := s.db.Query(
		`SELECT transactions.uuid,
		transactions.date_created,
		transactions.date,
		transactions.amount,
		transactions.description,
		transactions.shared,
		transactions.shared_quota,
		transactions.geolocation,
		users.id,
		users.name,
		types.id,
		types.name,
		paymentmethods.id,
		paymentmethods.name,
		categories.id,
		categories.name
		FROM transactions,users,types,paymentmethods,categories
		WHERE 
				transactions.user_id=users.id AND
				transactions.type_id=types.id AND
				transactions.method_id=paymentmethods.id AND
				transactions.category_id=categories.id
		ORDER BY ` + orderBy + `
		LIMIT ` + strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []*model.Transaction

	var (
		id          string
		dateCreated string
		date        string
		amount      decimal.Decimal
		description string
		shared      string
		sharedQuota decimal.Decimal
		geoLoc      string
		userID      int
		userName    string
		typeID      int
		typeName    string
		methodID    int
		methodName  string
		catID       int
		catName     string
	)

	for rows.Next() {
		err := rows.Scan(
			&id,
			&dateCreated,
			&date,
			&amount,
			&description,
			&shared,
			&sharedQuota,
			&geoLoc,
			&userID,
			&userName,
			&typeID,
			&typeName,
			&methodID,
			&methodName,
			&catID,
			&catName,
		)
		if err != nil {
			return nil, err
		}

		t, err := model.NewTransaction(
			id,
			dateCreated,
			date,
			amount,
			description,
			shared,
			sharedQuota,
			geoLoc,
			userID,
			userName,
			typeID,
			typeName,
			methodID,
			methodName,
			catID,
			catName)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
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
	stmt, err := s.db.Prepare("INSERT INTO users(name) VALUES(?)")
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
	stmt, err := s.db.Prepare("INSERT INTO types(name) VALUES(?)")
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
	stmt, err := s.db.Prepare("INSERT INTO categories(name) VALUES(?)")
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
	stmt, err := s.db.Prepare("INSERT INTO paymentmethods(name) VALUES(?)")
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
