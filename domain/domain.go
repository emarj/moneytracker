package domain

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID   string `schema:"user_id"`
	Name string `schema:"user_name"`
}

type Account struct {
	OwnersID    []string
	Name        string
	DisplayName string
	Balance     decimal.Decimal
	Immutable   bool
}

func (a *Account) ID() string {
	return strings.Join(a.OwnersID, "&") + ":" + a.Name
}

func (a *Account) AlterBalance(delta decimal.Decimal) {
	if a.Immutable {
		return
	}

	a.Balance = a.Balance.Add(delta)
}

func (a *Account) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          string
		OwnersID    []string
		Name        string
		DisplayName string
		Balance     decimal.Decimal
		Immutable   bool
	}{
		ID:          a.ID(),
		OwnersID:    a.OwnersID,
		Name:        a.Name,
		DisplayName: a.DisplayName,
		Balance:     a.Balance,
		Immutable:   a.Immutable,
	})

}

var WorldAcc = Account{Name: "world", OwnersID: []string{"mt"}, DisplayName: "🌍", Immutable: true}

var Balances map[[2]User]Account

const (
	TypeExpense  = "expense"
	TypeIncome   = "income"
	TypeTransfer = "transfer"
)

// A transaction is always a transfer between 2 accounts
type Transaction struct {
	ID           uuid.UUID
	DateCreated  DateTime
	DateModified DateTime

	Date Date

	Description string
	Notes       string
	Amount      decimal.Decimal

	FromID string
	ToID   string

	Type string
}

/*
type Expense struct {
	Transaction
	GeoLocation string `schema:"geolocation"`
	Receipt     []byte

	Shares []*ExpenseShare
}

func NewExpense(date time.Time, desc string, notes string, amount decimal.Decimal, from *Account, expenses []*ExpenseShare) *Expense {
	return &Expense{
		Transaction: Transaction{
			Date:        Date{date},
			Description: desc,
			Notes:       notes,
			Amount:      amount,
			From:        from,
			To:          &WorldAcc,
		},
		GeoLocation: "",
		Receipt:     []byte{},
		Shares:      expenses,
	}
}

//These are virtual transaction
type ExpenseShare struct {
	WithName    string
	WithUser    *User
	Amount      decimal.Decimal
	TotalAmount decimal.Decimal
	AlreadyPaid bool
}*/
