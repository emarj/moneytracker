package mock

import (
	"fmt"

	"ronche.se/moneytracker/domain"
)

type mockUserStore struct {
	users map[string]*domain.User
}

func newMockUserStore() *mockUserStore {
	return &mockUserStore{
		users: map[string]*domain.User{},
	}
}

func (us *mockUserStore) GetUsers() ([]*domain.User, error) {
	ul := make([]*domain.User, len(us.users))

	k := 0
	for _, u := range us.users {
		ul[k] = u
		k++
	}

	return ul, nil
}

func (us *mockUserStore) GetUser(uID string) (*domain.User, error) {
	u, ok := us.users[uID]
	if !ok {
		return nil, fmt.Errorf("a user with id=%s does not exists", uID)
	}

	return u, nil
}

func (us *mockUserStore) AddUser(u *domain.User) error {
	_, ok := us.users[u.ID]
	if ok {
		return fmt.Errorf("a user with id=%s already not exists", u.ID)
	}

	us.users[u.ID] = u

	return nil
}
