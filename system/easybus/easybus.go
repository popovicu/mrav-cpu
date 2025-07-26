package easybus

import (
	"fmt"
	"log/slog"
	"strings"

	"mrav/core"
	"mrav/isa"
	"mrav/system"
	"mrav/system/easybus/device"
)

type EasyBusSystem struct {
	core    *core.Core
	devices []device.Device

	logger  *slog.Logger
	verbose bool
}

func NewEasyBusSystem(opts *system.SystemOpts, devices []device.Device) (*EasyBusSystem, error) {
	coreOpts := &core.CoreOpts{
		Logger:  opts.Logger,
		Verbose: opts.Verbose,
	}

	system := &EasyBusSystem{
		core:    core.NewCore(coreOpts),
		devices: devices,
		logger:  opts.Logger,
		verbose: opts.Verbose,
	}

	return system, nil
}

func (sys *EasyBusSystem) hitDevice(address isa.BusValue) (device.Device, error) {
	hits := make([]device.Device, 0)

	for _, dev := range sys.devices {
		if dev.Hit(address) {
			hits = append(hits, dev)
		}
	}

	if len(hits) == 0 {
		return nil, fmt.Errorf("no device found for address: %04X", address)
	}

	if len(hits) > 1 {
		hitNames := make([]string, 0, len(hits))

		for _, hitDev := range hits {
			hitNames = append(hitNames, hitDev.Name())
		}

		return nil, fmt.Errorf("multiple devices hit on the bus: %s", strings.Join(hitNames, ", "))
	}

	return hits[0], nil
}

func (sys *EasyBusSystem) readBus(address isa.BusValue) (isa.BusValue, error) {
	if sys.verbose {
		sys.logger.Info("[EasyBus system] Read", "address", fmt.Sprintf("%04X", address))
	}

	busDevice, err := sys.hitDevice(address)

	if err != nil {
		return 0, err
	}

	return busDevice.ReadBus(address)
}

func (sys *EasyBusSystem) writeBus(address isa.BusValue, value isa.BusValue) error {
	if sys.verbose {
		sys.logger.Info("[EasyBus system] Write", "address", fmt.Sprintf("%04X", address), "value", fmt.Sprintf("%04X", value))
	}

	busDevice, err := sys.hitDevice(address)

	if err != nil {
		return err
	}

	return busDevice.WriteBus(address, value)
}

func (sys *EasyBusSystem) CoreDebug(regsToDump []isa.RegisterId) (string, error) {
	return sys.core.DebugDump(regsToDump)
}

func (sys *EasyBusSystem) ProtoCoreDebugFile(filepath string) error {
	return sys.core.SnapshotToFile(filepath)
}

func (sys *EasyBusSystem) RunInstruction() error {
	done := false
	nextBusValue := isa.BusValue(0x0000)

	for !done {
		busAccess, signals, err := sys.core.MultiturnRunInstruction(nextBusValue)

		if err != nil {
			return fmt.Errorf("unable to run instruction in the system: %w", err)
		}

		for _, dev := range sys.devices {
			dev.TickCycle()
		}

		if busAccess != nil {
			if busAccess.Read != nil {
				addr := busAccess.Read.Address
				val, err := sys.readBus(isa.BusValue(addr))

				if err != nil {
					return fmt.Errorf("cannot read from RAM: %w", err)
				}

				nextBusValue = val
			} else if busAccess.Write != nil {
				addr := busAccess.Write.Address
				val := busAccess.Write.Value

				if err := sys.writeBus(isa.BusValue(addr), isa.BusValue(val)); err != nil {
					return fmt.Errorf("cannot write to bus: %w", err)
				}

				nextBusValue = isa.BusValue(0x0000)
			} else {
				return fmt.Errorf("bus access is neither read nor write")
			}
		}

		if (len(signals) == 1) && (signals[0] == core.SIGNAL_DONE) {
			done = true
		}
	}

	return nil
}
