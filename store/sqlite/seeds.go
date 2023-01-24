package sqlite

import (
	"fmt"

	mt "github.com/emarj/moneytracker"

	jt "github.com/emarj/moneytracker/.gen/table"
)

func (s *SQLiteStore) Seeding() error {
	var err error

	fmt.Print("Seeding...")

	s.accountTypes = mt.AccountTypes()
	_, err = jt.AccountType.INSERT(jt.AccountType.AllColumns).
		MODELS(&s.accountTypes).
		ON_CONFLICT().DO_NOTHING().
		Exec(s.db)
	if err != nil {
		return err
	}

	s.operationTypes = mt.OperationTypes()
	_, err = jt.OperationType.INSERT(jt.OperationType.AllColumns).
		MODELS(&s.operationTypes).
		ON_CONFLICT().DO_NOTHING().
		Exec(s.db)
	if err != nil {
		return err
	}

	categories := mt.SystemCategories()
	_, err = jt.Category.INSERT(jt.Category.AllColumns).
		MODELS(&categories).
		ON_CONFLICT().DO_NOTHING().
		Exec(s.db)
	if err != nil {
		return err
	}

	ents := mt.SystemEntities()
	_, err = jt.Entity.INSERT(jt.Entity.AllColumns).
		MODELS(&ents).
		ON_CONFLICT().DO_NOTHING().
		Exec(s.db)
	if err != nil {
		return err
	}

	accs := mt.SystemAccounts()
	_, err = jt.Account.INSERT(jt.Account.AllColumns).
		MODELS(&accs).
		ON_CONFLICT().DO_NOTHING().
		Exec(s.db)
	if err != nil {
		return err
	}

	fmt.Println("OK")

	return nil

}
