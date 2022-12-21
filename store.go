package moneytracker

import "ronche.se/moneytracker/datetime"

type Store interface {
	GetEntities() ([]Entity, error)
	GetEntity(eID int) (*Entity, error)
	AddEntity(e *Entity) error

	GetAccounts() ([]Account, error)
	GetAccountsByEntity(eID int) ([]Account, error)
	GetAccount(aID int) (*Account, error)
	AddAccount(a *Account, initialBalance *Balance) error
	DeleteAccount(aID int, onlyIfEmpty bool) error

	GetBalanceAt(aID int, time datetime.DateTime) (*Balance, error)
	GetBalanceNow(aID int) (*Balance, error)
	GetHistory(aID int) ([]Balance, error)
	SnapshotBalance(aID int) error
	SetBalance(b Balance) error

	GetOperation(int) (*Operation, error)
	//GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID int, limit int) ([]Transaction, error)
	GetOperationsByEntity(eID int, limit int) ([]Operation, error)
	AddOperation(op *Operation) error
	DeleteOperation(tID int) error

	GetCategories() ([]Category, error)

	Login(user string, passwordHash []byte) (bool, error)
}
