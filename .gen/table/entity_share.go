//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var EntityShare = newEntityShareTable("", "entity_share", "")

type entityShareTable struct {
	sqlite.Table

	//Columns
	UserID   sqlite.ColumnInteger
	EntityID sqlite.ColumnInteger
	Quota    sqlite.ColumnInteger

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type EntityShareTable struct {
	entityShareTable

	EXCLUDED entityShareTable
}

// AS creates new EntityShareTable with assigned alias
func (a EntityShareTable) AS(alias string) *EntityShareTable {
	return newEntityShareTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new EntityShareTable with assigned schema name
func (a EntityShareTable) FromSchema(schemaName string) *EntityShareTable {
	return newEntityShareTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new EntityShareTable with assigned table prefix
func (a EntityShareTable) WithPrefix(prefix string) *EntityShareTable {
	return newEntityShareTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new EntityShareTable with assigned table suffix
func (a EntityShareTable) WithSuffix(suffix string) *EntityShareTable {
	return newEntityShareTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newEntityShareTable(schemaName, tableName, alias string) *EntityShareTable {
	return &EntityShareTable{
		entityShareTable: newEntityShareTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newEntityShareTableImpl("", "excluded", ""),
	}
}

func newEntityShareTableImpl(schemaName, tableName, alias string) entityShareTable {
	var (
		UserIDColumn   = sqlite.IntegerColumn("user_id")
		EntityIDColumn = sqlite.IntegerColumn("entity_id")
		QuotaColumn    = sqlite.IntegerColumn("quota")
		allColumns     = sqlite.ColumnList{UserIDColumn, EntityIDColumn, QuotaColumn}
		mutableColumns = sqlite.ColumnList{QuotaColumn}
	)

	return entityShareTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:   UserIDColumn,
		EntityID: EntityIDColumn,
		Quota:    QuotaColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
