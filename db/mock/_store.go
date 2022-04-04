package mock

import (
	"github.com/gofrs/uuid"
)

type Store interface {
	UserStore
	AccountStore
	TransactionStore
}

type UserStore interface {
	GetUsers() ([]*User, error)
	GetUser(uID string) (*User, error)
	AddUser(u *User) error
}

type AccountStore interface {
	GetAccount(aID uuid.UUID) (*Account, error)
	GetAccountsByUser(uID string) ([]*Account, error)
	GetAccountsByUserAndName(uID string, name string) ([]*Account, error)

	AddAccount(a *Account) error
}

type TransactionStore interface {
	//GET
	GetTransaction(tID uuid.UUID) (*Transaction, error)
	GetTransactionsByAccount(aID string) ([]*Transaction, error)
	GetTransactionsByUser(uID string) ([]*Transaction, error)
	//CREATE
	AddTransaction(t *Transaction) error
	//UPDATE
	UpdateTransaction(t *Transaction) error
	//DELETE
	DeleteTransaction(tID uuid.UUID) error
}

type ShareStore interface {
	GetShares(tID uuid.UUID) ([]*Share, error)
}
