package db

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/domain"
)

type Store interface {
	UserStore
	AccountStore
	TransactionStore
	//ExpensesStore
}

type UserStore interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uID string) (*domain.User, error)
	AddUser(u *domain.User) (string, error)
}

type AccountStore interface {
	GetAccount(aID string) (*domain.Account, error)
	GetAccountsOfUser(uID string) ([]*domain.Account, error)

	AddAccount(a *domain.Account) (string, error)

	Balance(aID string, amount decimal.Decimal) (*domain.Account, error)
}

type TransactionStore interface {
	//CREATE
	InsertTransaction(t *domain.Transaction) (uuid.UUID, error)
	//DELETE
	DeleteTransaction(tID uuid.UUID) error
	//UPDATE
	UpdateTransaction(t *domain.Transaction) (uuid.UUID, error)
	//GET
	GetTransaction(id uuid.UUID) (*domain.Transaction, error)
	GetTransactionsByAccount(aID string) ([]*domain.Transaction, error)
	GetTransactionsByUser(uID string) ([]*domain.Transaction, error)
}

/*
type ExpensesStore interface {
	GetExpense(id uuid.UUID) (domain.Expense, error)

	GetExpensesByAccount(aID uuid.UUID) ([]*domain.Expense, error)
	GetExpensesByUser(uID string) ([]*domain.Expense, error)
}*/
