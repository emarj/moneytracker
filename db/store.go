package db

import (
	"github.com/gofrs/uuid"
	"ronche.se/moneytracker/domain"
)

type Store interface {
	//User
	GetUsers() ([]*domain.User, error)
	GetUser(uID string) (*domain.User, error)
	AddUser(u *domain.User) error
	//Account.Get
	GetAccount(aID uuid.UUID) (*domain.Account, error)
	GetAccountsByUser(uID string) ([]*domain.Account, error)
	GetAccountsByUserAndName(uID string, name string) ([]*domain.Account, error)

	AddAccount(a *domain.Account) error
	//Transactions.Get
	GetTransaction(tID uuid.UUID) (*domain.Transaction, error)
	GetTransactionsByAccount(aID uuid.UUID) ([]*domain.Transaction, error)
	GetTransactionsByUser(uID string) ([]*domain.Transaction, error)
	//Transactions.Create
	AddTransaction(t *domain.Transaction) error
	//Transactions.Update
	UpdateTransaction(t *domain.Transaction) error
	//Transactions.Delete
	DeleteTransaction(tID uuid.UUID) error
}
