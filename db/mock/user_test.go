package mock

import (
	"testing"

	"ronche.se/moneytracker/db"
)

func TestGetUser(t *testing.T) {

	ms := NewMockStore()
	db.Populate(ms)

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

	ms := NewMockStore()
	db.Populate(ms)

	ul, err := ms.GetUsers()
	if err != nil {
		t.Fatal(err)
	}

	if len(ul) != 2 {
		t.Fatal(ul)
	}
}
