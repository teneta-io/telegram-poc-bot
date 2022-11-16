package entities

import (
	"errors"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

var AvailableProtocols = []string{"TCP", "UDP"}

var ErrProtocolIsNotSupported = errors.New("protocol is not supported")
var ErrNumberCanNotBePort = errors.New("number can not be port")
var ErrWrongPortFormat = errors.New("wrong ports format")

type Provider struct {
	UUID   uuid.UUID `gorm:"primary_key"`
	ChatID int64

	UserUUID uuid.UUID
	User     User

	VCPU    int64 `gorm:"column:vcpu"`
	Ram     int64
	Storage int64
	Network int64
	Ports   Ports `gorm:"serializer:json"`
}

type Ports []Port

func (ports *Ports) String() string {
	return strings.Join(lo.Map(*ports, func(item Port, index int) string {
		return item.String()
	}), ", ")
}

type Port struct {
	Protocol string
	Number   int
}

func (port *Port) String() string {
	number := strconv.Itoa(port.Number)
	return port.Protocol + ":" + number
}

func (provider *Provider) SetPorts(ports []string) error {
	newPorts := make([]Port, len(ports))

	for i, port := range ports {
		couple := strings.Split(port, ":")
		if len(couple) != 2 {
			return ErrWrongPortFormat
		}

		protocol, numberStr := strings.ToUpper(couple[0]), couple[1]

		if !lo.Contains(AvailableProtocols, protocol) {
			return ErrProtocolIsNotSupported
		}

		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return ErrNumberCanNotBePort
		}

		if number < 0 || number > (1<<16) {
			return ErrNumberCanNotBePort
		}

		newPorts[i] = Port{Protocol: protocol, Number: number}
	}

	provider.Ports = append(provider.Ports, newPorts...)
	provider.Ports = lo.Uniq(provider.Ports)

	return nil
}
