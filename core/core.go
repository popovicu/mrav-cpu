package core

import (
	"bytes"
	"fmt"
	"log/slog"
	"slices"

	"mrav/isa"
)

type State int

const (
	STATE_READY State = iota
	STATE_FETCH_AND_RUN
	STATE_LW_WAITING
	STATE_SW_WAITING
)

type Core struct {
	Registers isa.GeneralRegisters
	Pc        isa.Register
	Sp        isa.Register

	logger  *slog.Logger
	verbose bool

	state       State
	instruction isa.Register
}

type CoreOpts struct {
	Logger  *slog.Logger
	Verbose bool
}

func NewCore(opts *CoreOpts) *Core {
	core := &Core{
		Pc:    isa.Register(0x0000),
		Sp:    isa.Register(0x0000),
		state: STATE_READY,

		logger:  opts.Logger,
		verbose: opts.Verbose,
	}

	for idx := range core.Registers {
		core.Registers[idx] = isa.Register(0x0000)
	}

	return core
}

type ExecutionSignal int

const (
	SIGNAL_FETCH_INSTRUCTION ExecutionSignal = iota
	SIGNAL_DONE
	SIGNAL_LOADING_DATA
	SIGNAL_WRITING_DATA
)

// MultiturnRunInstruction is supposed to be called repeatedly until the instruction is executed.
//
// The idea is to very loosely couple the bus and the core. The system will invoke the method on the core, which can generate a bus access, but not be done.
// After that, the system responds by servicing the bus access and invoking the next stage of the multiturn run.
// This obviously happens when the core requests the next instruction from the bus, gets the data and then executes it, potentially generating another bus access.
//
// Execution signals can be multiple. For the simplest invocation here, keep invoking the method until a stop signal is generated.
func (c *Core) MultiturnRunInstruction(busValue isa.BusValue) (*isa.BusAccess, []ExecutionSignal, error) {
	if c.state == STATE_READY {
		c.state = STATE_FETCH_AND_RUN

		if c.verbose {
			c.logger.Info("[Core] Reading instruction", "pc", fmt.Sprintf("%04X", c.Pc))
		}

		return &isa.BusAccess{
			Read: &isa.BusAccessRead{
				Address: c.Pc,
			},
		}, []ExecutionSignal{SIGNAL_FETCH_INSTRUCTION}, nil
	}

	if c.state == STATE_LW_WAITING {
		c.state = STATE_READY
		rd := isa.ParseRd(c.instruction)
		c.Registers[rd] = isa.Register(busValue)
		c.Pc += isa.INSTRUCTION_SIZE
		return nil, []ExecutionSignal{SIGNAL_DONE}, nil
	}

	if c.state == STATE_SW_WAITING {
		c.state = STATE_READY
		c.Pc += isa.INSTRUCTION_SIZE
		return nil, []ExecutionSignal{SIGNAL_DONE}, nil
	}

	if c.state == STATE_FETCH_AND_RUN {
		c.state = STATE_READY
		c.instruction = isa.Register(busValue)

		instrCode := isa.ParseInstructionCode(c.instruction)
		instrString, err := isa.InstructionToString(instrCode)

		if err != nil {
			return nil, nil, fmt.Errorf("cannot find the name of the instruction: %w", err)
		}

		if c.verbose {
			c.logger.Info("[Core] Running instruction", "instruction", instrString, "hex", fmt.Sprintf("%04X", c.instruction))
		}

		switch instrCode {
		case isa.InstructionCode(isa.ADD):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			rs2 := isa.ParseRs2(c.instruction)
			c.Registers[rd] = isa.Register(c.Registers[rs1] + c.Registers[rs2])
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.SUB):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			rs2 := isa.ParseRs2(c.instruction)
			c.Registers[rd] = isa.Register(c.Registers[rs1] - c.Registers[rs2])
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.LW):
			rs := isa.ParseRs1(c.instruction)
			c.state = STATE_LW_WAITING
			return &isa.BusAccess{
				Read: &isa.BusAccessRead{
					Address: c.Registers[rs],
				},
			}, []ExecutionSignal{SIGNAL_LOADING_DATA}, nil
		case isa.InstructionCode(isa.SW):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			c.state = STATE_SW_WAITING
			return &isa.BusAccess{
				Write: &isa.BusAccessWrite{
					Address: c.Registers[rd],
					Value:   c.Registers[rs1],
				},
			}, []ExecutionSignal{SIGNAL_WRITING_DATA}, nil
		case isa.InstructionCode(isa.XOR):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			rs2 := isa.ParseRs2(c.instruction)
			c.Registers[rd] = isa.Register(c.Registers[rs1] ^ c.Registers[rs2])
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.AND):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			rs2 := isa.ParseRs2(c.instruction)
			c.Registers[rd] = isa.Register(c.Registers[rs1] & c.Registers[rs2])
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.OR):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			rs2 := isa.ParseRs2(c.instruction)
			c.Registers[rd] = isa.Register(c.Registers[rs1] | c.Registers[rs2])
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.ADDI):
			rd := isa.ParseRd(c.instruction)
			imm8 := isa.Imm8(c.instruction)
			c.Registers[rd] = isa.Register(uint16(c.Registers[rd]) + uint16(imm8))
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.LDHI):
			rd := isa.ParseRd(c.instruction)
			imm8 := isa.Imm8(c.instruction)
			c.Registers[rd] = isa.Register(uint16(uint16(imm8)<<8) | uint16(c.Registers[rd]&0x00FF))
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.BZ):
			rd := isa.ParseRd(c.instruction)
			imm8 := isa.Imm8(c.instruction)
			if c.Registers[rd] == 0 {
				c.Pc = isa.Register(imm8)
			} else {
				c.Pc += isa.INSTRUCTION_SIZE
			}
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.BNZ):
			rd := isa.ParseRd(c.instruction)
			imm8 := isa.Imm8(c.instruction)
			if c.Registers[rd] != 0 {
				c.Pc = isa.Register(imm8)
			} else {
				c.Pc += isa.INSTRUCTION_SIZE
			}
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.JAL):
			rd := isa.ParseRd(c.instruction)
			imm8 := isa.Imm8(c.instruction)
			c.Registers[rd] = c.Pc + isa.INSTRUCTION_SIZE
			c.Pc = isa.Register(imm8)
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.JALR):
			rd := isa.ParseRd(c.instruction)
			rs1 := isa.ParseRs1(c.instruction)
			c.Registers[rd] = c.Pc + isa.INSTRUCTION_SIZE
			c.Pc = c.Registers[rs1]
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.SHL):
			rd := isa.ParseRd(c.instruction)
			imm4 := isa.Imm4(c.instruction)
			c.Registers[rd] <<= imm4
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.SHR):
			rd := isa.ParseRd(c.instruction)
			imm4 := isa.Imm4(c.instruction)
			c.Registers[rd] = isa.Register(uint16(c.Registers[rd]) >> imm4)
			c.Pc += isa.INSTRUCTION_SIZE
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		case isa.InstructionCode(isa.SHRA):
			rd := isa.ParseRd(c.instruction)
			imm4 := isa.Imm4(c.instruction)
			c.Pc += isa.INSTRUCTION_SIZE
			c.Registers[rd] = isa.Register(int16(c.Registers[rd]) >> imm4)
			return nil, []ExecutionSignal{SIGNAL_DONE}, nil
		default:
			return nil, nil, fmt.Errorf("unknown instruction code %#x", instrCode)
		}
	}

	return nil, nil, fmt.Errorf("unknown state") // TODO: better error message
}

func (c *Core) DebugDump(regsToDump []isa.RegisterId) (string, error) {
	var buf bytes.Buffer

	if _, err := fmt.Fprintf(&buf, "PC = %04X, [ ", c.Pc); err != nil {
		return "", fmt.Errorf("unable to produce a debug snapshot: %w", err)
	}

	for i := 0; i < int(isa.RegsNumber); i++ {
		if !slices.Contains(regsToDump, isa.RegisterId(i)) {
			continue
		}

		if _, err := fmt.Fprintf(&buf, "r%d = %04X ", i, c.Registers[i]); err != nil {
			return "", fmt.Errorf("cannot produce registers debug snapshot: %w", err)
		}
	}

	fmt.Fprintf(&buf, "]")

	return buf.String(), nil
}
