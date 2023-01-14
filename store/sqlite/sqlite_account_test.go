package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/require"
)

func TestAccountCRUD(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	e := mt.Entity{
		Name:       "ent",
		IsSystem:   false,
		IsExternal: false,
	}
	err = store.AddEntity(&e)
	require.NoError(t, err)
	require.True(t, e.ID.Valid)
	require.NotEqual(t, e.ID.Int64, 0)

	a := mt.Account{
		Name:        "testacc",
		DisplayName: "Test Acc",
		OwnerID:     e.ID.Int64,
		TypeID:      mt.AccountMoney,
	}

	err = store.AddAccount(&a)
	require.NoError(t, err)
	require.True(t, a.ID.Valid)

	b, err := store.GetAccount(a.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, &a, b)

	b.Name = "newname"
	b.DisplayName = "NewName"
	b.IsSystem = true

	err = store.UpdateAccount(b)
	require.NoError(t, err)

	c, err := store.GetAccount(a.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, b, c)

	err = store.DeleteAccount(a.ID.Int64)
	require.NoError(t, err)

	_, err = store.GetAccount(a.ID.Int64)
	require.Error(t, err)
}
