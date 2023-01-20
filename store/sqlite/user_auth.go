package sqlite

import (
	"errors"

	mt "github.com/emarj/moneytracker"
	jt "github.com/emarj/moneytracker/.gen/table"
	"github.com/go-jet/jet/v2/qrm"
	jet "github.com/go-jet/jet/v2/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type userWithHashedPassword struct {
	mt.User
	Password []byte
}

func (s *SQLiteStore) Login(user string, password string) (mt.User, error) {
	var err error
	var u mt.User

	stmt := jt.User.LEFT_JOIN(
		jt.EntityShare,
		jt.EntityShare.UserID.EQ(jt.User.ID)).
		LEFT_JOIN(
			jt.Entity,
			jt.Entity.ID.EQ(jt.EntityShare.EntityID)).
		SELECT(
			jt.User.AllColumns.Except(jt.User.Password),
			jt.User.Password.AS("password"),
			jt.EntityShare.AllColumns,
			jt.Entity.AllColumns,
		).WHERE(
		jt.User.Name.EQ(jet.String(user)),
	)

	result := userWithHashedPassword{}
	err = stmt.Query(s.db, &result)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return u, mt.ErrNotFound
		}
		return u, err
	}

	err = bcrypt.CompareHashAndPassword(result.Password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return u, mt.ErrUnauthorized
		}
		return u, err
	}

	u = result.User
	return u, nil
}

func (s *SQLiteStore) RegisterUser(user *mt.User, password string) error {
	var err error
	u := *user

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	payload := userWithHashedPassword{u, hashedPassword}

	err = jt.User.INSERT(jt.User.AllColumns).MODEL(&payload).RETURNING(jt.User.AllColumns).
		Query(s.db, &u)
	if err != nil {
		return err
	}

	*user = u

	return nil
}
