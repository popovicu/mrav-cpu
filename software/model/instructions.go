package model

import (
	"github.com/samber/mo"

	"mrav/isa"
)

type ImmOrSymb = mo.Either[uint8, MravSymbol]

func ImmOrSymbFromImm(val uint8) ImmOrSymb {
	return ImmOrSymb(mo.Left[uint8, MravSymbol](val))
}

func ImmOrSymbFromSymb(val MravSymbol) ImmOrSymb {
	return ImmOrSymb(mo.Right[uint8, MravSymbol](val))
}

type MravAdd struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
	Rs2 isa.RegisterId
}

type MravSub struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
	Rs2 isa.RegisterId
}

type MravLw struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
}

type MravSw struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
}

type MravXor struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
	Rs2 isa.RegisterId
}

type MravAnd struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
	Rs2 isa.RegisterId
}

type MravOr struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
	Rs2 isa.RegisterId
}

type MravAddi struct {
	Rd    isa.RegisterId
	Value ImmOrSymb
}

type MravLdhi struct {
	Rd    isa.RegisterId
	Value ImmOrSymb
}

type MravBz struct {
	Rd   isa.RegisterId
	Addr ImmOrSymb
}

type MravBnz struct {
	Rd   isa.RegisterId
	Addr ImmOrSymb
}

type MravJal struct {
	Rd   isa.RegisterId
	Addr ImmOrSymb
}

type MravJalr struct {
	Rd  isa.RegisterId
	Rs1 isa.RegisterId
}

type MravShl struct {
	Rd   isa.RegisterId
	Imm4 uint8
}

type MravShr struct {
	Rd   isa.RegisterId
	Imm4 uint8
}

type MravShra struct {
	Rd   isa.RegisterId
	Imm4 uint8
}
