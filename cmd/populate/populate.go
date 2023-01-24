package main

import (
	"fmt"

	mt "github.com/emarj/moneytracker"
	"gopkg.in/guregu/null.v4"
)

func populate(s mt.Store) error {
	fmt.Print("Populating...")
	var err error

	user1 := mt.User{
		Name:        "arianna",
		DisplayName: "Arianna",
		IsAdmin:     false,
	}
	err = s.AddUser(&user1, "prova")
	if err != nil {
		return err
	}

	user2 := mt.User{
		Name:        "marco",
		DisplayName: "Marco",
		IsAdmin:     false,
	}
	err = s.AddUser(&user2, "pippo")
	if err != nil {
		return err
	}

	entUser1 := mt.Entity{
		ID:          null.IntFrom(1),
		Name:        "arianna",
		DisplayName: "Arianna",
		IsSystem:    false,
		IsExternal:  false,
		Shares: []mt.EntityShare{{
			UserID:   user1.ID.Int64,
			Quota:    100,
			Priority: null.IntFrom(0),
		}},
	}

	err = s.AddEntity(&entUser1)
	if err != nil {
		return err
	}

	entUser2 := mt.Entity{
		ID:          null.IntFrom(2),
		Name:        "marco",
		DisplayName: "Marco",
		IsSystem:    false,
		IsExternal:  false,
		Shares: []mt.EntityShare{{
			UserID:   user2.ID.Int64,
			Quota:    100,
			Priority: null.IntFrom(0),
		}},
	}

	err = s.AddEntity(&entUser2)
	if err != nil {
		return err
	}

	entUser3 := mt.Entity{
		ID:          null.IntFrom(3),
		Name:        "family",
		DisplayName: "Family",
		IsSystem:    false,
		IsExternal:  false,
	}

	err = s.AddEntity(&entUser3)
	if err != nil {
		return err
	}
	err = s.AddSharesForEntity(mt.EntityShare{
		UserID:   user1.ID.Int64,
		EntityID: entUser3.ID,
		Quota:    50,
	}, mt.EntityShare{
		UserID:   user2.ID.Int64,
		EntityID: entUser3.ID,
		Quota:    50,
	})
	if err != nil {
		return err
	}

	accounts := []mt.Account{{
		Name:        "portafoglio",
		DisplayName: "Portafoglio",
		OwnerID:     entUser1.ID.Int64,
		IsDefault:   true,
	}, {
		Name:        "cassa",
		DisplayName: "Cassa",
		OwnerID:     entUser1.ID.Int64,
	}, {
		Name:        "illimity",
		DisplayName: "Illimity",
		OwnerID:     entUser1.ID.Int64,
	}, {
		Name:        "bancoposta",
		DisplayName: "BancoPosta",
		OwnerID:     entUser1.ID.Int64,
	}, {
		Name:        "moneyfarm",
		DisplayName: "MoneyFarm",
		OwnerID:     entUser1.ID.Int64,
		IsDefault:   true,
		TypeID:      mt.AccTypeInvestment,
	}, {
		Name:        "crediti",
		DisplayName: "Crediti",
		OwnerID:     entUser1.ID.Int64,
		TypeID:      mt.AccTypeCredit,
		IsDefault:   true,
	}, {
		Name:        "cred_lavoro_m",
		DisplayName: "Crediti per Lavoro",
		OwnerID:     entUser1.ID.Int64,
		TypeID:      mt.AccTypeCredit,
	}, {
		Name:        "illimity",
		DisplayName: "Illimity",
		OwnerID:     entUser2.ID.Int64,
	}, {
		Name:        "wallet",
		DisplayName: "Wallet",
		OwnerID:     entUser2.ID.Int64,
		IsDefault:   true,
	}, {
		Name:        "cassa",
		DisplayName: "Cassa",
		OwnerID:     entUser2.ID.Int64,
	}, {
		Name:        "credits",
		DisplayName: "Crediti",
		OwnerID:     entUser2.ID.Int64,
		TypeID:      mt.AccTypeCredit,
		IsDefault:   true,
	}, {
		Name:        "moneyfarm",
		DisplayName: "MoneyFarm",
		OwnerID:     entUser2.ID.Int64,
		TypeID:      mt.AccTypeInvestment,
	}, {
		Name:        "cassa",
		DisplayName: "Cassa",
		OwnerID:     entUser3.ID.Int64,
		IsDefault:   true,
	}, {
		Name:        "credits",
		DisplayName: "Crediti",
		OwnerID:     entUser3.ID.Int64,
		TypeID:      mt.AccTypeCredit,
		IsDefault:   true,
	}}

	for _, a := range accounts {
		err = s.AddAccount(&a)
		if err != nil {
			return err
		}
	}

	// Add Categories
	categories := []string{
		"Spesa",
		"Pasti",
		"Merende",
		"Utenze",
		"Utenze/Energia",
		"Utenze/Internet",
		"Utenze/Telefonia",
		"Salute",
		"Salute/Visite",
		"Salute/Farmacia",
		"Uscite",
		"Uscite/Mangiare",
		"Sport",
		"Sport/Tennis",
		"Trasporti",
		"Trasporti/Carburante",
		"Tasse",
		"Tasse/Rifiuti",
		"Regali",
		"Viaggi",
	}

	for _, n := range categories {
		_, err = s.AddCategory(n)
		if err != nil {
			return err
		}
	}

	fmt.Println("OK")

	return nil
}
