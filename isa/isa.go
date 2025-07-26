package isa

import (
	"fmt"
	"strings"
)

//     Add: 0, // 0x0, add rd rs1 rs2
//     Sub: 1, // 0x1, sub rd rs1 rs2
//     Lw: 2, // 0x2, lw rd rs1 xxxx
//     Sw: 3, // 0x3, sw rd rs1 xxxx
//     Xor: 4, // 0x4, xor rd rs1 rs2
//     And: 5, // 0x5, and rd rs1 rs2
//     Or: 6, // 0x6, or rd rs1 rs2
//     Addi: 7, // 0x7, addi rd imm8
//     Ldhi: 8, // 0x8, ldhi rd imm8
//     Bz: 9, // 0x9, bz rd imm8
//     Bnz: 10, // 0xA, bnz rd imm8
//     Jal: 11, // 0xB, jal rd imm8
//     Jalr: 12, // 0xC, jalr rd rs1 xxxx
//     Shl: 13, // 0xD, shl rd imm4 xxxx
//     Shr: 14, // 0xE, shr rd imm4 xxxx
//     Shra: 15, // 0xF, shra rd imm4 xxxx

const (
	RegsNumber uint8      = 16
	MinRegId   RegisterId = 0
	MaxRegId   RegisterId = RegisterId(RegsNumber - 1)
)

type RegisterId uint8
type Register uint16
type GeneralRegisters [RegsNumber]Register

type InstructionCode uint8

const (
	ADD  InstructionCode = 0x0
	SUB  InstructionCode = 0x1
	LW   InstructionCode = 0x2
	SW   InstructionCode = 0x3
	XOR  InstructionCode = 0x4
	AND  InstructionCode = 0x5
	OR   InstructionCode = 0x6
	ADDI InstructionCode = 0x7
	LDHI InstructionCode = 0x8
	BZ   InstructionCode = 0x9
	BNZ  InstructionCode = 0xA
	JAL  InstructionCode = 0xB
	JALR InstructionCode = 0xC
	SHL  InstructionCode = 0xD
	SHR  InstructionCode = 0xE
	SHRA InstructionCode = 0xF
)

func StringToInstruction(instruction string) (InstructionCode, error) {
	uppered := strings.ToUpper(instruction)

	switch uppered {
	case "ADD":
		return ADD, nil
	case "SUB":
		return SUB, nil
	case "LW":
		return LW, nil
	case "SW":
		return SW, nil
	case "XOR":
		return XOR, nil
	case "AND":
		return AND, nil
	case "OR":
		return OR, nil
	case "ADDI":
		return ADDI, nil
	case "LDHI":
		return LDHI, nil
	case "BZ":
		return BZ, nil
	case "BNZ":
		return BNZ, nil
	case "JAL":
		return JAL, nil
	case "JALR":
		return JALR, nil
	case "SHL":
		return SHL, nil
	case "SHR":
		return SHR, nil
	case "SHRA":
		return SHRA, nil
	default:
		return 0, fmt.Errorf("unknown instruction: '%s'", uppered)
	}
}

func InstructionToString(instruction InstructionCode) (string, error) {
	switch instruction {
	case ADD:
		return "ADD", nil
	case SUB:
		return "SUB", nil
	case LW:
		return "LW", nil
	case SW:
		return "SW", nil
	case XOR:
		return "XOR", nil
	case AND:
		return "AND", nil
	case OR:
		return "OR", nil
	case ADDI:
		return "ADDI", nil
	case LDHI:
		return "LDHI", nil
	case BZ:
		return "BZ", nil
	case BNZ:
		return "BNZ", nil
	case JAL:
		return "JAL", nil
	case JALR:
		return "JALR", nil
	case SHL:
		return "SHL", nil
	case SHR:
		return "SHR", nil
	case SHRA:
		return "SHRA", nil
	default:
		return "", fmt.Errorf("unknown instruction: '0x%02X'", instruction)
	}
}

type BusAccessRead struct {
	Address Register
}

type BusAccessWrite struct {
	Address Register
	Value   Register
}

type BusAccess struct {
	// Only one should be non nil
	Read  *BusAccessRead
	Write *BusAccessWrite
}

type BusValue Register // Mrav connects to the bus of the same width as the register

const INSTRUCTION_SIZE = 2
