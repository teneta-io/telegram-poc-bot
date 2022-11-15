package entities

import "github.com/google/uuid"

type Provider struct {
	UUID   uuid.UUID `gorm:"primary_key"`
	ChatID int64

	UserUUID uuid.UUID
	User     User

	VCPU    int64 `gorm:"column:vcpu"`
	Ram     int64
	Storage int64
	Network int64
	Ports   []Port `gorm:"serializer:json"`
}

type Port struct {
	Protocol string
	Number   int
}
