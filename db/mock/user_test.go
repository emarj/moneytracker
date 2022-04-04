package mock

import (
	"testing"

	"ronche.se/moneytracker/domain"
)

func populateUserStore(us *mockUserStore) {
	u1 := domain.User{ID: "marco", Name: "Marco"}
	u2 := domain.User{ID: "arianna", Name: "Arianna"}
	us.AddUser(&u1)
	us.AddUser(&u2)
}

func TestGetUser(t *testing.T) {

	ms := newMockUserStore()
	populateUserStore(ms)

	u, err := ms.GetUser("marco")
	if err != nil {
		t.Fatal(err)
	}

	if u.ID != "marco" || u.Name != "Marco" {
		t.Fatal(u)
	}

	u, err = ms.GetUser("F")
	if err == nil {
		t.Fatal(u)
	}
}

func TestGetUsers(t *testing.T) {

	ms := newMockUserStore()
	populateUserStore(ms)

	ul, err := ms.GetUsers()
	if err != nil {
		t.Fatal(err)
	}

	if len(ul) != 2 {
		t.Fatal(ul)
	}
}
