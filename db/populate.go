package db

import (
	"time"

	"github.com/shopspring/decimal"
	"ronche.se/moneytracker/domain"
)

func Populate(s Store) {

	u1 := domain.User{ID: "marco", Name: "Marco"}
	u2 := domain.User{ID: "arianna", Name: "Arianna"}
	s.AddUser(&u1)
	s.AddUser(&u2)

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

	s.AddAccount(&acc1)
	s.AddAccount(&acc2)
	s.AddAccount(&acc3)
	s.AddAccount(&acc4)

	s.AddTransaction(&domain.Transaction{
		Date:        time.Now(),
		Description: "Transferimento",
		Notes:       "",
		Amount:      decimal.New(30, 0),
		FromID:      acc1.ID,
		ToID:        acc2.ID,
	})

	s.AddTransaction(&domain.Transaction{
		Date:        time.Now(),
		Description: "TX2",
		Notes:       "",
		Amount:      decimal.New(120, 0),
		FromID:      acc2.ID,
		ToID:        acc3.ID,
	})
}
