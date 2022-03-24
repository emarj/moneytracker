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
		OwnersID:    []string{"marco"},
		DisplayName: "Primary M",
	}
	acc1.AlterBalance(decimal.New(100, 0))

	acc2 := domain.Account{
		Name:        "secondary",
		OwnersID:    []string{"marco"},
		DisplayName: "Secondary M",
	}
	acc2.AlterBalance(decimal.New(50, 0))

	acc3 := domain.Account{
		Name:        "shared",
		OwnersID:    []string{"marco", "arianna"},
		DisplayName: "Shared AM 🌍",
	}
	acc3.AlterBalance(decimal.New(200, 0))

	acc4 := domain.Account{
		Name:        "primary",
		OwnersID:    []string{"arianna"},
		DisplayName: "Primary A",
	}
	acc4.AlterBalance(decimal.New(100, 0))

	s.AddAccount(&acc1)
	s.AddAccount(&acc2)
	s.AddAccount(&acc3)
	s.AddAccount(&acc4)

	s.InsertTransaction(&domain.Transaction{
		Date:        domain.Date{time.Now()},
		Description: "TX1",
		Notes:       "",
		Amount:      decimal.New(30, 0),
		FromID:      acc1.ID(),
		ToID:        acc2.ID(),
		Type:        "",
	})

	s.InsertTransaction(&domain.Transaction{
		Date:        domain.Date{time.Now()},
		Description: "TX2",
		Notes:       "",
		Amount:      decimal.New(120, 0),
		FromID:      acc2.ID(),
		ToID:        acc3.ID(),
		Type:        "",
	})
}
