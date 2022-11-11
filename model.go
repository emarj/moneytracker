package moneytracker

import (
	"time"

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
	Owner    string
	Virtual  bool            //Real money or assets?
	External bool            //External means that can't accept transaction
	Balance  decimal.Decimal // This should probably be computed stored
	System   bool            //Is this a system account
	// More properties
}

/*type Balance struct {
	AccountID string
	Computed  DateTime
	At        DateTime
	Balance   decimal.Decimal
}*/

type Transaction struct {
	ID          int
	Timestamp   time.Time
	From        string
	To          string
	Amount      decimal.Decimal
	OperationID string
}

type Operation struct {
	ID           string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    string
	////////
	Timestamp   time.Time
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
