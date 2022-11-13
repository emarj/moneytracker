package moneytracker

type Store interface {
	GetEntities() ([]Entity, error)
	GetEntity(eID int) (*Entity, error)
	AddEntity(e Entity) (int, error)

	GetAccounts() ([]Account, error)
	GetAccountsByEntity(eID int) ([]Account, error)
	GetAccount(aID int) (*Account, error)
	AddAccount(a Account) (int, error)
	DeleteAccount(aID int) error

	GetBalance(aID int) (*Balance, error)
	GetBalances(aID int) ([]Balance, error)
	ComputeBalance(aID int) error
	AddBalance(b Balance) error

	GetTransaction(int) (*Transaction, error)
	GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID int) ([]Transaction, error)
	GetTransactionsByEntity(eID int) ([]Transaction, error)
	AddTransaction(t Transaction) (int, error)
	DeleteTransaction(tID int) error
}
