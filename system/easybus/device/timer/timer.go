package timer

import (
	"fmt"
	"mrav/isa"
)

type Timer struct {
	status  isa.Register
	counter isa.Register
}

func (t *Timer) Name() string {
	return "Timer"
}

const (
	cAddressLo isa.Register = 253
	cAddressHi isa.Register = 255

	cCounterReg isa.Register = 253
	cControlReg isa.Register = 254
	cStatusReg  isa.Register = 255
)

func (t *Timer) Hit(address isa.BusValue) bool {
	return (isa.Register(address) >= cAddressLo) && (isa.Register(address) <= cAddressHi)
}

func (t *Timer) ReadBus(address isa.BusValue) (isa.BusValue, error) {
	switch isa.Register(address) {
	case cCounterReg:
		return isa.BusValue(t.counter), nil
	case cControlReg:
		return isa.BusValue(0x0000), nil
	case cStatusReg:
		return isa.BusValue(t.status), nil
	}

	return 0, fmt.Errorf("device %s, reading, address out of bounds: %04X", t.Name(), address)
}

func (t *Timer) WriteBus(address isa.BusValue, value isa.BusValue) error {
	switch isa.Register(address) {
	case cCounterReg:
		t.counter = isa.Register(value)
		return nil
	case cControlReg:
		startRunning := value & 0x01

		if startRunning == 1 {
			if (t.status & 0x01) == 1 {
				// Timer is running already, nothing to do.
				return nil
			}

			t.status |= 0x01
			return nil
		}

		// Other bits don't matter.
		return nil
	case cStatusReg:
		return nil // Do nothing
	}

	return fmt.Errorf("device %s, writing, address out of bounds: %04X", t.Name(), address)
}

func (t *Timer) TickCycle() {
	if (t.status & 0x01) == 0 {
		return // Not running
	}

	t.counter--

	if t.counter == 0 {
		t.status &= ^isa.Register(0x01) // Flip the bit
	}
}
