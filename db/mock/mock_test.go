package mock

import (
	"testing"

	"ronche.se/moneytracker/domain"
)

func TestStoreGetUser(t *testing.T) {

	ms := NewMockStore()

	ms.AddUser(&domain.User{ID: "marco", Name: "Marco"})

	al, err := ms.GetAccountsByUser("marco")
	if err != nil {
		t.Fatal(err)
	}
	if len(al) == 0 {
		t.Fatal(al)
	}
}
