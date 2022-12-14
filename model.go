package moneytracker

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

type Record struct {
	ID           null.Int `json:"id"`
	DateCreated  DateTime `json:"date_created"`
	DateModified DateTime `json:"date_modified"`
	CreatedByID  int      `json:"created_by_id"`
}

type User struct {
	ID          null.Int `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Admin       bool     `json:"admin"`
}

type Entity struct {
	ID       null.Int `json:"id"`
	Name     string   `json:"name"`
	System   bool     `json:"is_system"`
	External bool     `json:"is_external"` // For example a friend that owes me
}

type Account struct {
	ID          null.Int `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Owner       Entity   `json:"owner"`     // TODO: Allow for shared accounts
	IsSystem    bool     `json:"is_system"` // Was this created by system or by user?
	IsWorld     bool     `json:"is_world"`
	Type        int      `json:"type"`
}

/*type AccountType struct {
	ID          null.Int `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
}*/

const (
	AccountMoney int = iota
	AccountCredit
	AccountGroup
)

type Balance struct {
	AccountID null.Int        `json:"account_id"`
	Timestamp DateTime        `json:"timestamp"`
	Value     decimal.Decimal `json:"value"`
	Comment   string          `json:"comment"`
}

type Transaction struct {
	ID        null.Int        `json:"id"`
	From      Account         `json:"from"`
	To        Account         `json:"to"`
	Amount    decimal.Decimal `json:"amount"`
	Operation Operation       `json:"operation"`
}

type Operation struct {
	ID null.Int `json:"id"`
	//DateCreated  DateTime `json:"date_created"`
	//DateModified DateTime `json:"date_modified"`
	CreatedByID int `json:"created_by_id"` // This should be User.ID
	////////
	Timestamp    DateTime      `json:"timestamp"`
	Description  string        `json:"description"`
	Transactions []Transaction `json:"transactions"`
	Type         int           `json:"type"`
	CategoryID   int           `json:"category_id"`
	//Details      string        `json:"details"`
}

const (
	OpTypeOther int = iota
	OpTypeExpense
	OpTypeTransfer
	OpTypeBalanceAdjust
)

type Category struct {
	ID   null.Int `json:"id"`
	Name string   `json:"name"`
}
