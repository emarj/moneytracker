package mock

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"ronche.se/moneytracker/domain"
)

func (m *MockStore) GetTransaction(tID uuid.UUID) (*domain.Transaction, error) {
	idstr := tID.String()
	t, ok := m.transactions[idstr]
	if !ok {
		return nil, fmt.Errorf("an transaction with uuid=%s does not exists", tID.String())
	}
	return t, nil
}

func (m *MockStore) GetTransactionsByAccount(aID string) ([]*domain.Transaction, error) {
	result := []*domain.Transaction{}
	for _, t := range m.transactions {
		if t.FromID == aID || t.ToID == aID {
			result = append(result, t)
		}
	}

	return result, nil
}

func (m *MockStore) GetTransactionsByUser(uID string) ([]*domain.Transaction, error) {
	al, err := m.GetAccountsOfUser(uID)
	if err != nil {
		return nil, nil
	}

	tm := map[string]*domain.Transaction{}
	for _, a := range al {
		for _, t := range m.transactions {
			if t.FromID == a.ID() || t.ToID == a.ID() {
				tm[t.ID.String()] = t
			}
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

func (m *MockStore) InsertTransaction(t *domain.Transaction) (uuid.UUID, error) {
	id, err := uuid.NewV1()
	if err != nil {
		return uuid.Nil, err
	}
	t.ID = id
	t.DateCreated = domain.DateTime{time.Now()}
	t.DateModified = domain.DateTime{time.Now()}

	m.transactions[t.ID.String()] = t
	m.Balance(t.FromID, t.Amount.Neg())
	m.Balance(t.ToID, t.Amount)

	return id, nil

}

func (m *MockStore) DeleteTransaction(tID uuid.UUID) error {
	idstr := tID.String()
	t, ok := m.transactions[idstr]
	if !ok {
		return fmt.Errorf("an transaction with uuid=%s does not exists", idstr)
	}

	m.Balance(t.FromID, t.Amount)
	m.Balance(t.ToID, t.Amount.Neg())
	delete(m.transactions, idstr)

	return nil

}

func (m *MockStore) UpdateTransaction(t *domain.Transaction) (uuid.UUID, error) {
	idstr := t.ID.String()
	tOld, ok := m.transactions[idstr]
	if !ok {
		return uuid.Nil, fmt.Errorf("an transaction with uuid=%s does not exists", idstr)
	}

	delta := t.Amount.Sub(tOld.Amount)

	if !delta.IsZero() {
		m.Balance(t.FromID, delta)
		m.Balance(t.ToID, delta.Neg())
	}
	t.DateModified = domain.DateTime{time.Now()}
	m.transactions[idstr] = t

	return t.ID, nil

}
