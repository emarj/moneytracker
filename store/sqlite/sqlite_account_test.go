package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestAccountCreateGet(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	a := mt.Account{
		Name:        "testacc",
		DisplayName: "Test Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}

	err = store.AddAccount(&a)
	require.NoError(t, err)
	require.True(t, a.ID.Valid)

	b, err := store.GetAccount(a.ID.Int64)
	require.NoError(t, err)

	require.Equal(t, &a, b)
}
