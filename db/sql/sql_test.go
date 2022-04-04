package sql

import (
	"testing"

	"ronche.se/moneytracker/domain"
)

func TestMigrate(t *testing.T) {
	s, err := NewInMemoryStore()
	if err != nil {
		t.Fatal(err)
	}

	err = s.Migrate()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUser(t *testing.T) {
	s, err := NewInMemoryStore()
	if err != nil {
		t.Fatal(err)
	}

	err = s.Migrate()
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddUser(&domain.User{
		ID:   "user1",
		Name: "User 1",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddUser(&domain.User{
		ID:   "user2",
		Name: "User 2",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddUser(&domain.User{
		ID:   "user2",
		Name: "User 2",
	})
	if err == nil {
		t.Fatal(err)
	}

	u, err := s.GetUser("user2")
	if err != nil {
		t.Fatal(err)
	}

	if u.Name != "User 2" {
		t.Fatal(u)
	}

	ul, err := s.GetUsers()
	if err != nil {
		t.Fatal(err)
	}

	if len(ul) != 2 {
		t.Fatal(ul)
	}

}
