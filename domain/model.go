package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Balance struct {
	At      DateTime
	Balance decimal.Decimal
}
type User struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

type Account struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Owners      []User    `gorm:"many2many:account_users;"`
	Name        string
	DisplayName string
	Description string
	Default     bool
	//ParentID    uuid.UUID
}

// Transaction represent a transfer between 2 accounts
type Transaction struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	DateCreated  time.Time
	DateModified time.Time

	OwnerID string
	Owner   User

	Date time.Time

	Description string
	Notes       string
	Amount      decimal.Decimal

	FromID uuid.UUID
	From   Account
	ToID   uuid.UUID
	To     Account `gorm:"foreignKey:ToID"`

	Related []Transaction `gorm:"many2many:transaction_transactions;"`

	Shared bool
	Shares []Share

	//Details
	PaymentMethod string
	GeoLocation   string
	Receipt       string
}

type ShareAccount struct {
	ID          uuid.UUID
	OwnersID    []string
	Name        string
	DisplayName string
	Description string
	Default     bool
	ParentID    uuid.UUID
}

type Share struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
	OwnerID       string
	Owner         User
	OtherUserID   string
	OtherUser     User
	Amount        decimal.Decimal
	AlreadyPaid   bool //if true, this is only for the user
}
