package moneytracker

import (
	"encoding/json"

	"github.com/emarj/moneytracker/datetime"
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

//go:generate easytags $GOFILE json

type Record struct {
	ID          null.Int          `json:"id" sql:"primary_key"`
	CreatedOn   datetime.DateTime `json:"created_on"`
	ModifiedOn  datetime.DateTime `json:"modified_on"`
	CreatedByID int64             `json:"created_by_id"`
}

type User struct {
	ID          null.Int      `json:"id" sql:"primary_key"`
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	IsAdmin     bool          `json:"is_admin"`
	Shares      []EntityShare `json:"shares"`
}

type EntityShare struct {
	UserID   int64 `json:"user_id"`
	EntityID int64 `json:"entity_id"`
	Quota    int64 `json:"quota"`
}

type Entity struct {
	ID         null.Int      `json:"id" sql:"primary_key"`
	Name       string        `json:"name"`
	IsSystem   bool          `json:"is_system"`
	IsExternal bool          `json:"is_external"` // For example a friend that owes me
	Shares     []EntityShare `json:"shares"`
}

const EntSystemID int64 = 0

func SystemEntities() []Entity {
	return []Entity{
		{ID: null.IntFrom(EntSystemID), Name: "system", IsSystem: true, IsExternal: false},
	}
}

type Account struct {
	ID          null.Int     `json:"id" sql:"primary_key"`
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	OwnerID     int64        `json:"owner_id"`
	Owner       *Entity      `json:"owner,omitempty" alias:"owner"` // TODO: Allow for shared accounts
	IsSystem    bool         `json:"is_system"`                     // This can't be deleted by the user
	IsGroup     bool         `json:"is_group"`                      // This should not be type inside type
	TypeID      int64        `json:"type_id"`
	Type        *AccountType `json:"type,omitempty"`
	GroupID     null.Int     `json:"group_id"`
}

const AccWorldID int64 = 0

func SystemAccounts() []Account {
	return []Account{
		{
			ID:          null.IntFrom(AccWorldID),
			Name:        "world",
			DisplayName: "World",
			OwnerID:     EntSystemID,
			IsSystem:    true,
			TypeID:      AccTypeWorld,
		},
	}
}

type AccountType struct {
	ID     int64  `json:"id" sql:"primary_key"`
	Name   string `json:"name"`
	System bool   `json:"system"` //System reserved
}

const (
	AccTypeMoney int64 = iota // This is first since must be the default
	AccTypeWorld
	AccTypeCredit
	AccTypeInvestment
)

func AccountTypes() []AccountType {
	return []AccountType{
		{ID: AccTypeWorld, Name: "world", System: true},
		{ID: AccTypeMoney, Name: "money"},
		{ID: AccTypeCredit, Name: "credit"},
		{ID: AccTypeInvestment, Name: "investment"},
	}
}

type Balance struct {
	AccountID   null.Int            `json:"account_id" sql:"primary_key"`
	Account     *Account            `json:"account,omitempty"`
	Timestamp   datetime.DateTime   `json:"timestamp" sql:"primary_key"`
	Value       decimal.Decimal     `json:"value"`
	Delta       decimal.NullDecimal `json:"delta"`
	IsComputed  bool                `json:"is_computed"`
	Comment     string              `json:"comment"`
	OperationID null.Int            `json:"operation_id"`
	Operation   *Operation          `json:"operation,omitempty"`
}

type Transaction struct {
	ID          null.Int          `json:"id" sql:"primary_key"`
	Timestamp   datetime.DateTime `json:"timestamp"`
	FromID      int64             `json:"from_id"`
	From        *Account          `json:"from" alias:"from"`
	ToID        int64             `json:"to_id"`
	To          *Account          `json:"to" alias:"to"`
	Amount      decimal.Decimal   `json:"amount"`
	Comment     string            `json:"comment"`
	OperationID int64             `json:"operation_id"`
	Operation   *Operation        `json:"operation"`
}

type Operation struct {
	ID          null.Int          `json:"id" sql:"primary_key"`
	CreatedOn   datetime.DateTime `json:"created_on"`
	ModifiedOn  datetime.DateTime `json:"modified_on"`
	CreatedByID int64             `json:"created_by_id"`
	//Shares []Entity
	Description  string         `json:"description"`
	TypeID       int64          `json:"type_id"`
	Type         *OperationType `json:"type,omitempty"`
	Transactions []Transaction  `json:"transactions,omitempty"`
	Balances     []Balance      `json:"balances,omitempty"`
	//////////////////////////////////////////////
	CategoryID int64           `json:"category_id"`
	Details    json.RawMessage `json:"details,omitempty"`
	ParentID   null.Int        `json:"parent_id"`
	//Parent       *Operation    `json:"parent,omitempty"`
}

type OperationType struct {
	ID   int64  `json:"id" sql:"primary_key"`
	Name string `json:"name"`
}

// This must be the same as in schema.sql
const (
	OpTypeOther         int64 = iota
	OpTypeBalanceAdjust       // A balance adjust
	OpTypeExpense             // Something that enters the system?
	OpTypeIncome              // Something that exits the system?
	OpTypeTransfer            // Just a transfer
)

func OperationTypes() []OperationType {
	return []OperationType{
		{OpTypeOther, "other"},
		{OpTypeBalanceAdjust, "balance"},
		{OpTypeExpense, "expense"},
		{OpTypeIncome, "income"},
		{OpTypeTransfer, "transfer"},
	}
}

type Category struct {
	ID       null.Int        `json:"id" sql:"primary_key"`
	Name     string          `json:"name"`
	FullName string          `json:"full_name,omitempty"`
	ParentID null.Int        `json:"parent_id"`
	Parent   *ParentCategory `json:"parent" alias:"parent"`
}

type ParentCategory Category

const CatUncategorized int64 = 0

func SystemCategories() []Category {
	return []Category{
		{ID: null.IntFrom(CatUncategorized), Name: "Uncategorized"},
	}
}
