package mock

import (
	"testing"

	"ronche.se/moneytracker/domain"
)

func populateAccounts(as *mockAccountStore) {

	u1 := domain.User{ID: "marco", Name: "Marco"}
	u2 := domain.User{ID: "arianna", Name: "Arianna"}

	acc1 := domain.Account{
		Name:        "primary",
		Owners:      []domain.User{u1},
		DisplayName: "Primary M",
	}

	acc2 := domain.Account{
		Name:        "secondary",
		Owners:      []domain.User{u1},
		DisplayName: "Secondary M",
	}

	acc3 := domain.Account{
		Name:        "shared",
		Owners:      []domain.User{u1, u2},
		DisplayName: "Shared AM",
	}

	acc4 := domain.Account{
		Name:        "primary",
		Owners:      []domain.User{u2},
		DisplayName: "Primary A",
	}

	as.AddAccount(&acc1)
	as.AddAccount(&acc2)
	as.AddAccount(&acc3)
	as.AddAccount(&acc4)
}

func TestGetAccountByUserAndName(t *testing.T) {

	as := newMockAccountStore()

	populateAccounts(as)

	var al []*domain.Account
	var err error

	al, err = as.GetAccountsByUserAndName("marco", "primary")
	if err != nil || len(al) != 1 {
		t.Fatal(err)
	}

	al, err = as.GetAccountsByUserAndName("arianna", "shared")
	if err != nil || len(al) != 1 {
		t.Fatal(err)
	}

	al, _ = as.GetAccountsByUserAndName("marco", "NotExistent")
	if len(al) != 0 {
		t.Fatal(err)
	}

	al, _ = as.GetAccountsByUserAndName("sda", "shared")
	if len(al) != 0 {
		t.Fatal(err)
	}

	al, _ = as.GetAccountsByUserAndName("sdsdasda", "gdfgfshared")
	if len(al) != 0 {
		t.Fatal(err)
	}
}

func TestGetAccountsByUser(t *testing.T) {

	as := newMockAccountStore()

	populateAccounts(as)

	al, _ := as.GetAccountsByUser("nonexistentuser")
	if len(al) != 0 {
		t.Fatal(al)
	}

	al, _ = as.GetAccountsByUser("marco")
	if len(al) != 3 {
		t.Fatal(al)
	}

	al, _ = as.GetAccountsByUser("arianna")
	if len(al) != 2 {
		t.Fatal(al)
	}
}
