package moneytracker

import (
	"github.com/shopspring/decimal"
)

type Entity struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	System bool   `json:"is_system"`
	//External bool //?
}

type Account struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	OwnerID     int    `json:"owner_id"`    // TODO: Allow for shared accounts
	IsMoney     bool   `json:"is_money"`    // Is it money or assets?
	IsExternal  bool   `json:"is_external"` // Allow for direct balance manipulation
	IsSystem    bool   `json:"is_system"`   // Is this a system account
	//Balance  *decimal.Decimal `json:"omitempty"` // This would be computed from the transactions
}

type Balance struct {
	AccountID  int              `json:"id"`
	Timestamp  DateTime         `json:"timestamp"`
	Value      *decimal.Decimal `json:"value"`
	IsComputed bool             `json:"is_computed"`
	Notes      string           `json:"notes"`
}

type Transaction struct {
	ID          int             `json:"id"`
	Timestamp   DateTime        `json:"timestamp"`
	FromID      string          `json:"from_id"`
	ToID        string          `json:"to_id"`
	Amount      decimal.Decimal `json:"amount"`
	OperationID string          `json:"operation_id"`
}

type Operation struct {
	ID           string   `json:"id"`
	DateCreated  DateTime `json:"date_created"`
	DateModified DateTime `json:"date_modified"`
	CreatedByID  string   `json:"create_by_id"`
	////////
	Timestamp   DateTime `json:"timestamp"`
	Description string   `json:"description"`
	Movements   []string `json:"movements"`
	CategoryID  int      `json:"cateogry_id"`
	Details     string   `json:"details"`
}

const (
	CategoryExpense = iota
	CategoryTransfer
	CategorySharedExpense
	CategoryGift
)
