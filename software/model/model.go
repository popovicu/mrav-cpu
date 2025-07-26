package model

type MravSymbol string

type MravValue uint16 // Physically represented as 16-bit, but could be small enough to fit 8-bit.

type MravModule struct {
	Labels          []MravLabel
	AssignedSymbols []MravDefinition
	Instructions    []MravInstruction
}

type MravLabel struct {
	Symbol  MravSymbol
	Address MravValue
}

type MravDefinition struct {
	Symbol MravSymbol
	Value  MravValue
}

type MravInstruction struct {
	// Exactly one should be non-null. Hard to enforce in Go, honor system.
	Add  *MravAdd
	Sub  *MravSub
	Lw   *MravLw
	Sw   *MravSw
	Xor  *MravXor
	And  *MravAnd
	Or   *MravOr
	Addi *MravAddi
	Ldhi *MravLdhi
	Bz   *MravBz
	Bnz  *MravBnz
	Jal  *MravJal
	Jalr *MravJalr
	Shl  *MravShl
	Shr  *MravShr
	Shra *MravShra
}
