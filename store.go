package moneytracker

import "github.com/emarj/moneytracker/datetime"

type Store interface {
	GetEntities() ([]Entity, error)
	GetEntity(eID int64) (*Entity, error)
	AddEntity(e *Entity) error

	GetAccounts() ([]Account, error)
	GetAccountsByEntity(eID int64) ([]Account, error)
	GetAccount(aID int64) (*Account, error)
	AddAccount(a *Account) error
	DeleteAccount(aID int64, onlyIfEmpty bool) error

	GetBalanceAt(aID int64, time datetime.DateTime) (Balance, error)
	GetBalanceNow(aID int64) (Balance, error)
	GetLastBalance(aID int64) (Balance, error)
	GetBalanceHistory(aID int64) ([]Balance, error)
	SnapshotBalance(aID int64) error
	SetBalance(b *Balance) error

	GetOperation(int64) (*Operation, error)
	//GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID int64, limit int64) ([]Transaction, error)
	GetOperationsByEntity(eID int64, limit int64) ([]Operation, error)
	AddOperation(op *Operation) error
	DeleteOperation(opID int64) error

	GetCategories() ([]Category, error)

	Login(user string, passwordHash []byte) (bool, error)
}
