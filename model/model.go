package model

import (
	"database/sql/driver"
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

type User struct {
	ID   int    `json:"user_id"`
	Name string `json:"user_name"`
}

type Type struct {
	ID   int    `json:"type_id"`
	Name string `json:"type_name"`
}

type Category struct {
	ID   int    `json:"cat_id"`
	Name string `json:"cat_name"`
}

type PaymentMethod struct {
	ID   int    `json:"pm_id"`
	Name string `json:"pm_name"`
}

type Share struct {
	TxID     uuid.UUID       `json:"tx_uuid"`
	WithID   int             `json:"with_id"`
	WithName string          `json:"with_name"`
	Quota    decimal.Decimal `json:"quota"`
}

type Transaction struct {
	UUID         uuid.UUID
	DateCreated  DateTime
	DateModified DateTime
	Date         DateTime
	Description  string
	Amount       decimal.Decimal
	Shared       bool
	GeoLocation  string `json:"geolocation"`

	Shares []*Share
	User
	PaymentMethod
	Category
	Type
}

func (t Transaction) SharedWith() []interface{} {
	type SharePrint struct {
		UserName string
		Quota    decimal.Decimal
	}
	shares := make([]interface{}, 0, len(t.Shares))
	for _, shr := range t.Shares {
		sp := SharePrint{shr.WithName, shr.Quota}
		shares = append(shares, sp)
	}
	return shares
}

func (t Transaction) SharedWithID1() int {
	if len(t.Shares) < 1 {
		return -1
	}
	return t.Shares[0].WithID
}

func (t Transaction) SharedWithName1() string {
	if len(t.Shares) < 1 {
		return ""
	}
	return t.Shares[0].WithName
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

	TransactionsGetNOrderBy(limit int, orderby string) ([]*Transaction, error)
	TransactionsGetNOrderByDate(limit int) ([]*Transaction, error)
	TransactionsGetNOrderByInserted(limit int) ([]*Transaction, error)
	TransactionsGetNOrderByModified(limit int) ([]*Transaction, error)
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
	PaymentMethodsGetAll() ([]*PaymentMethod, error)
	//PaymentMethodGet(int) (*PaymentMethod, error)
	PaymentMethodInsert(Name string) (*PaymentMethod, error)
	/*PaymentMethodUpdate(*PaymentMethod) (*PaymentMethod, error)
	PaymentMethodDelete(int) error*/

	Close() error
}
