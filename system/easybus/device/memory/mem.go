package memory

import (
	"fmt"

	"mrav/isa"
)

type Mem struct {
	ram []byte
}

func NewMem(size int, image []byte) (*Mem, error) {
	ram := make([]byte, size)

	if image != nil {
		if len(image) > size {
			return nil, fmt.Errorf("unable to store image of size %d bytes into RAM of %d bytes", len(image), size)
		}

		copy(ram, image)
	}

	return &Mem{
		ram: ram,
	}, nil
}

func (m *Mem) Name() string {
	return "Memory"
}

func (m *Mem) TickCycle() {} // Nothing to do

func (m *Mem) Hit(address isa.BusValue) bool {
	return address < isa.BusValue(len(m.ram))
}

func (m *Mem) ReadBus(address isa.BusValue) (isa.BusValue, error) {
	if address >= isa.BusValue(len(m.ram)-1) {
		return 0, fmt.Errorf("%s device address %X out of bounds", m.Name(), address)
	}

	hiByte := m.ram[address]
	loByte := m.ram[address+1]

	return isa.BusValue(uint16(uint16(hiByte)<<8) | uint16(loByte)), nil
}

func (m *Mem) WriteBus(address isa.BusValue, value isa.BusValue) error {
	if address >= isa.BusValue(len(m.ram)-1) {
		return fmt.Errorf("%s device address %X out of bounds", m.Name(), address)
	}

	hiByte := byte((value >> 8) & 0xFF)
	loByte := byte(value & 0xFF)

	m.ram[address] = hiByte
	m.ram[address+1] = loByte

	return nil
}

func (m *Mem) GetMemoryBytes() []byte {
	memCopy := make([]byte, len(m.ram))
	copy(memCopy, m.ram)
	return memCopy
}
