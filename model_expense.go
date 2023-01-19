package moneytracker

import (
	"github.com/emarj/moneytracker/timestamp"
	"github.com/shopspring/decimal"
)

//go:generate easytags $GOFILE json
type ExpenseShare struct {
	Amount        decimal.Decimal `json:"amount"`
	Quota         int64           `json:"quota"`
	Total         decimal.Decimal `json:"total"`
	WithID        int64           `json:"with_id"`
	CredAccountID int64           `json:"cred_account_id"`
	DebAccountID  int64           `json:"deb_account_id"`
	IsCredit      bool            `json:"is_credit"`
}

type Expense struct {
	Timestamp   timestamp.Timestamp `json:"timestamp"`
	Amount      decimal.Decimal     `json:"amount"`
	Description string              `json:"description"`
	AccountID   int64               `json:"account_id"`
	CategoryID  int64               `json:"category_id"`
	Shares      []ExpenseShare      `json:"shares"`
}

func (e Expense) ToOperation() Operation {
	op := Operation{
		Description: e.Description,
		TypeID:      OpTypeExpense,
		Transactions: []Transaction{
			{
				Timestamp: e.Timestamp,
				FromID:    e.AccountID,
				ToID:      AccWorldID,
				Amount:    e.Amount,
			},
		},
		CategoryID: e.CategoryID,
	}

	for _, s := range e.Shares {
		op.Transactions = append(op.Transactions, Transaction{
			Timestamp: e.Timestamp,
			Amount:    s.Amount,
			ToID:      s.CredAccountID,
			FromID:    s.DebAccountID,
		})
	}

	return op
}
