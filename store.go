package moneytracker

import (
	"github.com/emarj/moneytracker/timestamp"
)

type Store interface {
	GetEntities() ([]Entity, error)
	GetEntitiesOfUser(int64) ([]Entity, error)
	GetEntity(eID int64) (*Entity, error)
	AddEntity(e *Entity) error
	AddSharesForEntity(...EntityShare) error

	GetAccounts() ([]Account, error)
	GetAccountsByEntity(eID int64) ([]Account, error)
	GetUserAccounts(uID int64) ([]Account, error)
	GetAccount(aID int64) (*Account, error)
	AddAccount(a *Account) error
	DeleteAccount(aID int64) error

	GetBalanceAt(aID int64, time timestamp.Timestamp) (Balance, error)
	GetBalanceNow(aID int64) (Balance, error)
	GetLastBalance(aID int64) (Balance, error)
	GetBalanceHistory(aID int64) ([]Balance, error)
	SnapshotBalance(aID int64) error
	SetBalance(b *Balance) error

	GetOperation(int64) (*Operation, error)
	//GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID int64, limit int64) ([]Transaction, error)
	GetOperationsOfUser(uID int64, limit int64) ([]Operation, error)
	GetOperationsByEntity(eID int64, limit int64) ([]Operation, error)
	AddOperation(op *Operation) error
	DeleteOperation(opID int64) error

	GetCategories() ([]Category, error)
	AddCategory(fullName string) (Category, error)

	GetOperationTypes() []OperationType
	GetAccountTypes() []AccountType

	Login(user string, password string) (User, error)
	RegisterUser(user *User, password string) error
}
