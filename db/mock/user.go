package mock

import (
	"fmt"

	"ronche.se/moneytracker/domain"
)

func (m *MockStore) GetUsers() ([]*domain.User, error) {
	ul := make([]*domain.User, len(m.users))

	k := 0
	for _, u := range m.users {
		ul[k] = u
		k++
	}

	return ul, nil
}

func (m *MockStore) GetUser(uID string) (*domain.User, error) {
	u, ok := m.users[uID]
	if !ok {
		return nil, fmt.Errorf("a user with id=%s does not exists", uID)
	}

	return u, nil
}

func (m *MockStore) AddUser(u *domain.User) (string, error) {
	_, ok := m.users[u.ID]
	if ok {
		return "", fmt.Errorf("a user with id=%s already not exists", u.ID)
	}

	m.users[u.ID] = u

	return u.ID, nil
}
