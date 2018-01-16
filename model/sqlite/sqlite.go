package sqlite

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	//sqlite driver
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"ronche.se/expensetracker/model"
)

var schema = [...]string{

	`CREATE TABLE IF NOT EXISTS expenses (
	uuid	TEXT NOT NULL,
	datecreated	TEXT NOT NULL,
	date	TEXT NOT NULL,
	who		INTEGER NOT NULL,
	amount	INTEGER NOT NULL,
	method	INTEGER,
	description	TEXT NOT NULL,
	category	INTEGER NOT NULL,
	shared	TEXT NOT NULL,
	quota	INTEGER,
	PRIMARY KEY(uuid)
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
)`,

	`CREATE TABLE IF NOT EXISTS users (
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

func (s *sqlite) ExpensesGetN(limit int) ([]*model.Expense, error) {
	rows, err := s.db.Query(
		`SELECT expenses.uuid,
		expenses.datecreated,
		expenses.date,
		expenses.amount,
		expenses.description,
		expenses.shared,
		expenses.quota,
		users.id,
		users.name,
		paymentmethods.id,
		paymentmethods.name,
		categories.id,
		categories.name
		 FROM expenses,users,paymentmethods,categories
		WHERE expenses.who=users.id AND expenses.method=paymentmethods.id AND expenses.category=categories.id
		LIMIT ` + strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var es []*model.Expense

	var (
		id          string
		dateCreated string
		date        string
		amount      int
		description string
		shared      string
		quota       int
		whoID       int
		whoName     string
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
			&quota,
			&whoID,
			&whoName,
			&methodID,
			&methodName,
			&catID,
			&catName,
		)
		if err != nil {
			return nil, err
		}

		e, err := model.ParseExpense(id, dateCreated,
			date,
			amount,
			description,
			shared,
			quota,
			whoID,
			whoName,
			methodID,
			methodName,
			catID,
			catName)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return es, nil
}

func (s *sqlite) ExpenseInsert(e *model.Expense) (*model.Expense, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	e.UUID = id
	e.DateCreated = time.Now().Local()

	stmt, err := s.db.Prepare(
		`INSERT INTO expenses(
			uuid,
			datecreated,
			date,
			who,
			amount,
			method,
			description,
			category,
			shared,
			quota
		) VALUES(?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(
		e.UUID.String(),
		e.DateCreated.Format("2006-02-01"),
		e.Date.Format("2006-02-01"),
		e.Who.ID,
		e.Amount,
		e.Method.ID,
		e.Description,
		e.Category.ID,
		strconv.FormatBool(e.Shared),
		e.ShareQuota,
	)

	if err != nil {
		return nil, err
	}

	return e, nil
}
func (s *sqlite) ExpenseGet(id uuid.UUID) (*model.Expense, error) {

	return &model.Expense{}, nil
}

func (s *sqlite) ExpenseUpdate(expense *model.Expense) (*model.Expense, error) {
	return &model.Expense{}, nil
}

func (s *sqlite) ExpenseDelete(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM expenses WHERE uuid=?", id)
	if err != nil {
		return err
	}
	return nil
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
