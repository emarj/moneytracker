package moneytracker

import "gopkg.in/guregu/null.v4"

type Store interface {
	GetEntities() ([]Entity, error)
	GetEntity(eID int) (*Entity, error)
	AddEntity(e Entity) (null.Int, error)

	GetAccounts() ([]Account, error)
	GetAccountsByEntity(eID int) ([]Account, error)
	GetAccount(aID int) (*Account, error)
	AddAccount(a Account) (null.Int, error)
	DeleteAccount(aID int) error

	GetBalance(aID int) (*Balance, error)
	GetBalances(aID int) ([]Balance, error)
	ComputeBalance(aID int) error
	AddBalance(b Balance) error

	GetOperation(int) (*Operation, error)
	//GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID int) ([]Transaction, error)
	GetOperationsByEntity(eID int) ([]Operation, error)
	AddOperation(op Operation) (null.Int, error)
	DeleteOperation(tID int) error

	GetCategories() ([]Category, error)

	Login(user string, passwordHash []byte) (bool, error)
}
