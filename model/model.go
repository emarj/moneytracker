package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type DateTime struct{ time.Time }

func (t *DateTime) Scan(v interface{}) error {

	var s string
	switch z := v.(type) {
	case []uint8:
		s = string(z)
	case string:
		s = z
	default:
		return errors.New("cannot convert time to string")
	}

	vt, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	t.Time = vt
	return nil
}

func (t DateTime) Value() (driver.Value, error) {
	return driver.Value(t.Format("2006-01-02T15:04:05")), nil
}

func (t *DateTime) UnmarshalJSON(json []byte) error {
	str := string(json[1 : len(json)-1])
	vt, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		vt, err = time.Parse("2006-01-02", str)
		if err != nil {
			return err
		}
	}
	t.Time = vt
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format("2006-01-02T15:04:05"))
}

type Date struct{ time.Time }

func (t *Date) Scan(v interface{}) error {

	var s string
	switch z := v.(type) {
	case []uint8:
		s = string(z)
	case string:
		s = z
	default:
		return errors.New("cannot convert time to string")
	}

	vt, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	t.Time = vt
	return nil
}

func (t Date) Value() (driver.Value, error) {
	return driver.Value(t.Format("2006-01-02")), nil
}

func (t *Date) UnmarshalJSON(json []byte) error {
	str := string(json[1 : len(json)-1])
	vt, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	t.Time = vt
	return nil
}

func (t Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format("2006-01-02"))
}

type User struct {
	ID   int    `schema:"user_id"`
	Name string `schema:"user_name"`
}

type Type struct {
	ID   int    `schema:"type_id"`
	Name string `schema:"type_name"`
}

type Category struct {
	ID   int    `schema:"cat_id"`
	Name string `schema:"cat_name"`
}

type Method struct {
	ID   int    `schema:"pm_id"`
	Name string `schema:"pm_name"`
}

type Share struct {
	TxID     uuid.UUID       `schema:"tx_uuid"`
	WithID   int             `schema:"with_id"`
	WithName string          `schema:"with_name"`
	Quota    decimal.Decimal `schema:"quota"`
}

type Transaction struct {
	UUID         uuid.UUID
	DateCreated  DateTime
	DateModified DateTime
	Date         Date
	Description  string
	Amount       decimal.Decimal
	Shared       bool
	GeoLocation  string `schema:"geolocation"`

	Shares []*Share

	User     `json:"User"`
	Method   `json:"Method"`
	Category `json:"Category"`
	Type     `json:"Type"`
}

func (t Transaction) SharedQuota() decimal.Decimal {
	var total decimal.Decimal
	for _, shr := range t.Shares {
		total = total.Add(shr.Quota)
	}
	return total
}

type Service interface {

	/*Transactions*/
	TransactionGet(uuid.UUID) (*Transaction, error)
	TransactionInsert(*Transaction) error
	TransactionUpdate(*Transaction) error
	TransactionDelete(uuid.UUID) error

	TransactionsGetNOrderBy(limit int, offset int, orderby string) ([]*Transaction, error)
	TransactionsGetNOrderByDate(limit int, offset int) ([]*Transaction, error)
	TransactionsGetNOrderByInserted(limit int, offset int) ([]*Transaction, error)
	TransactionsGetNOrderByModified(limit int, offset int) ([]*Transaction, error)
	TransactionsGetNByUser(id int, limit int) ([]*Transaction, error)

	TransactionsGetBalance(userID int) (decimal.Decimal, error)
	TransactionsGetCredit(userID1 int, userID2 int) (decimal.Decimal, error)

	/*Types*/
	TypesGetAll() ([]*Type, error)

	/*Users*/
	UsersGetAll() ([]*User, error)

	/*Shares*/

	/*Categories*/
	CategoriesGetAll() ([]*Category, error)
	/*CategoryGet(int) (*Category, error)*/
	CategoryInsert(Name string) (*Category, error)
	/*CategoryUpdate(*Category) (*Category, error)
	CategoryDelete(int,int) error*/

	/*PaymentsMethods*/
	PaymentMethodsGetAll() ([]*Method, error)
	//PaymentMethodGet(int) (*Method, error)
	PaymentMethodInsert(Name string) (*Method, error)
	/*PaymentMethodUpdate(*Method) (*Method, error)
	PaymentMethodDelete(int) error*/

	Close() error
}
