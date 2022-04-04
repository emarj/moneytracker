package mock

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"ronche.se/moneytracker/domain"
)

type mockTransactionStore struct {
	transactions map[string]*domain.Transaction
}

func newMockTransactionStore() *mockTransactionStore {
	return &mockTransactionStore{
		transactions: map[string]*domain.Transaction{},
	}
}

func (ts *mockTransactionStore) GetTransaction(tID uuid.UUID) (*domain.Transaction, error) {
	idstr := tID.String()
	t, ok := ts.transactions[idstr]
	if !ok {
		return nil, fmt.Errorf("an transaction with uuid=%s does not exists", tID.String())
	}
	return t, nil
}

func (ts *mockTransactionStore) GetTransactionsByAccount(aID uuid.UUID) ([]*domain.Transaction, error) {
	result := []*domain.Transaction{}
	for _, t := range ts.transactions {
		if t.FromID == aID || t.ToID == aID {
			result = append(result, t)
		}
	}

	return result, nil
}

func (ts *mockTransactionStore) GetTransactionsByUser(uID string) ([]*domain.Transaction, error) {

	tm := map[string]*domain.Transaction{}
	for _, t := range ts.transactions {
		if t.OwnerID == uID {
			tm[t.ID.String()] = t
		}
	}

	result := make([]*domain.Transaction, len(tm))
	i := 0
	for _, t := range tm {
		result[i] = t
		i++
	}

	return result, nil
}

func (ts *mockTransactionStore) AddTransaction(t *domain.Transaction) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	t.ID = id
	t.DateCreated = time.Now()
	t.DateModified = time.Now()

	ts.transactions[t.ID.String()] = t

	return nil

}

func (ts *mockTransactionStore) DeleteTransaction(tID uuid.UUID) error {
	idstr := tID.String()
	_, ok := ts.transactions[idstr]
	if !ok {
		return fmt.Errorf("an transaction with uuid=%s does not exists", idstr)
	}
	delete(ts.transactions, idstr)

	return nil

}

func (ts *mockTransactionStore) UpdateTransaction(t *domain.Transaction) error {
	idstr := t.ID.String()
	_, ok := ts.transactions[idstr]
	if !ok {
		return fmt.Errorf("an transaction with uuid=%s does not exists", idstr)
	}

	t.DateModified = time.Now()
	ts.transactions[idstr] = t

	return nil

}
