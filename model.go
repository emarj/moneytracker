package moneytracker

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

type Entity struct {
	ID       null.Int `json:"id" sql:"primary_key"`
	Name     string   `json:"name"`
	System   bool     `json:"is_system"`
	External bool     `json:"is_external"` // For example a friend that owes me
}

type Account struct {
	ID          null.Int `json:"id" sql:"primary_key"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Owner       Entity   `json:"owner"`     // TODO: Allow for shared accounts
	IsSystem    bool     `json:"is_system"` // Was this created by system or by user?
	IsWorld     bool     `json:"is_world"`
	IsCredit    bool     `json:"is_credit"`
}

type Balance struct {
	AccountID int              `json:"account_id"`
	Timestamp DateTime         `json:"timestamp"`
	Value     *decimal.Decimal `json:"value"` // The reason to use a pointer is to have a nice api, see server.go
}

type Transaction struct {
	ID        null.Int        `json:"id" sql:"primary_key"`
	From      Account         `json:"from" alias:"From"`
	To        Account         `json:"to" alias:"To"`
	Amount    decimal.Decimal `json:"amount"`
	Operation Operation       `json:"operation"`
}

type Operation struct {
	ID null.Int `json:"id" sql:"primary_key"`
	//DateCreated  DateTime `json:"date_created"`
	//DateModified DateTime `json:"date_modified"`
	CreatedByID int `json:"created_by_id"`
	////////
	Timestamp    *DateTime     `json:"timestamp"`
	Description  string        `json:"description"`
	Transactions []Transaction `json:"transactions"`
	CategoryID   int           `json:"category_id"`
	//Details      string        `json:"details"`
}

const (
	CategoryExpense = iota
	CategoryTransfer
	CategorySharedExpense
	CategoryGift
)
