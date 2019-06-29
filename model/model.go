package model

import (
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

var Users = []string{"A", "M"}

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
	Who         string
	Method      *PaymentMethod
	Amount      decimal.Decimal
	Shared      bool
	ShareQuota  int
	Category    *Category
	InSheet     bool
	Type        int
}

func NewExpenseNoID(
	dateCreated string,
	date string,
	amount decimal.Decimal,
	description string,
	shared string,
	quota int,
	who string,
	inSheet bool,
	typ int,
	methodID int,
	methodName string,
	catID int,
	catName string) (*Expense, error) {
	c := Category{ID: catID, Name: catName}
	pm := PaymentMethod{ID: methodID, Name: methodName}
	e := Expense{Category: &c, Method: &pm, Who: who, Type: typ, InSheet: inSheet}

	e.Amount = amount
	e.ShareQuota = quota
	e.Description = description

	shrd, err := strconv.ParseBool(shared)
	if err != nil {
		return nil, err
	}
	e.Shared = shrd

	dc, err := time.Parse("2006-01-02T15:04:05", dateCreated)
	if err != nil {
		return nil, err
	}
	e.DateCreated = dc

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	e.Date = d

	return &e, nil
}

func NewExpense(id string, dateCreated string,
	date string,
	amount decimal.Decimal,
	description string,
	shared string,
	quota int,
	who string,
	inSheet bool,
	typ int,
	methodID int,
	methodName string,
	catID int,
	catName string) (*Expense, error) {
	e, err := NewExpenseNoID(
		dateCreated,
		date,
		amount,
		description,
		shared,
		quota,
		who,
		inSheet,
		typ,
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
	ExpensesGetNOrderBy(limit int, orderby string) ([]*Expense, error)
	ExpensesGetNOrderByDate(limit int) ([]*Expense, error)
	ExpensesGetNOrderByInserted(limit int) ([]*Expense, error)
	ExpenseGet(uuid.UUID) (*Expense, error)
	ExpenseInsert(*Expense) error
	ExpenseUpdate(*Expense) error
	ExpenseDelete(uuid.UUID) error

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

	GetAmount(User string) (decimal.Decimal, error)

	Close() error
}
