package moneytracker

type Store interface {
	GetAccounts() ([]Account, error)
	AddAccount(a Account) (int, error)

	GetTransaction(int) (Transaction, error)
	GetTransactions() ([]Transaction, error)
	GetTransactionsByAccount(aID int) ([]Transaction, error)
	AddTransaction(t Transaction) (int, error)
}
