package moneytracker

import (
	"github.com/gofrs/uuid"
)

type Store interface {
	GetAccounts() ([]Account, error)
	AddAccount(a Account) error
	GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID uuid.UUID) ([]Transaction, error)
	AddTransaction(t Transaction) error
}
