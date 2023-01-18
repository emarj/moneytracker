package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/emarj/moneytracker/timestamp"
	"github.com/shopspring/decimal"
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

	op := mt.Operation{
		Description: "Desc",
		TypeID:      mt.OpTypeExpense,
		Transactions: []mt.Transaction{{
			Timestamp: timestamp.Now(),
			FromID:    a.ID.Int64,
			ToID:      mt.AccWorldID,
			Amount:    decimal.NewFromInt(50),
		}},
		CategoryID: 0,
	}
	err = store.AddOperation(&op)
	require.NoError(t, err)
	require.True(t, op.ID.Valid)
	require.Equal(t, "Desc", op.Description)

}
