package model

import (
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID   int
	Name string
}

type Type struct {
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

type Transaction struct {
	UUID        uuid.UUID
	DateCreated time.Time
	Date        time.Time
	Description string
	Amount      decimal.Decimal
	//AmountShared decimal.Decimal
	//SharedWith *User Maybe Wallets?
	Shared     bool
	ShareQuota int

	User     *User
	Method   *PaymentMethod
	Category *Category
	Type     *Type
}

func NewTransactionNoID(
	dateCreated string,
	date string,
	amount decimal.Decimal,
	description string,
	shared string,
	quota int,
	userID int,
	userName string,
	typeID int,
	typeName string,
	methodID int,
	methodName string,
	catID int,
	catName string) (*Transaction, error) {

	u := User{ID: userID, Name: userName}
	tp := Type{ID: typeID, Name: typeName}
	c := Category{ID: catID, Name: catName}
	pm := PaymentMethod{ID: methodID, Name: methodName}
	t := Transaction{Category: &c, Method: &pm, User: &u, Type: &tp}

	t.Amount = amount
	t.ShareQuota = quota
	t.Description = description

	shrd, err := strconv.ParseBool(shared)
	if err != nil {
		return nil, err
	}
	t.Shared = shrd

	dc, err := time.Parse("2006-01-02T15:04:05", dateCreated)
	if err != nil {
		return nil, err
	}
	t.DateCreated = dc

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	t.Date = d

	return &t, nil
}

func NewTransaction(
	id string,
	dateCreated string,
	date string,
	amount decimal.Decimal,
	description string,
	shared string,
	quota int,
	userID int,
	userName string,
	typeID int,
	typeName string,
	methodID int,
	methodName string,
	catID int,
	catName string) (*Transaction, error) {
	t, err := NewTransactionNoID(
		dateCreated,
		date,
		amount,
		description,
		shared,
		quota,
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

	uID, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}
	t.UUID = uID
	return t, nil
}

type Service interface {

	/*Transactions*/
	TransactionGet(uuid.UUID) (*Transaction, error)
	TransactionInsert(*Transaction) error
	TransactionUpdate(*Transaction) error
	TransactionDelete(uuid.UUID) error

	TransactionsGetNOrderBy(limit int, orderby string) ([]*Transaction, error)
	TransactionsGetNOrderByDate(limit int) ([]*Transaction, error)
	TransactionsGetNOrderByInserted(limit int) ([]*Transaction, error)

	/*Types*/
	TypesGetAll() ([]*Type, error)

	/*Users*/
	UsersGetAll() ([]*User, error)

	/*Categories*/
	CategoriesGetAll() ([]*Category, error)
	/*CategoryGet(int) (*Category, error)*/
	CategoryInsert(Name string) (*Category, error)
	/*CategoryUpdate(*Category) (*Category, error)
	CategoryDelete(int,int) error*/

	/*PaymentsMethods*/
	PaymentMethodsGetAll() ([]*PaymentMethod, error)
	//PaymentMethodGet(int) (*PaymentMethod, error)
	PaymentMethodInsert(Name string) (*PaymentMethod, error)
	/*PaymentMethodUpdate(*PaymentMethod) (*PaymentMethod, error)
	PaymentMethodDelete(int) error*/

	Close() error
}
