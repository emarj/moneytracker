package mockup

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"ronche.se/moneytracker/model"
)

type mockup struct {
	transactions map[string]*model.Transaction
	types        map[int]*model.Type
	users        map[int]*model.User
	categories   map[int]*model.Category
	pmethods     map[int]*model.Method
}

func New() (model.Service, error) {
	return &mockup{make(map[string]*model.Transaction),
		make(map[int]*model.Type),
		make(map[int]*model.User),
		make(map[int]*model.Category),
		make(map[int]*model.Method),
	}, nil

}

func (m *mockup) TransactionGet(id uuid.UUID) (*model.Transaction, error) {
	idstr := id.String()
	t := m.transactions[idstr]
	if t == nil {
		return nil, fmt.Errorf("an transaction with uuid=%s does not exists", id.String())
	}
	return t, nil
}

func (m *mockup) TransactionInsert(t *model.Transaction) error {
	id := uuid.NewV4()

	t.UUID = id

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return err
	}
	t.DateCreated.Time = time.Now().In(loc)

	//Actual inser
	m.transactions[id.String()] = t

	return nil
}

func (m *mockup) TransactionUpdate(t *model.Transaction) error {
	idstr := t.UUID.String()
	if m.transactions[idstr] == nil {
		return fmt.Errorf("an transaction with uuid=%s does not exists", t.UUID.String())
	}
	m.transactions[idstr] = t

	return nil
}

func (m *mockup) TransactionDelete(id uuid.UUID) error {
	idstr := id.String()
	m.transactions[idstr] = nil
	return nil
}

func (m *mockup) TransactionsGetNOrderBy(limit int, orderBy string) ([]*model.Transaction, error) {
	ts := make([]*model.Transaction, len(m.transactions))
	i := 0
	for _, value := range m.transactions {
		ts[i] = value
		i++
	}

	return ts, nil
}

func (m *mockup) TransactionsGetNOrderByDate(limit int) ([]*model.Transaction, error) {
	return m.TransactionsGetNOrderBy(limit, "")
}

func (m *mockup) TransactionsGetNOrderByInserted(limit int) ([]*model.Transaction, error) {
	return m.TransactionsGetNOrderBy(limit, "")
}

func (m *mockup) TypesGetAll() ([]*model.Type, error) {
	tps := make([]*model.Type, len(m.types))
	i := 0
	for _, value := range m.types {
		tps[i] = value
		i++
	}
	return tps, nil
}
func (m *mockup) TypeInsert(name string) (*model.Type, error) {
	id := len(m.types)
	t := model.Type{id, name}
	m.types[id] = &t
	return &t, nil
}

func (m *mockup) UsersGetAll() ([]*model.User, error) {
	us := make([]*model.User, len(m.users))
	i := 0
	for _, value := range m.users {
		us[i] = value
		i++
	}
	return us, nil
}
func (m *mockup) UserInsert(name string) (*model.User, error) {
	id := len(m.users)
	u := model.User{id, name}
	m.users[id] = &u
	return &u, nil
}
func (m *mockup) CategoriesGetAll() ([]*model.Category, error) {
	cat := make([]*model.Category, len(m.categories))
	i := 0
	for _, value := range m.categories {
		cat[i] = value
		i++
	}
	return cat, nil
}

func (m *mockup) CategoryInsert(name string) (*model.Category, error) {
	id := len(m.categories)
	cat := model.Category{id, name}
	m.categories[id] = &cat
	return &cat, nil
}

func (m *mockup) PaymentMethodsGetAll() ([]*model.Method, error) {
	pm := make([]*model.Method, len(m.pmethods))
	i := 0
	for _, value := range m.pmethods {
		pm[i] = value
		i++
	}
	return pm, nil
}

func (m *mockup) PaymentMethodInsert(name string) (*model.Method, error) {
	id := len(m.pmethods)
	pm := model.Method{id, name}
	m.pmethods[id] = &pm
	return &pm, nil
}

func (m *mockup) Close() error {
	//I'm fake so I do nothing
	return nil
}
