package sql

import (
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ronche.se/moneytracker/domain"
)

type SQLStore struct {
	db  *gorm.DB
	dsn string
}

func NewStore(dsn string) (*SQLStore, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &SQLStore{
		db:  db,
		dsn: dsn,
	}, nil
}

func NewInMemoryStore() (*SQLStore, error) {
	return NewStore("file::memory:?cache=shared")
}

func (s *SQLStore) Migrate() error {
	return s.db.AutoMigrate(&domain.User{}, &domain.Account{}, &domain.Transaction{}, &domain.Share{})
}

func (s *SQLStore) CreateAccount() {
	id := uuid.Must(uuid.NewV4())

	s.db.Create(domain.Account{
		ID:          id,
		Owners:      []domain.User{{"marco", "sdsd"}},
		Name:        "primary",
		DisplayName: "Primary",
		Description: "",
		Default:     false,
	})

	a := domain.Account{}

	s.db.Preload("Owners").First(&a, "ID = ?", id)

	fmt.Println(a)
}

/*
func (s *SQLStore) CreateTransaction() {
	err := s.db.AutoMigrate(&User{}, &Account{}, &Transaction{}, &Share{})
	if err != nil {
		fmt.Println(err)
	}

	//s.db.Create(User{"marco", "Marco"})

	id := uuid.Must(uuid.NewV4())
	fid := uuid.Must(uuid.NewV4())
	tid := uuid.Must(uuid.NewV4())

	t := Transaction{
		ID:           id,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		OwnerID:      "marco",
		Owner: User{
			ID:   "sds",
			Name: "FFSFSF",
		},
		Date:        time.Now(),
		Description: "fdfdf",
		Notes:       "",
		Amount:      decimal.NewFromInt(32),
		From: Account{
			ID:          fid,
			Owners:      []User{{"marco", "Marco"}},
			Name:        "primary",
			DisplayName: "Primary",
			Description: "",
			Default:     false,
		},
		ToID: tid,
		To: Account{
			ID:          tid,
			Owners:      []User{{"marco", "Marco"}},
			Name:        "secondary",
			DisplayName: "Secondary",
			Description: "",
			Default:     false,
		},
		Related: []Transaction{},
		Shared:  false,
		Shares: []Share{
			{
				ID:            uuid.Must(uuid.NewV4()),
				TransactionID: id,
				OwnerID:       "marco",
				/*Owner: User{
					ID:   "marco",
					Name: "Marco",
				},
				//OtherUserID   string
				OtherUser: domain.User{
					ID:   "arianna",
					Name: "Arianna",
				},
				Amount: decimal.NewFromInt(22),
				//AlreadyPaid   bool //if true, this is only for the user
			},
		},
		PaymentMethod: "",
		GeoLocation:   "",
		Receipt:       "",
	}

	s.db.Create(&t)

	s.db.Preload("Owner").First(&t, "ID = ?", t.ID)

	//fmt.Println(t)
}*/

//User
func (s *SQLStore) GetUsers() ([]domain.User, error) {
	var users []domain.User
	res := s.db.Find(&users)
	if res.Error != nil {
		return users, res.Error
	}

	return users, nil

}

func (s *SQLStore) GetUser(uID string) (domain.User, error) {
	var user domain.User
	if err := s.db.Find(&user, "id = ?", uID).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (s *SQLStore) AddUser(u *domain.User) error {
	if err := s.db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

func (s *SQLStore) GetAccount(aID uuid.UUID) (domain.Account, error) {
	var a domain.Account
	if err := s.db.Find(&a, "id = ?", aID).Error; err != nil {
		return a, err
	}

	return a, nil
}

/*func (s *SQLStore) GetAccountsByUser(uID string) ([]domain.Account, error) {
	var al []domain.Account
	if err := s.db.Preload("users").Find(&al, "user = ?", uID).Error; err != nil {///ALMOST SURELY WRONG
		return al, err
	}

	return al, nil
}/*

func (s *SQLStore) AddAccount(a *domain.Account) error {
	if err := s.db.Create(a).Error; err != nil {
		return err
	}

	return nil
}

/*
func (s *SQLStore) GetAccount(aID uuid.UUID) (*domain.Account, error) {

}
func (s *SQLStore) GetAccountsByUser(uID string) ([]*domain.Account, error) {

}
func (s *SQLStore) GetAccountsByUserAndName(uID string, name string) ([]*domain.Account, error) {

}
func (s *SQLStore) AddAccount(a *domain.Account) error {

}
func (s *SQLStore) GetTransaction(tID uuid.UUID) (*domain.Transaction, error) {

}
func (s *SQLStore) GetTransactionsByAccount(aID uuid.UUID) ([]*domain.Transaction, error) {

}
func (s *SQLStore) GetTransactionsByUser(uID string) ([]*domain.Transaction, error) {

}
func (s *SQLStore) AddTransaction(t *domain.Transaction) error {

}
func (s *SQLStore) UpdateTransaction(t *domain.Transaction) error {

}
func (s *SQLStore) DeleteTransaction(tID uuid.UUID) error {

}*/
