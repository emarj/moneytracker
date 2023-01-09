package sqlite

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func TestAccountsCreateGet(t *testing.T) {
	store := New(":memory:", true)
	err := store.Open()
	require.NoError(t, err)
	defer store.Close()

	a := mt.Account{
		Name:        "testacc",
		DisplayName: "Test Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}

	err = store.AddAccount(&a, nil)
	require.NoError(t, err)
	require.Equal(t, a.ID.Valid, true)

	b, err := store.GetAccount(int(a.ID.Int64))
	require.NoError(t, err)

	require.Equal(t, &a, b)
}
