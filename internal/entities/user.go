package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	UUID      uuid.UUID `gorm:"primary_key"`
	ChatID    int64
	FirstName string
	LastName  string
	WalletID  *uuid.UUID

	ProviderConfig *Provider `gorm:"references:UUID"`

	Language string
	State    int
}

func (u *User) IsProvider() bool {
	return u.ProviderConfig != nil
}
