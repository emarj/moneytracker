package mock

import (
	"fmt"

	"github.com/gofrs/uuid"
	"ronche.se/moneytracker/domain"
)

type mockAccountStore struct {
	accounts map[string]*domain.Account
}

func newMockAccountStore() *mockAccountStore {
	return &mockAccountStore{
		accounts: map[string]*domain.Account{},
	}
}

func (as *mockAccountStore) GetAccountsByUser(uID string) ([]*domain.Account, error) {
	al := []*domain.Account{}

	for _, a := range as.accounts {
		for _, o := range a.Owners {
			if o.ID == uID {
				al = append(al, a)
			}
		}
	}

	return al, nil
}

func (as *mockAccountStore) GetAccountsByUserAndName(uID string, name string) ([]*domain.Account, error) {
	al := []*domain.Account{}

	for _, a := range as.accounts {
		if a.Name == name {
			for _, o := range a.Owners {
				if o.ID == uID {
					al = append(al, a)
				}
			}
		}
	}

	return al, nil
}

func (as *mockAccountStore) GetAccount(aID uuid.UUID) (*domain.Account, error) {
	a, ok := as.accounts[aID.String()]
	if !ok {
		return nil, fmt.Errorf("an account with id=%s does not exist", aID.String())
	}

	return a, nil
}

func (as *mockAccountStore) AddAccount(a *domain.Account) error {
	var err error
	a.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	// Double check for uuid collision
	_, ok := as.accounts[a.ID.String()]
	if ok {
		return fmt.Errorf("CRITICAL: UUID collision detected. An account with id=%s already exists", a.ID)
	}

	as.accounts[a.ID.String()] = a

	return nil
}
