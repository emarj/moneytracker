package sqlite

import (
	"errors"

	mt "github.com/emarj/moneytracker"
	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/go-jet/jet/v2/qrm"
	jet "github.com/go-jet/jet/v2/sqlite"
)

func selectUserStmtWhere(where jet.BoolExpression, includePassword bool) jet.SelectStatement {

	cols := []jet.Projection{
		jt.User.AllColumns.Except(jt.User.Password),
		jt.EntityShare.AllColumns,
		jt.Entity.AllColumns,
	}

	if includePassword {
		cols = append(cols, jt.User.Password.AS("password"))
	}

	stmt := jet.SELECT(cols[0], cols[1:]...).FROM(jt.User.LEFT_JOIN(
		jt.EntityShare,
		jt.EntityShare.UserID.EQ(jt.User.ID)).
		LEFT_JOIN(
			jt.Entity,
			jt.Entity.ID.EQ(jt.EntityShare.EntityID))).WHERE(where)

	return stmt
}

func (s *SQLiteStore) GetUserByName(user string) (mt.User, error) {
	var err error
	var u mt.User

	stmt := selectUserStmtWhere(jt.User.Name.EQ(jet.String(user)), false)

	err = stmt.Query(s.db, &u)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return u, mt.ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func (s *SQLiteStore) GetUserByID(uID int64) (mt.User, error) {
	var err error
	var u mt.User

	stmt := selectUserStmtWhere(jt.User.ID.EQ(jet.Int(uID)), false)

	err = stmt.Query(s.db, &u)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return u, mt.ErrNotFound
		}
		return u, err
	}

	return u, nil
}
