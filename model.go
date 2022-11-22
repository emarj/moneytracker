package moneytracker

import (
	"github.com/shopspring/decimal"
)

type Entity struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	System   bool   `json:"is_system"`
	External bool   // For example a friend that owes me
}

type Account struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	EntityID    int    `json:"entity_id"` // TODO: Allow for shared accounts
	IsSystem    bool   `json:"is_system"` // Was this created by system or by user?
	IsWorld     bool   `json:"is_world"`
	IsCredit    bool   `json:"is_credits"`
}

type Balance struct {
	AccountID int              `json:"id"`
	Timestamp DateTime         `json:"timestamp"`
	Value     *decimal.Decimal `json:"value"` // The reason to use a pointer is to have a nice api, see server.go
}

type Transaction struct {
	ID        int             `json:"id"`
	From      Account         `json:"from"`
	To        Account         `json:"to"`
	Amount    decimal.Decimal `json:"amount"`
	Operation Operation       `json:"operation"`
}

type Credit struct {
	ID          int             `json:"id"`
	Debtor      Entity          `json:"debtor"`
	Creditor    Entity          `json:"creditor"`
	Account     Account         `json:"account"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Operation   Operation       `json:"operation"`
}

type Operation struct {
	ID int `json:"id"`
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
