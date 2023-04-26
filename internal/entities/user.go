package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ChatID    int64
	FirstName string
	LastName  string
	WalletID  *uuid.UUID

	ProviderConfig *Provider

	Language string
	State    int
}

func (u *User) IsProvider() bool {
	return u.ProviderConfig != nil
}

func (u *User) SetPorts(ports []string) map[string]error {
	return u.ProviderConfig.SetPorts(ports)
}
