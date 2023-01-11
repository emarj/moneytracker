package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mt "ronche.se/moneytracker"
)

func TestEntity(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	store := New(":memory:", true)
	err := store.Open()
	require.NoError(err)
	defer func() {
		store.Close()
	}()

	ents, err := store.GetEntities()
	require.NoError(err)
	// There is already the system entity present
	require.EqualValues(1, len(ents))
	assert.True(ents[0].ID.Valid, "system id must be valid")
	assert.True(ents[0].ID.Int64 == 0, "system id must be zero")
	assert.True(ents[0].Name == "_system")
	assert.True(ents[0].IsSystem)

	ent, err := store.GetEntity(ents[0].ID.Int64)
	require.NoError(err)
	assert.Equal(*ent, ents[0])

	*ent = mt.Entity{
		Name: "Ent1",
	}
	err = store.AddEntity(ent)
	require.NoError(err)
	assert.True(ent.ID.Valid)
	assert.True(ent.ID.Int64 > 0)
	assert.True(ent.Name == "Ent1")
	assert.False(ent.IsSystem)
	assert.False(ent.IsExternal)

	ents, err = store.GetEntities()
	require.NoError(err)
	require.EqualValues(2, len(ents))

}

func TestGetEntityNotFound(t *testing.T) {
	require := require.New(t)

	store := New(":memory:", true)
	err := store.Open()
	require.NoError(err)
	defer func() {
		store.Close()
	}()

	_, err = store.GetEntity(99999)
	require.ErrorIs(err, mt.ErrNotFound)

}
