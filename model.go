package moneytracker

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
	"ronche.se/moneytracker/datetime"
)

//go:generate easytags $GOFILE json

type Record struct {
	ID          null.Int          `json:"id" sql:"primary_key"`
	CreatedOn   datetime.DateTime `json:"created_on"`
	ModifiedOn  datetime.DateTime `json:"modified_on"`
	CreatedByID int               `json:"created_by_id"`
}

type User struct {
	ID          null.Int `json:"id" sql:"primary_key"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	IsAdmin     bool     `json:"is_admin"`
	Password    string   `json:"password"`
}

type Entity struct {
	ID         null.Int `json:"id" sql:"primary_key"`
	Name       string   `json:"name"`
	IsSystem   bool     `json:"is_system"`
	IsExternal bool     `json:"is_external"` // For example a friend that owes me
}

type Account struct {
	ID          null.Int `json:"id" sql:"primary_key"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Owner       Entity   `json:"owner" alias:"owner" mapping:".ID:owner_id"` // TODO: Allow for shared accounts
	IsSystem    bool     `json:"is_system"`                                  // This can't be deleted by the user
	IsWorld     bool     `json:"is_world"`
	IsGroup     bool     `json:"is_group"` // This should not be type inside type
	Type        int      `json:"type"`
	ParentID    null.Int `json:"parent_id"`
}

// These must the same as in schema.sql
const (
	AccountMoney int = iota
	AccountCredit
	AccountInvestment
)

type Balance struct {
	AccountID  null.Int            `json:"account_id" sql:"primary_key"`
	Timestamp  datetime.DateTime   `json:"timestamp" sql:"primary_key"`
	Value      decimal.Decimal     `json:"value"`
	Delta      decimal.NullDecimal `json:"delta"`
	IsComputed bool                `json:"is_computed"`
	Operation  Operation           `json:"operation" mapping:".ID:operation_id"`
}

type Transaction struct {
	ID        null.Int          `json:"id" sql:"primary_key"`
	Timestamp datetime.DateTime `json:"timestamp"`
	From      Account           `json:"from" alias:"from" mapping:".ID:from_id"`
	To        Account           `json:"to" alias:"to" mapping:".ID:to_id"`
	Amount    decimal.Decimal   `json:"amount"`
	Comment   string            `json:"comment"`
	Operation Operation         `json:"operation" mapping:".ID:operation_id"`
}

type Operation struct {
	ID          null.Int          `json:"id" sql:"primary_key"`
	CreatedOn   datetime.DateTime `json:"created_on"`
	ModifiedOn  datetime.DateTime `json:"modified_on"`
	CreatedByID int               `json:"created_by_id"`
	//Shares []Entity
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Transactions []Transaction `json:"transactions"`
	Balances     []Balance     `json:"balances"`
	TypeID       int           `json:"type_id"`
	CategoryID   int           `json:"category_id"`
	//Parent       *Operation    `json:"parent"`
}

// This must be the same as in schema.sql
const (
	OpTypeOther    int = iota
	OpTypeExpense      // Something that enters the system?
	OpTypeIncome       // Something that exits the system?
	OpTypeTransfer     // Just a transfer
	OpTypeBalance      // A balance adjust
)

type Category struct {
	ID       null.Int `json:"id" sql:"primary_key"`
	Name     string   `json:"name"`
	ParentID null.Int `json:"parent_id"`
}
