package mock

import (
	"fmt"

	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/domain"
)

func (m *MockStore) GetAccountsOfUser(uID string) ([]*domain.Account, error) {
	al := []*domain.Account{}

	for _, a := range m.accounts {
		for _, id := range a.OwnersID {
			if id == uID {
				al = append(al, a)
			}
		}
	}

	return al, nil
}

func (m *MockStore) GetAccount(aID string) (*domain.Account, error) {
	a, ok := m.accounts[aID]
	if !ok {
		return nil, fmt.Errorf("an account with id=%s does not exist", aID)
	}

	return a, nil
}

func (m *MockStore) AddAccount(a *domain.Account) (string, error) {
	aID := a.ID()
	_, ok := m.accounts[aID]
	if ok {
		return "", fmt.Errorf("an account with id=%s already exists", aID)
	}

	m.accounts[aID] = a

	return aID, nil
}

func (m *MockStore) Balance(id string, delta decimal.Decimal) (*domain.Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, fmt.Errorf("a account with id=%s does not exist", id)
	}

	a.AlterBalance(delta)

	return a, nil
}
