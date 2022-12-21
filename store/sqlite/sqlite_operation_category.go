package sqlite

import (
	"fmt"

	jet "github.com/go-jet/jet/v2/sqlite"
	mt "ronche.se/moneytracker"
	jt "ronche.se/moneytracker/.gen/table"
)

func (s *SQLiteStore) GetCategories() ([]mt.Category, error) {

	stmt := jet.SELECT(jt.Category.AllColumns).FROM(jt.Category)

	categories := []mt.Category{}
	err := stmt.Query(s.db, &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *SQLiteStore) AddCategory(c *mt.Category) error {

	stmt := jt.Category.INSERT(jt.Category.AllColumns).MODEL(&c).RETURNING(jt.Category.AllColumns)
	err := stmt.Query(s.db, c)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}

	return nil

}
