package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/require"
)

func TestOperationCRUD(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	a := mt.Account{
		Name:        "testacc",
		DisplayName: "Test Acc",
		OwnerID:     0,
		TypeID:      mt.AccTypeMoney,
	}

	err = store.AddAccount(&a)
	require.NoError(t, err)
	require.True(t, a.ID.Valid)

}
