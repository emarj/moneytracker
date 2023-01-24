package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestEntity(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	store := NewTemp()
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
	assert.True(ents[0].Name == "system")
	assert.True(ents[0].IsSystem)

	ent, err := store.GetEntity(ents[0].ID.Int64)
	require.NoError(err)
	assert.Equal(ent, ents[0])

	ent1 := mt.Entity{
		Name: "Ent1",
	}
	err = store.AddEntity(&ent1)
	require.NoError(err)
	assert.True(ent1.ID.Valid)
	assert.True(ent1.ID.Int64 > 0)
	assert.True(ent1.Name == "Ent1")
	assert.False(ent1.IsSystem)
	assert.False(ent1.IsExternal)

	ent2 := mt.Entity{
		ID:   null.IntFrom(99),
		Name: "Ent2",
	}
	err = store.AddEntity(&ent2)
	require.NoError(err)
	assert.True(ent2.ID.Valid)
	assert.EqualValues(ent2.ID.Int64, 99)
	assert.True(ent2.Name == "Ent2")

	ents, err = store.GetEntities()
	require.NoError(err)
	require.EqualValues(3, len(ents))

}

func TestGetEntityNotFound(t *testing.T) {
	require := require.New(t)

	store := NewTemp()
	err := store.Open()
	require.NoError(err)
	defer func() {
		store.Close()
	}()

	_, err = store.GetEntity(99999)
	require.ErrorIs(err, mt.ErrNotFound)

}
