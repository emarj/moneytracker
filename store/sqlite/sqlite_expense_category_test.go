package sqlite

import (
	"testing"

	mt "github.com/emarj/moneytracker"
	"github.com/stretchr/testify/require"
)

func TestInsertCategory(t *testing.T) {
	store := NewTemp()
	err := store.Open()
	require.NoError(t, err)
	defer func() {
		store.Close()
	}()

	var cats []mt.Category
	var res mt.Category

	cats, err = store.GetCategories()
	require.NoError(t, err)
	require.EqualValues(t, 1, len(cats))

	res, err = store.AddCategory("Cat1")
	require.NoError(t, err)
	require.EqualValues(t, "Cat1", res.Name)

	cats, err = store.GetCategories()
	require.NoError(t, err)
	require.EqualValues(t, 2, len(cats))

	res, err = store.AddCategory("Cat2")
	require.NoError(t, err)
	require.EqualValues(t, "Cat2", res.Name)

	cats, err = store.GetCategories()
	require.NoError(t, err)
	require.EqualValues(t, 3, len(cats))

	res, err = store.AddCategory("Cat2/SubCat1")
	require.NoError(t, err)
	t.Log(res)
	require.EqualValues(t, "SubCat1", res.Name)

	cats, err = store.GetCategories()
	require.NoError(t, err)
	require.EqualValues(t, 4, len(cats))

	_, err = store.AddCategory("NonExistingCat/SubCat")
	require.Error(t, err)

	cats, err = store.GetCategories()
	require.NoError(t, err)
	require.EqualValues(t, 4, len(cats))

	_, err = store.AddCategory("Cat1/Sub/SubCat")
	require.Error(t, err)

	cats, err = store.GetCategories()
	require.NoError(t, err)
	require.EqualValues(t, 4, len(cats))

}
