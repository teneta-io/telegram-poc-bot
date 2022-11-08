package entities

import "github.com/google/uuid"

type Provider struct {
	UUID uuid.UUID `gorm:"primary_key"`

	ChatID int64
	User   User

	VCPU    int64
	Ram     int64
	Storage int64
	Network int64
	Ports   string
}
