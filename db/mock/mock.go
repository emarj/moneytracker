package mock

import "ronche.se/moneytracker/domain"

type MockStore struct {
	users        map[string]*domain.User
	accounts     map[string]*domain.Account
	transactions map[string]*domain.Transaction
}

func NewMockStore() *MockStore {
	return &MockStore{
		users:        map[string]*domain.User{},
		accounts:     map[string]*domain.Account{},
		transactions: map[string]*domain.Transaction{},
	}
}
