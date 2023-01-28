package sqlite

import (
	"fmt"
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestEntityShare(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	e1 := mt.Entity{
		Name: "ent1",
	}
	err = store.AddEntity(&e1)
	require.NoError(t, err)

	u := mt.User{
		Name:        "user",
		DisplayName: "User",
	}
	err = store.AddUser(&u, "sdsd")
	require.NoError(t, err)

	err = store.AddSharesForEntity(mt.EntityShare{
		UserID:   u.ID.Int64,
		EntityID: e1.ID,
		Quota:    100,
		Priority: null.IntFrom(23),
	})
	require.NoError(t, err)

	e1b, err := store.GetEntity(e1.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, 1, len(e1b.Shares))
	assert.EqualValues(t, 100, e1b.Shares[0].Quota)
	assert.EqualValues(t, 23, e1b.Shares[0].Priority.Int64)

	uws, err := store.GetUserWithShares(u.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, 1, len(uws.Shares))

	e2 := mt.Entity{
		Name: "ent2",
	}
	err = store.AddEntity(&e2)
	require.NoError(t, err)

	err = store.AddSharesForEntity(mt.EntityShare{
		UserID:   u.ID.Int64,
		EntityID: e2.ID,
		Quota:    100,
	})
	require.NoError(t, err)
}

func TestEntityShare2(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	u := mt.User{
		Name:        "user",
		DisplayName: "User",
	}
	err = store.AddUser(&u, "sdsd")
	require.NoError(t, err)

	e1 := mt.Entity{
		Name: "ent1",
		Shares: []mt.EntityShare{{
			UserID:   u.ID.Int64,
			Quota:    100,
			Priority: null.IntFrom(23),
		}},
	}
	err = store.AddEntity(&e1)
	require.NoError(t, err)

	e1b, err := store.GetEntity(e1.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, 1, len(e1b.Shares))
	assert.EqualValues(t, 100, e1b.Shares[0].Quota)
	assert.EqualValues(t, 23, e1b.Shares[0].Priority.Int64)

	uws, err := store.GetUserWithShares(u.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, 1, len(uws.Shares))

	e2 := mt.Entity{
		Name: "ent2",
		Shares: []mt.EntityShare{{
			UserID: u.ID.Int64,
			Quota:  100,
		}},
	}
	err = store.AddEntity(&e2)
	require.NoError(t, err)

}

func TestEntityShare2Users(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	e := mt.Entity{
		Name: "ent1",
	}
	err = store.AddEntity(&e)
	require.NoError(t, err)

	u1 := mt.User{
		Name:        "user1",
		DisplayName: "User1",
	}
	err = store.AddUser(&u1, "sdsd")
	require.NoError(t, err)
	u2 := mt.User{
		Name:        "user2",
		DisplayName: "User2",
	}
	err = store.AddUser(&u2, "sdsd")
	require.NoError(t, err)

	err = store.AddSharesForEntity(
		[]mt.EntityShare{{
			UserID:   u1.ID.Int64,
			EntityID: e.ID,
			Quota:    70,
		}, {
			UserID:   u2.ID.Int64,
			EntityID: e.ID,
			Quota:    30,
		}}...)
	require.NoError(t, err)

	e2, err := store.GetEntity(e.ID.Int64)
	require.NoError(t, err)
	fmt.Println(e2.Shares)
	require.Equal(t, 2, len(e2.Shares))

}
