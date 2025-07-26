package device

import (
	"mrav/isa"
)

type Device interface {
	Name() string
	Hit(address isa.BusValue) bool
	TickCycle()
	ReadBus(address isa.BusValue) (isa.BusValue, error)
	WriteBus(address isa.BusValue, value isa.BusValue) error
}
