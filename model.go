package moneytracker

import (
	"github.com/shopspring/decimal"
)

type Entity struct {
	ID     string
	Name   string
	System bool
}

type Account struct {
	ID       int
	Name     string
	Owner    string          // TODO: Allow for shared accounts
	Money    bool            // Is it money or assets?
	External bool            // Allow for direct balance manipulation
	System   bool            // Is this a system account
	Balance  decimal.Decimal `json:"omitempty"` // This would be computed from the transactions
}

type Balance struct {
	AccountID int
	Timestamp DateTime
	Value     decimal.Decimal
	Computed  bool
	Notes     string
}

type Transaction struct {
	ID          int
	Timestamp   DateTime
	From        string
	To          string
	Amount      decimal.Decimal
	OperationID string
}

type Operation struct {
	ID           string
	DateCreated  DateTime
	DateModified DateTime
	CreatedBy    string
	////////
	Timestamp   DateTime
	Description string
	Movements   []string
	Category    Category
}

type Category int

const (
	CategoryExpense = iota
	CategoryTransfer
	CategorySharedExpense
	CategoryGift
)
