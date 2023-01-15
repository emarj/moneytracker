package sqlite

import mt "github.com/emarj/moneytracker"

func (s *SQLiteStore) GetOperationTypes() []mt.OperationType {
	return s.operationTypes
}
func (s *SQLiteStore) GetAccountTypes() []mt.AccountType {
	return s.accountTypes
}
