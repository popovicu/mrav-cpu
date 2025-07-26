package core

import (
	"fmt"
	"math"
	"os"

	protobuf "google.golang.org/protobuf/proto"

	"mrav/core/proto"
	"mrav/isa"
)

func (c *Core) Serialize() *proto.CoreState {
	regs := make([]uint32, len(c.Registers))

	for i := range c.Registers {
		regs[i] = uint32(c.Registers[i])
	}

	return &proto.CoreState{
		Pc: uint32(c.Pc),
		R:  regs,
	}
}

func Deserialize(protoCore *proto.CoreState) (*Core, error) {
	if protoCore.Pc > math.MaxUint16 {
		return nil, fmt.Errorf("cannot deserialize the core, PC too large: %X", protoCore.Pc)
	}

	regs := make([]isa.Register, 0, isa.RegsNumber)

	for i := range protoCore.R {
		if protoCore.R[i] > math.MaxUint16 {
			return nil, fmt.Errorf("cannot deserialize the core, register %d too large: %X", i, protoCore.R[i])
		}

		regs[i] = isa.Register(protoCore.R[i])
	}

	return &Core{
		Pc:        isa.Register(protoCore.Pc),
		Registers: isa.GeneralRegisters(regs),
	}, nil
}

func (c *Core) SerializeToBytes() ([]byte, error) {
	coreBytes, err := protobuf.Marshal(c.Serialize())

	if err != nil {
		return nil, err
	}

	return coreBytes, nil
}

func DeserializeFromBytes(coreBytes []byte) (*Core, error) {
	coreProto := &proto.CoreState{}

	if err := protobuf.Unmarshal(coreBytes, coreProto); err != nil {
		return nil, fmt.Errorf("cannot unmarshal from bytes, proto error: %w", err)
	}

	modelCore, err := Deserialize(coreProto)

	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal from bytes, model error: %w", err)
	}

	return modelCore, nil
}

func (c *Core) SnapshotToFile(filePath string) error {
	storageBytes, err := c.SerializeToBytes()

	if err != nil {
		return fmt.Errorf("unable to snapshot core to the file, cannot serialize: %w", err)
	}

	if err := os.WriteFile(filePath, storageBytes, 0644); err != nil {
		return fmt.Errorf("cannot snapshot core to the file, file writing error: %w", err)
	}

	return nil
}

func ReadSnapshotFromFile(filePath string) (*Core, error) {
	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("cannot read core from file, file read error: %w", err)
	}

	modelCore, err := DeserializeFromBytes(fileBytes)

	if err != nil {
		return nil, fmt.Errorf("cannot read core from file, deserialization error: %w", err)
	}

	return modelCore, nil
}
