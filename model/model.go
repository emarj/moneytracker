package model

import (
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID   int
	Name string
}

type Category struct {
	ID   int
	Name string
}

type PaymentMethod struct {
	ID   int
	Name string
}

type Expense struct {
	UUID        uuid.UUID
	DateCreated time.Time
	Date        time.Time
	Description string
	Who         *User
	Method      *PaymentMethod
	Amount      int
	Shared      bool
	ShareQuota  int
	Category    *Category
}

func ParseExpenseNoID(
	dateCreated string,
	date string,
	amount int,
	description string,
	shared string,
	quota int,
	whoID int,
	whoName string,
	methodID int,
	methodName string,
	catID int,
	catName string) (*Expense, error) {
	c := Category{ID: catID, Name: catName}
	pm := PaymentMethod{ID: methodID, Name: methodName}
	u := User{ID: whoID, Name: whoName}
	e := Expense{Category: &c, Method: &pm, Who: &u}

	e.Amount = amount
	e.ShareQuota = quota
	e.Description = description

	shrd, err := strconv.ParseBool(shared)
	if err != nil {
		return nil, err
	}
	e.Shared = shrd

	dc, err := time.Parse("2006-02-01", dateCreated)
	if err != nil {
		return nil, err
	}
	e.DateCreated = dc

	d, err := time.Parse("2006-02-01", date)
	if err != nil {
		return nil, err
	}
	e.Date = d

	return &e, nil
}

func ParseExpense(id string,
	dateCreated string,
	date string,
	amount int,
	description string,
	shared string,
	quota int,
	whoID int,
	whoName string,
	methodID int,
	methodName string,
	catID int,
	catName string) (*Expense, error) {
	e, err := ParseExpenseNoID(
		dateCreated,
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

	uID, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}
	e.UUID = uID
	return e, nil
}

type Service interface {
	ExpensesGetN(limit int) ([]*Expense, error)
	ExpenseGet(uuid.UUID) (*Expense, error)
	ExpenseInsert(*Expense) (*Expense, error)
	ExpenseUpdate(*Expense) (*Expense, error)
	ExpenseDelete(uuid.UUID) error

	UsersGetAll() ([]*User, error)
	//UserGet(int) (*User, error)
	UserInsert(Name string) (*User, error)
	/*UserUpdate(*User) (*User, error)
	UserDelete(int) error*/

	CategoriesGetAll() ([]*Category, error)
	/*CategoryGet(int) (*Category, error)*/
	CategoryInsert(Name string) (*Category, error)
	/*CategoryUpdate(*Category) (*Category, error)
	CategoryDelete(int,int) error*/

	PaymentMethodsGetAll() ([]*PaymentMethod, error)
	//PaymentMethodGet(int) (*PaymentMethod, error)
	PaymentMethodInsert(Name string) (*PaymentMethod, error)
	/*PaymentMethodUpdate(*PaymentMethod) (*PaymentMethod, error)
	PaymentMethodDelete(int) error*/

	Close() error
}
