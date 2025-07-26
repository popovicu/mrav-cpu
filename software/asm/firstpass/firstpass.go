package firstpass

import (
	"fmt"
	"mrav/isa"
	"mrav/software/asm/parsing"
	"mrav/software/model"
	"slices"
)

func ModuleFirstPass(m *parsing.Module) (*model.MravModule, error) {
	assignments := make([]model.MravDefinition, 0)
	labels := make([]model.MravLabel, 0)
	instructions := make([]model.MravInstruction, 0)

	var runningPc isa.Register = 0

	for _, line := range m.Lines {
		lineContent := line.Content

		var matchErr error = nil

		// Match API is a bit weird...
		lineContent.Match(
			func(_ parsing.BlankLine) parsing.AssemblyLine {
				return parsing.AssemblyBlankLine() // This should do nothing
			},
			func(assignment parsing.AssignmentLine) parsing.AssemblyLine {
				val, err := parsing.NumberValue(string(assignment.Value))

				if err != nil {
					matchErr = err
					return parsing.AssemblyBlankLine()
				}

				assignments = append(assignments, model.MravDefinition{
					Symbol: model.MravSymbol(assignment.Sym),
					Value:  model.MravValue(val),
				})

				return parsing.AssemblyBlankLine() // This should do nothing
			},
			func(ll parsing.LabelLine) parsing.AssemblyLine {
				labels = append(labels, model.MravLabel{
					Symbol:  model.MravSymbol(ll.Label),
					Address: model.MravValue(runningPc),
				})

				return parsing.AssemblyBlankLine() // This should do nothing
			},
			func(il parsing.InstructionLine) parsing.AssemblyLine {
				instr, err := processInstruction(il.Instruction)

				if err != nil {
					matchErr = fmt.Errorf("error on line %d, cannot parse instruction: %w", line.Number, err)
					return parsing.AssemblyBlankLine()
				}

				instructions = append(instructions, instr)
				runningPc += 2
				return parsing.AssemblyBlankLine() // This should do nothing
			},
			func(lil parsing.LabeledInstructionLine) parsing.AssemblyLine {
				labels = append(labels, model.MravLabel{
					Symbol:  model.MravSymbol(lil.Label),
					Address: model.MravValue(runningPc),
				})

				instr, err := processInstruction(lil.Instruction)

				if err != nil {
					matchErr = fmt.Errorf("error on line %d, cannot parse instruction: %w", line.Number, err)
					return parsing.AssemblyBlankLine()
				}

				instructions = append(instructions, instr)
				runningPc += 2
				return parsing.AssemblyBlankLine() // This should do nothing
			},
		)

		if matchErr != nil {
			return &model.MravModule{}, matchErr
		}
	}

	return &model.MravModule{
		Labels:          labels,
		AssignedSymbols: assignments,
		Instructions:    instructions,
	}, nil
}

func processInstruction(inst parsing.Instruction) (model.MravInstruction, error) {
	switch inst.CpuInstruction {
	case isa.ADD, isa.SUB, isa.XOR, isa.AND, isa.OR:
		instr, err := processRdRs1Rs2Instruction(inst)

		if err != nil {
			return model.MravInstruction{}, err
		}

		return instr, nil
	case isa.ADDI, isa.LDHI, isa.BZ, isa.BNZ, isa.JAL:
		instr, err := processRdImm8(inst)

		if err != nil {
			return model.MravInstruction{}, err
		}

		return instr, nil
	case isa.SHL, isa.SHR, isa.SHRA:
		instr, err := processRdImm4(inst)

		if err != nil {
			return model.MravInstruction{}, err
		}

		return instr, nil
	case isa.LW, isa.SW, isa.JALR:
		stringInstruction, err := isa.InstructionToString(inst.CpuInstruction)

		if err != nil {
			return model.MravInstruction{}, err
		}

		if len(inst.Args) != 1 {
			return model.MravInstruction{}, fmt.Errorf("%s instruction should have arguments rd, rs1,", stringInstruction)
		}

		rs1, err := parsing.ParseRegister(inst.Args[0].UnprocessedValue)

		if err != nil {
			return model.MravInstruction{}, fmt.Errorf("cannot parse rs1 of %s instruction: %w", stringInstruction, err)
		}

		if inst.CpuInstruction == isa.JALR {
			return model.MravInstruction{
				Jalr: &model.MravJalr{
					Rd:  inst.Rd,
					Rs1: rs1,
				},
			}, nil
		}

		if inst.CpuInstruction == isa.LW {
			return model.MravInstruction{
				Lw: &model.MravLw{
					Rd:  inst.Rd,
					Rs1: rs1,
				},
			}, nil
		}

		return model.MravInstruction{
			Sw: &model.MravSw{
				Rd:  inst.Rd,
				Rs1: rs1,
			},
		}, nil
	default:
		return model.MravInstruction{}, fmt.Errorf("unknown instruction (this should absolutely never happen!)")
	}
}

func processRdRs1Rs2Instruction(inst parsing.Instruction) (model.MravInstruction, error) {
	instructions := []isa.InstructionCode{isa.ADD, isa.SUB, isa.XOR, isa.AND, isa.OR}

	stringInstruction, err := isa.InstructionToString(inst.CpuInstruction)

	if err != nil {
		return model.MravInstruction{}, err
	}

	if !slices.Contains(instructions, inst.CpuInstruction) {
		return model.MravInstruction{}, fmt.Errorf("%s is not a 'rd rs1 rs2' instruction", stringInstruction)
	}

	if len(inst.Args) != 2 {
		return model.MravInstruction{}, fmt.Errorf("%s instruction should have arguments rd, rs1, rs2", stringInstruction)
	}

	rs1, err := parsing.ParseRegister(inst.Args[0].UnprocessedValue)

	if err != nil {
		return model.MravInstruction{}, fmt.Errorf("cannot parse rs1 of %s instruction: %w", stringInstruction, err)
	}

	rs2, err := parsing.ParseRegister(inst.Args[1].UnprocessedValue)

	if err != nil {
		return model.MravInstruction{}, fmt.Errorf("cannot parse rs2 of %s instruction: %w", stringInstruction, err)
	}

	var mravInstruction model.MravInstruction

	switch inst.CpuInstruction {
	case isa.ADD:
		mravInstruction = model.MravInstruction{
			Add: &model.MravAdd{
				Rd:  inst.Rd,
				Rs1: rs1,
				Rs2: rs2,
			},
		}
	case isa.SUB:
		mravInstruction = model.MravInstruction{
			Sub: &model.MravSub{
				Rd:  inst.Rd,
				Rs1: rs1,
				Rs2: rs2,
			},
		}
	case isa.XOR:
		mravInstruction = model.MravInstruction{
			Xor: &model.MravXor{
				Rd:  inst.Rd,
				Rs1: rs1,
				Rs2: rs2,
			},
		}
	case isa.AND:
		mravInstruction = model.MravInstruction{
			And: &model.MravAnd{
				Rd:  inst.Rd,
				Rs1: rs1,
				Rs2: rs2,
			},
		}
	case isa.OR:
		mravInstruction = model.MravInstruction{
			Or: &model.MravOr{
				Rd:  inst.Rd,
				Rs1: rs1,
				Rs2: rs2,
			},
		}
	default:
		return mravInstruction, fmt.Errorf("cannot process a 'rd rs1 rs2' instruction, though this should never happen!")
	}

	return mravInstruction, nil
}

func processRdImm8(inst parsing.Instruction) (model.MravInstruction, error) {
	instructions := []isa.InstructionCode{isa.ADDI, isa.LDHI, isa.BZ, isa.BNZ, isa.JAL}

	stringInstruction, err := isa.InstructionToString(inst.CpuInstruction)

	if err != nil {
		return model.MravInstruction{}, err
	}

	if !slices.Contains(instructions, inst.CpuInstruction) {
		return model.MravInstruction{}, fmt.Errorf("%s is not a 'rd imm8' instruction", stringInstruction)
	}

	if len(inst.Args) != 1 {
		return model.MravInstruction{}, fmt.Errorf("%s instruction should have arguments rd, imm8", stringInstruction)
	}

	var immOrSymb model.ImmOrSymb

	if inst.Args[0].ArgType == parsing.INSTRUCTION_ARG_TYPE_IDENTIFIER {
		if _, err := parsing.ParseRegister(inst.Args[0].UnprocessedValue); err == nil {
			// Weird condition, but this is what we need: if this actually parses as a register reference, we want to raise an error.
			return model.MravInstruction{}, fmt.Errorf("%s instruction cannot use a reigster as its second argument", stringInstruction)
		}

		immOrSymb = model.ImmOrSymbFromSymb(model.MravSymbol(inst.Args[0].UnprocessedValue))
	} else {
		imm8, err := parsing.NumberValue(inst.Args[0].UnprocessedValue)

		if err != nil {
			return model.MravInstruction{}, fmt.Errorf("cannot parse imm8 of %s instruction: %w", stringInstruction, err)
		}

		if imm8 > 0xFF {
			return model.MravInstruction{}, fmt.Errorf("cannot parse imm8 of %s instruction, value too large: %w", stringInstruction, err)
		}

		immOrSymb = model.ImmOrSymbFromImm(uint8(imm8))
	}

	var mravInstruction model.MravInstruction

	switch inst.CpuInstruction {
	case isa.ADDI:
		mravInstruction = model.MravInstruction{
			Addi: &model.MravAddi{
				Rd:    inst.Rd,
				Value: immOrSymb,
			},
		}
	case isa.LDHI:
		mravInstruction = model.MravInstruction{
			Ldhi: &model.MravLdhi{
				Rd:    inst.Rd,
				Value: immOrSymb,
			},
		}
	case isa.BZ:
		mravInstruction = model.MravInstruction{
			Bz: &model.MravBz{
				Rd:   inst.Rd,
				Addr: immOrSymb,
			},
		}
	case isa.BNZ:
		mravInstruction = model.MravInstruction{
			Bnz: &model.MravBnz{
				Rd:   inst.Rd,
				Addr: immOrSymb,
			},
		}
	case isa.JAL:
		mravInstruction = model.MravInstruction{
			Jal: &model.MravJal{
				Rd:   inst.Rd,
				Addr: immOrSymb,
			},
		}
	default:
		return mravInstruction, fmt.Errorf("cannot process a 'rd imm8' instruction, though this should never happen!")
	}

	return mravInstruction, nil
}

func processRdImm4(inst parsing.Instruction) (model.MravInstruction, error) {
	instructions := []isa.InstructionCode{isa.SHL, isa.SHR, isa.SHRA}

	stringInstruction, err := isa.InstructionToString(inst.CpuInstruction)

	if err != nil {
		return model.MravInstruction{}, err
	}

	if !slices.Contains(instructions, inst.CpuInstruction) {
		return model.MravInstruction{}, fmt.Errorf("%s is not a 'rd imm4' instruction", stringInstruction)
	}

	if len(inst.Args) != 1 {
		return model.MravInstruction{}, fmt.Errorf("%s instruction should have arguments rd, imm4", stringInstruction)
	}

	imm4, err := parsing.NumberValue(inst.Args[0].UnprocessedValue)

	if err != nil {
		return model.MravInstruction{}, fmt.Errorf("cannot parse imm4 of %s instruction: %w", stringInstruction, err)
	}

	if imm4 > 0xF {
		return model.MravInstruction{}, fmt.Errorf("cannot parse imm4 of %s instruction, value too large: %w", stringInstruction, err)
	}

	var mravInstruction model.MravInstruction

	switch inst.CpuInstruction {
	case isa.SHL:
		mravInstruction = model.MravInstruction{
			Shl: &model.MravShl{
				Rd:   inst.Rd,
				Imm4: uint8(imm4),
			},
		}
	case isa.SHR:
		mravInstruction = model.MravInstruction{
			Shr: &model.MravShr{
				Rd:   inst.Rd,
				Imm4: uint8(imm4),
			},
		}
	case isa.SHRA:
		mravInstruction = model.MravInstruction{
			Shra: &model.MravShra{
				Rd:   inst.Rd,
				Imm4: uint8(imm4),
			},
		}
	default:
		return mravInstruction, fmt.Errorf("cannot process a 'rd imm4' instruction, though this should never happen!")
	}

	return mravInstruction, nil
}
