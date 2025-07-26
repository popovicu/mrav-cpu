package isa

func ParseInstructionCode(instructionRegister Register) InstructionCode {
	return InstructionCode(instructionRegister >> 12)
}

func ParseRd(instructionRegister Register) RegisterId {
	return RegisterId((instructionRegister & 0x0F00) >> 8)
}

func ParseRs1(instructionRegister Register) RegisterId {
	return RegisterId((instructionRegister & 0x00F0) >> 4)
}

func ParseRs2(instructionRegister Register) RegisterId {
	return RegisterId(instructionRegister & 0x000F)
}

func Imm8(instructionRegister Register) uint8 {
	return uint8(instructionRegister & 0x00FF)
}

func Imm4(instructionRegister Register) uint8 {
	return uint8((instructionRegister & 0x00F0) >> 4)
}
