package sqlite

import (
	"fmt"
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/emarj/moneytracker/datetime"
	tt "github.com/emarj/moneytracker/datetime/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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
		b, err := store.GetBalanceAt(m.account.ID.Int64, pair.Key)
		assert.NoError(t, err)
		assert.True(t, b.Value.Equal(pair.Value), fmt.Sprintf("%s@%s should be %s. Got %s", m.account.Name, pair.Key.String(), pair.Value, b.Value))
	}
}

func (m ValueMap) MultiSet(lists ...[]Point) {
	for _, list := range lists {
		for _, p := range list {
			m.Set(p.T, p.V)
		}
	}
}

func NewValueMap(acc mt.Account, lists ...[]Point) ValueMap {
	vm := ValueMap{
		OrderedMap: orderedmap.New[datetime.DateTime, decimal.Decimal](),
		account:    acc,
	}

	vm.MultiSet(lists...)

	return vm
}

func PointList(v decimal.Decimal, times ...datetime.DateTime) []Point {
	list := make([]Point, len(times))
	for i := 0; i < len(times); i++ {
		list[i] = Point{times[i], v}
	}
	return list
}

///////////////////////////////////////////////////////////////////////////

func TestAccountNoBalance(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)
	assert.Equal(t, acc.OwnerID, mt.EntSystemID)

	// With no balances the balance should be zero for all times
	m := NewValueMap(acc, PointList(decimal.Zero, tt.Times...))

	m.Test(t, store)

}

func TestAccountZeroBalance(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	b := mt.Balance{
		AccountID: acc.ID,
		Timestamp: tt.Now,
		Value:     decimal.Zero,
	}
	err = store.SetBalance(&b)
	require.NoError(t, err)
	assert.True(t, b.Delta.Valid)
	assert.True(t, b.Delta.Decimal.IsZero())

	m := NewValueMap(acc, PointList(b.Value, tt.Times...))
	m.Test(t, store)

}

func TestAccountBalance(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	b := mt.Balance{
		AccountID: acc.ID,
		Timestamp: tt.Now,
		Value:     decimal.NewFromInt(1000),
	}
	err = store.SetBalance(&b)
	require.NoError(t, err)
	assert.True(t, b.Delta.Valid)
	assert.True(t, b.Delta.Decimal.Equal(b.Value))

	m := NewValueMap(acc,
		PointList(decimal.Zero, tt.Past...),
		PointList(b.Value, tt.Future...))
	m.Test(t, store)

}

func TestAccountBalanceWithTransactions(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	delta := decimal.NewFromInt(500)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		Transactions: []mt.Transaction{
			{
				Timestamp: tt.Now,
				FromID:    acc.ID.Int64,
				ToID:      0,
				Amount:    delta,
			},
		},
	})
	require.NoError(t, err)

	amount := delta.Neg()

	m := NewValueMap(acc,
		PointList(decimal.Zero, tt.Past...),
		PointList(amount, tt.Future...))
	m.Test(t, store)
}

func TestAccountBalanceWithBalanceAndTransactions(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	balNow := mt.Balance{
		AccountID: acc.ID,
		Timestamp: tt.Now,
		Value:     decimal.NewFromInt(1000),
	}
	err = store.SetBalance(&balNow)
	require.NoError(t, err)
	assert.True(t, balNow.Delta.Valid)
	assert.True(t, balNow.Delta.Decimal.Equal(balNow.Value))

	delta := decimal.NewFromInt(500)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		Transactions: []mt.Transaction{
			{
				Timestamp: tt.Now,
				FromID:    acc.ID.Int64,
				ToID:      0,
				Amount:    delta,
			},
		},
	})
	require.NoError(t, err)

	// 1000 - 500 = 500
	newValueNow := balNow.Value.Sub(delta)

	m := NewValueMap(acc,
		PointList(decimal.Zero, tt.Past...),
		PointList(newValueNow, tt.Now.Plus(tt.Epsilon)),
		PointList(newValueNow, tt.Future...),
	)
	m.Test(t, store)

	delta1 := decimal.NewFromInt(200)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		Transactions: []mt.Transaction{
			{
				Timestamp: tt.Later,
				FromID:    acc.ID.Int64,
				ToID:      0,
				Amount:    delta1,
			},
		},
	})
	require.NoError(t, err)

	// 500 - 200 = 300
	newValueLater := newValueNow.Sub(delta1)

	m = NewValueMap(acc,
		PointList(decimal.Zero, tt.Past...),
		PointList(newValueNow, tt.Now),
		PointList(newValueLater, tt.Later.Plus(tt.Epsilon), tt.LATER, tt.EndOfTime),
	)
	m.Test(t, store)

	balLATER := mt.Balance{
		AccountID: acc.ID,
		Timestamp: tt.LATER,
		Value:     decimal.NewFromInt(2000),
	}
	err = store.SetBalance(&balLATER)
	require.NoError(t, err)
	assert.True(t, balLATER.Delta.Valid)
	assert.True(t, balLATER.Delta.Decimal.Equal(balLATER.Value.Sub(newValueLater)))

	m = NewValueMap(acc, PointList(balLATER.Value, tt.LATER.Plus(tt.Epsilon), tt.EndOfTime))
	m.Test(t, store)
}

func TestAccountBalancePastTransactionDeltas(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	bNow := mt.Balance{
		AccountID: acc.ID,
		Timestamp: tt.Now,
		Value:     decimal.NewFromInt(1000),
	}
	err = store.SetBalance(&bNow)
	require.NoError(t, err)
	assert.True(t, bNow.Delta.Valid)
	assert.True(t, bNow.Delta.Decimal.Equal(bNow.Value))

	delta := decimal.NewFromInt(500)

	err = store.AddOperation(&mt.Operation{
		Description: "",
		Transactions: []mt.Transaction{
			{
				Timestamp: tt.Now.Plus(tt.Epsilon),
				FromID:    acc.ID.Int64,
				ToID:      0,
				Amount:    delta,
			},
		},
	})
	require.NoError(t, err)
	// 1000 - 500 = 500
	newValueNow := bNow.Value.Sub(delta)

	m := NewValueMap(acc,
		PointList(decimal.Zero, tt.Past...),
		PointList(newValueNow, tt.Future...),
	)
	m.Test(t, store)

	b1, err := store.GetLastBalance(acc.ID.Int64)
	require.NoError(t, err)
	// The delta should not have changed
	assert.True(t, b1.Delta.Decimal.Equal(bNow.Delta.Decimal))

	// Let us add a past transaction
	err = store.AddOperation(&mt.Operation{
		Description: "",
		Transactions: []mt.Transaction{
			{
				Timestamp: tt.BEFORE,
				FromID:    0,
				ToID:      acc.ID.Int64,
				Amount:    delta,
			},
		},
	})
	require.NoError(t, err)

	// Now only the dates between BEFORE and now (excluded) should be affected
	m.MultiSet(PointList(delta, tt.BEFORE, tt.Before))
	m.Test(t, store)

	//TODO: Fix deltas

	/* expectedDelta := bNow.Delta.Decimal.Sub(delta)

	// The delta should be updated
	b2, err := store.GetLastBalance(acc.ID.Int64)
	require.NoError(t, err)
	// The delta should be updated
	assert.True(t, b2.Delta.Decimal.Equal(expectedDelta),
		fmt.Sprintf("want %s, got %s", expectedDelta, b2.Delta.Decimal),
	) */
}

func TestAccountBalanceDelete(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	acc := mt.Account{
		Name:        "acc",
		DisplayName: "Acc",
		OwnerID:     mt.EntSystemID,
		TypeID:      mt.AccTypeMoney,
	}
	err = store.AddAccount(&acc)
	require.NoError(t, err)
	assert.True(t, acc.ID.Valid)

	b := mt.Balance{
		AccountID: acc.ID,
		Timestamp: tt.Now,
		Value:     decimal.NewFromInt(1000),
	}
	err = store.SetBalance(&b)
	require.NoError(t, err)
	assert.True(t, b.Delta.Valid)
	assert.True(t, b.Delta.Decimal.Equal(b.Value))

	m := NewValueMap(acc, PointList(decimal.Zero, tt.Past...), PointList(b.Value, tt.Future...))
	m.Test(t, store)

	err = store.DeleteBalance(acc.ID.Int64, tt.Now)
	require.NoError(t, err)

	m.MultiSet(PointList(decimal.Zero, tt.Times...))
	m.Test(t, store)

}
