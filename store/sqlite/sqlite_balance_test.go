package sqlite

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
)

func TestAccountNoBalance(t *testing.T) {
	store := New(":memory:", true)
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	for _, tm := range times {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		require.NoError(t, err)
		require.True(t, b.Value.IsZero())
	}

}

func TestAccountZeroBalance(t *testing.T) {
	store := New(":memory:", true)
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	err = store.SetBalance(mt.Balance{
		AccountID: acc.ID,
		ValueAt: mt.ValueAt{
			Timestamp: now,
			Value:     decimal.Zero,
		},
		Delta: decimal.NewNullDecimal(decimal.Zero),
	})
	require.NoError(t, err)

	for _, tm := range times {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		require.NoError(t, err)
		require.True(t, b.Value.IsZero())
	}

}

func TestAccountBalance(t *testing.T) {
	store := New(":memory:", true)
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	val := decimal.NewFromInt(1000)

	err = store.SetBalance(mt.Balance{
		AccountID: acc.ID,
		ValueAt: mt.ValueAt{
			Timestamp: now,
			Value:     val,
		},
		Delta: decimal.NewNullDecimal(val),
	})
	require.NoError(t, err)

	for _, tm := range past {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		assert.NoError(t, err)
		assert.True(t, b.Value.IsZero(), fmt.Sprintf("Value should be zero. Got %s", b.Value))
	}

	for _, tm := range future {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		assert.NoError(t, err)
		assert.True(t, b.Value.Equal(val), fmt.Sprintf("Value should be %s. Got %s", val, b.Value))
	}
}

func TestAccountBalanceWithBalanceAndTransactions(t *testing.T) {
	store := New(":memory:", true)
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	val := decimal.NewFromInt(1000)

	err = store.SetBalance(mt.Balance{
		AccountID: acc.ID,
		ValueAt: mt.ValueAt{
			Timestamp: now,
			Value:     val,
		},
		Delta: decimal.NewNullDecimal(val),
	})
	require.NoError(t, err)

	delta := decimal.NewFromInt(500)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		TypeID:      mt.OpTypeExpense,
		Transactions: []mt.Transaction{
			{
				Timestamp: now,
				From:      acc,
				To:        mt.Account{ID: null.IntFrom(0)},
				Amount:    delta,
			},
		},
		Balances: []mt.Balance{},
	})
	require.NoError(t, err)

	for _, tm := range past {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		assert.NoError(t, err)
		assert.True(t, b.Value.IsZero(), fmt.Sprintf("Value should be zero. Got %s", b.Value))
	}

	amount := val.Sub(delta)

	for _, tm := range future {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		assert.NoError(t, err)
		assert.True(t, b.Value.Equal(amount), fmt.Sprintf("Value should be %s. Got %s", amount, b.Value))
	}
}

func TestAccountBalanceWithNoBalanceAndTransactions(t *testing.T) {
	store := New(":memory:", true)
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		Owner:       mt.Entity{ID: null.IntFrom(0)},
		TypeID:      mt.AccountMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	delta := decimal.NewFromInt(500)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		TypeID:      mt.OpTypeExpense,
		Transactions: []mt.Transaction{
			{
				Timestamp: now,
				From:      acc,
				To:        mt.Account{ID: null.IntFrom(0)},
				Amount:    delta,
			},
		},
		Balances: []mt.Balance{},
	})
	require.NoError(t, err)

	for _, tm := range past {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		assert.NoError(t, err)
		assert.True(t, b.Value.IsZero(), fmt.Sprintf("Value should be zero. Got %s", b.Value))
	}

	amount := delta.Neg()

	for _, tm := range future {
		b, err := store.GetValueAt(acc.ID.Int64, tm)
		assert.NoError(t, err)
		assert.True(t, b.Value.Equal(amount), fmt.Sprintf("Value should be %s. Got %s", amount, b.Value))
	}
}
