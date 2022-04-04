package mock

import (
	"ronche.se/moneytracker/db"
	"ronche.se/moneytracker/domain"
)

type MockStore struct {
	*mockUserStore
	*mockAccountStore
	*mockTransactionStore
}

func NewMockStore() db.Store {
	return &MockStore{
		newMockUserStore(),
		newMockAccountStore(),
		newMockTransactionStore(),
	}
}

//This overrides mockUserStore.AddUser
func (ms *MockStore) AddUser(u *domain.User) error {
	err := ms.mockUserStore.AddUser(u)
	if err != nil {
		return err
	}

	//Add default accounts
	defaultAccounts := []*domain.Account{
		{Name: "worldin", Owners: []domain.User{*u}, DisplayName: "🌍 World IN", Default: true},
		{Name: "worldout", Owners: []domain.User{*u}, DisplayName: "🌍 World OUT", Default: true},
		{Name: "debits", Owners: []domain.User{*u}, DisplayName: "Debits", Default: true},
		{Name: "credits", Owners: []domain.User{*u}, DisplayName: "Credits", Default: true},
	}

	for _, a := range defaultAccounts {
		ms.mockAccountStore.AddAccount(a)
	}

	return nil

}

/*
_,  */
