package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegister(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	user1 := mt.User{
		Name:        "user1",
		DisplayName: "User1",
		IsAdmin:     false,
	}
	user1Copy := user1

	password1 := "password1"
	err = store.RegisterUser(&user1, password1)
	require.NoError(t, err)
	assert.True(t, user1.ID.Valid)
	user1Copy.ID = user1.ID
	assert.Equal(t, user1, user1Copy)

}

func TestUserLoginNonExistingUser(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	_, err = store.Login("nonexistinguser", "password")
	require.ErrorIs(t, err, mt.ErrNotFound)

}

func TestUserLoginWithExistingUser(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	user1 := mt.User{
		Name:        "user1",
		DisplayName: "User1",
		IsAdmin:     false,
	}
	password1 := "password1"
	err = store.RegisterUser(&user1, password1)
	require.NoError(t, err)

	user1Returned, err := store.Login(user1.Name, password1)
	require.NoError(t, err)
	assert.Equal(t, user1, user1Returned)

	// Wrong password
	_, err = store.Login(user1.Name, password1+password1)
	require.Error(t, err)
	require.ErrorIs(t, err, mt.ErrUnauthorized)

}
