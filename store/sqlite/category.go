package sqlite

import (
	"errors"
	"fmt"
	"strings"

	mt "github.com/emarj/moneytracker"
	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/go-jet/jet/v2/qrm"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func (s *SQLiteStore) GetCategories() ([]mt.Category, error) {

	Parent := jt.Category.AS("parent")

	stmt := jet.SELECT(
		jt.Category.AllColumns,
		Parent.AllColumns,
		jet.COALESCE(Parent.Name.CONCAT(jet.String("/")).CONCAT(jt.Category.Name), jt.Category.Name).AS("category.full_name"),
	).
		FROM(
			jt.Category.LEFT_JOIN(
				Parent, Parent.ID.EQ(jt.Category.ParentID),
			),
		).ORDER_BY(jet.Raw("\"category.full_name\" ASC"))

	//fmt.Println(stmt.DebugSql())

	categories := []mt.Category{}
	err := stmt.Query(s.db, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func getCategoryByName(db TXDB, name string) (mt.Category, error) {
	var cat mt.Category
	stmt := jet.SELECT(jt.Category.AllColumns).FROM(jt.Category).WHERE(jt.Category.Name.EQ(jet.String(name)))
	err := stmt.Query(db, &cat)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func insertCategory(db TXDB, c *mt.Category) error {

	stmt := jt.Category.INSERT(jt.Category.AllColumns).MODEL(c).RETURNING(jt.Category.AllColumns)
	err := stmt.Query(db, c)
	if err != nil {
		return err
	}

	return nil
}

const MaxSubCategoryDepth int = 2

func (s *SQLiteStore) AddCategory(fullName string) (cat mt.Category, err error) {

	if fullName == "" {
		err = errors.New("category name cannot be empty")
		return
	}
	names := strings.Split(fullName, "/")
	if len(names) > MaxSubCategoryDepth {
		err = fmt.Errorf("sub-categories can have at max depth %d, got %d: %q", MaxSubCategoryDepth, len(names), fullName)
		return
	}

	var parentCat mt.Category

	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	if len(names) > 1 {
		parentCat, err = getCategoryByName(tx, names[0])
		if err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return cat, fmt.Errorf("parent category %q does not exists", names[0])
			}
			return
		}

	}

	cat = mt.Category{
		Name:     names[len(names)-1],
		ParentID: parentCat.ID,
	}

	err = insertCategory(tx, &cat)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return

}
