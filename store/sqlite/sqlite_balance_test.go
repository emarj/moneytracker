package sqlite

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/guregu/null.v4"
	mt "ronche.se/moneytracker"
	"ronche.se/moneytracker/datetime"
)

type Point struct {
	T datetime.DateTime
	V decimal.Decimal
}

type ValueMap struct {
	account mt.Account
	*orderedmap.OrderedMap[datetime.DateTime, decimal.Decimal]
}

func (m ValueMap) Test(t *testing.T, store *SQLiteStore) {
	for pair := m.Oldest(); pair != nil; pair = pair.Next() {
		b, err := store.GetValueAt(m.account.ID.Int64, pair.Key)
		assert.NoError(t, err)
		assert.True(t, b.Value.Equal(pair.Value), fmt.Sprintf("%s@%s should be %s. Got %s", m.account.Name, pair.Key.String(), pair.Value, b.Value))
	}
}

func NewValueMap(acc mt.Account, lists ...[]Point) ValueMap {
	m := orderedmap.New[datetime.DateTime, decimal.Decimal]()

	for _, list := range lists {
		for _, p := range list {
			m.Set(p.T, p.V)
		}
	}

	return ValueMap{
		OrderedMap: m,
		account:    acc,
	}
}

func PointList(v decimal.Decimal, times []datetime.DateTime) []Point {
	list := make([]Point, len(times))
	for i := 0; i < len(times); i++ {
		list[i] = Point{times[i], v}
	}
	return list
}

///////////////////////////////////////////////////////////////////////////

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

	m := NewValueMap(acc, PointList(decimal.Zero, times))

	m.Test(t, store)

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
	})
	require.NoError(t, err)

	m := NewValueMap(acc, PointList(decimal.Zero, times))
	m.Test(t, store)

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

	m := NewValueMap(acc, PointList(decimal.Zero, past), PointList(val, future))

	m.Test(t, store)

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

	amount := delta.Neg()

	m := NewValueMap(acc, PointList(decimal.Zero, past), PointList(amount, future))
	m.Test(t, store)
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

	valueNow := decimal.NewFromInt(1000)

	err = store.SetBalance(mt.Balance{
		AccountID: acc.ID,
		ValueAt: mt.ValueAt{
			Timestamp: now,
			Value:     valueNow,
		},
		Delta: decimal.NewNullDecimal(valueNow),
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

	newValueNow := valueNow.Sub(delta)

	m := NewValueMap(acc, PointList(decimal.Zero, past), PointList(newValueNow, future))
	m.Test(t, store)

	delta1 := decimal.NewFromInt(200)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		TypeID:      mt.OpTypeExpense,
		Transactions: []mt.Transaction{
			{
				Timestamp: Later,
				From:      acc,
				To:        mt.Account{ID: null.IntFrom(0)},
				Amount:    delta1,
			},
		},
		Balances: []mt.Balance{},
	})
	require.NoError(t, err)

	newValueLater := newValueNow.Sub(delta1)

	m = NewValueMap(acc,
		PointList(decimal.Zero, past),
		PointList(newValueNow, []datetime.DateTime{now, later}),
		PointList(newValueLater, []datetime.DateTime{Later, LATER, EndOfTime}),
	)
	m.Test(t, store)

	// Let us add a past transaction
	err = store.AddOperation(&mt.Operation{
		Description: "",
		TypeID:      mt.OpTypeExpense,
		Transactions: []mt.Transaction{
			{
				Timestamp: BEFORE,
				From:      acc,
				To:        mt.Account{ID: null.IntFrom(0)},
				Amount:    delta1,
			},
		},
		Balances: []mt.Balance{},
	})
	require.NoError(t, err)

	// Now only the dates between BEFORE and now (excluded) should be affected

	m.Set(BEFORE, delta1.Neg())
	m.Set(Before, delta1.Neg())
	m.Set(before, delta1.Neg())
	/* OR
	 m = NewValueMap(acc,
		NewList(decimal.Zero, []datetime.DateTime{BeginningOfTime}),
		NewList(delta1.Neg(), []datetime.DateTime{BEFORE, Before, before}),
		NewList(newValueNow, []datetime.DateTime{now, later}),
		NewList(newValueLater, []datetime.DateTime{Later, LATER, EndOfTime}),
	) */
	m.Test(t, store)

}
