package model

import (
	"errors"
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

type Sharing map[int]decimal.Decimal

/*func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}*/

func timeParse(format string, v interface{}) (time.Time, error) {
	var vt time.Time
	var s string
	switch z := v.(type) {
	case []uint8:
		s = string(z)
	case string:
		s = z
	default:
		return vt, errors.New("cannot convert to string")
	}
	vt, err := time.Parse(format, s)
	if err != nil {
		return vt, err
	}
	return vt, nil
}

type DateTime time.Time

func (t *DateTime) Scan(v interface{}) error {
	vt, err := timeParse("2006-01-02T15:04:05", v)
	if err != nil {
		return err
	}
	*t = DateTime(vt)
	return nil
}

func (t *DateTime) Format(format string) string {
	return time.Time(*t).Format(format)

}

type Date time.Time

func (t *Date) Scan(v interface{}) error {
	vt, err := timeParse("2006-01-02", v)
	if err != nil {
		return err
	}
	*t = Date(vt)
	return nil
}

func (t *Date) Format(format string) string {
	return time.Time(*t).Format(format)

}

type Transaction struct {
	UUID        uuid.UUID
	DateCreated DateTime
	Date        Date
	Description string
	Amount      decimal.Decimal
	Shared      bool
	SharedQuota decimal.Decimal
	GeoLocation string `json:"geolocation"`

	User
	PaymentMethod
	Category
	Type
}

func NewTransactionNoID(
	dateCreated string,
	date string,
	amount decimal.Decimal,
	description string,
	shared string,
	sharedQuota decimal.Decimal,
	geoLoc string,
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
	t := Transaction{Category: c, PaymentMethod: pm, User: u, Type: tp}

	t.Amount = amount
	t.SharedQuota = sharedQuota
	t.Description = description
	t.GeoLocation = geoLoc

	shrd, err := strconv.ParseBool(shared)
	if err != nil {
		return nil, err
	}
	t.Shared = shrd

	dc, err := time.Parse("2006-01-02T15:04:05", dateCreated)
	if err != nil {
		return nil, err
	}
	t.DateCreated = DateTime(dc)

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	t.Date = Date(d)

	return &t, nil
}

func NewTransaction(
	id string,
	dateCreated string,
	date string,
	amount decimal.Decimal,
	description string,
	shared string,
	sharedQuota decimal.Decimal,
	geoLoc string,
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
