package mockup

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"ronche.se/expensetracker/model"
)

func New() (model.Service, error) {
	return &mockup{make(map[string]*model.Expense),
		make(map[int]*model.Category),
		make(map[int]*model.User),
		make(map[int]*model.PaymentMethod),
		0,
		0,
		0,
	}, nil

}

type mockup struct {
	expenses   map[string]*model.Expense
	categories map[int]*model.Category
	users      map[int]*model.User
	pmethods   map[int]*model.PaymentMethod
	catid      int
	uid        int
	pmid       int
}

func (m *mockup) ExpenseInsert(expense *model.Expense) (*model.Expense, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	expense.UUID = id
	expense.DateCreated = time.Now().Local()

	if expense.Date.IsZero() {
		expense.Date = expense.DateCreated
	}

	m.expenses[id.String()] = expense
	return expense, nil
}
func (m *mockup) ExpenseGet(id uuid.UUID) (*model.Expense, error) {
	idstr := id.String()
	e := m.expenses[idstr]
	if e == nil {
		return nil, fmt.Errorf("an expense with uuid=%s does not exists", id.String())
	}
	return e, nil
}

func (m *mockup) ExpenseUpdate(expense *model.Expense) (*model.Expense, error) {
	idstr := expense.UUID.String()
	if m.expenses[idstr] == nil {
		return nil, fmt.Errorf("an expense with uuid=%s does not exists", expense.UUID.String())
	}
	m.expenses[idstr] = expense

	return expense, nil
}
func (m *mockup) ExpensesGetN(limit int) ([]*model.Expense, error) {
	es := make([]*model.Expense, len(m.expenses))
	i := 0
	for _, value := range m.expenses {
		es[i] = value
		es[i].Category = m.categories[es[i].Category.ID]
		es[i].Who = m.users[es[i].Who.ID]
		es[i].Method = m.pmethods[es[i].Method.ID]
		i++
	}
	return es, nil
}
func (m *mockup) ExpenseDelete(id uuid.UUID) error {
	idstr := id.String()
	m.expenses[idstr] = nil
	return nil
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
	cat := model.Category{m.catid, name}
	m.categories[cat.ID] = &cat
	m.catid++
	return &cat, nil
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
	u := model.User{m.uid, name}
	m.users[u.ID] = &u
	m.uid++
	return &u, nil
}

func (m *mockup) PaymentMethodsGetAll() ([]*model.PaymentMethod, error) {
	pm := make([]*model.PaymentMethod, len(m.pmethods))
	i := 0
	for _, value := range m.pmethods {
		pm[i] = value
		i++
	}
	return pm, nil
}

func (m *mockup) PaymentMethodInsert(name string) (*model.PaymentMethod, error) {
	pm := model.PaymentMethod{m.pmid, name}
	m.pmethods[pm.ID] = &pm
	m.pmid++
	return &pm, nil
}

func (m *mockup) Close() error {
	//I'm fake so I do nothing
	return nil
}
