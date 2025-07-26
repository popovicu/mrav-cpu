package linker

import (
	"fmt"

	"mrav/isa"
	"mrav/software/asm/secondpass"
	"mrav/software/model"
)

func Link(objects []*secondpass.MravObject) (*model.MravModule, error) {
	type objSymbol struct {
		module     int
		symbolType int // 1 for label, 2 otherwise
		value      uint8
	}

	symbolsToObjects := make(map[model.MravSymbol]objSymbol) // TODO: don't just map to int, but also the type

	// Check for duplicate symbols
	for i, obj := range objects {
		for _, symb := range obj.Module.AssignedSymbols {
			if objIdx, exists := symbolsToObjects[symb.Symbol]; exists {
				return nil, fmt.Errorf("symbol '%s' already defined in object %d", symb.Symbol, objIdx)
			}

			symbolsToObjects[symb.Symbol] = objSymbol{
				module:     i,
				symbolType: 2,
				value:      uint8(symb.Value),
			}
		}

		for _, label := range obj.Module.Labels {
			if objIdx, exists := symbolsToObjects[label.Symbol]; exists {
				return nil, fmt.Errorf("symbol '%s' already defined in object %d", label.Symbol, objIdx)
			}

			symbolsToObjects[label.Symbol] = objSymbol{
				module:     i,
				symbolType: 1,
				value:      uint8(label.Address),
			}
		}
	}

	// Check if it will be possible to resolve all unresolved symbols.
	for i, obj := range objects {
		for _, symb := range obj.UnresolvedSymbols {
			if _, found := symbolsToObjects[symb]; !found {
				return nil, fmt.Errorf("symbol '%s' from module %d is unresolved", symb, i)
			}
		}
	}

	objectsToOffset := make([]isa.Register, len(objects))
	offset := isa.Register(0)
	totalInstructions := 0

	for i, obj := range objects {
		objectsToOffset[i] = offset
		objectInstructions := len(obj.Module.Instructions)
		offset += isa.Register(2 * objectInstructions) // TODO: establish a limit on module instruction size
		totalInstructions += objectInstructions
	}

	linkedInstructions := make([]model.MravInstruction, 0, totalInstructions)

	for _, obj := range objects {
		for _, instr := range obj.Module.Instructions {
			if (instr.Addi != nil) && (instr.Addi.Value.IsRight()) {
				symb := instr.Addi.Value.MustRight()
				objMeta := symbolsToObjects[symb]
				finalValue := objMeta.value

				if objMeta.symbolType == 1 {
					symbObjOffset := objectsToOffset[objMeta.module]
					finalValue += uint8(symbObjOffset)
				}

				linkedInstructions = append(linkedInstructions, model.MravInstruction{
					Addi: &model.MravAddi{
						Rd:    instr.Addi.Rd,
						Value: model.ImmOrSymbFromImm(uint8(finalValue)),
					},
				})
				continue
			}

			if (instr.Ldhi != nil) && (instr.Ldhi.Value.IsRight()) {
				symb := instr.Ldhi.Value.MustRight()
				objMeta := symbolsToObjects[symb]
				finalValue := objMeta.value

				if objMeta.symbolType == 1 {
					symbObjOffset := objectsToOffset[objMeta.module]
					finalValue += uint8(symbObjOffset)
				}

				linkedInstructions = append(linkedInstructions, model.MravInstruction{
					Ldhi: &model.MravLdhi{
						Rd:    instr.Ldhi.Rd,
						Value: model.ImmOrSymbFromImm(uint8(finalValue)),
					},
				})
				continue
			}

			if (instr.Bz != nil) && (instr.Bz.Addr.IsRight()) {
				symb := instr.Bz.Addr.MustRight()
				objMeta := symbolsToObjects[symb]
				finalValue := objMeta.value

				if objMeta.symbolType == 1 {
					symbObjOffset := objectsToOffset[objMeta.module]
					finalValue += uint8(symbObjOffset)
				}

				linkedInstructions = append(linkedInstructions, model.MravInstruction{
					Bz: &model.MravBz{
						Rd:   instr.Bz.Rd,
						Addr: model.ImmOrSymbFromImm(uint8(finalValue)),
					},
				})
				continue
			}

			if (instr.Bnz != nil) && (instr.Bnz.Addr.IsRight()) {
				symb := instr.Bnz.Addr.MustRight()
				objMeta := symbolsToObjects[symb]
				finalValue := objMeta.value

				if objMeta.symbolType == 1 {
					symbObjOffset := objectsToOffset[objMeta.module]
					finalValue += uint8(symbObjOffset)
				}

				linkedInstructions = append(linkedInstructions, model.MravInstruction{
					Bnz: &model.MravBnz{
						Rd:   instr.Bnz.Rd,
						Addr: model.ImmOrSymbFromImm(uint8(finalValue)),
					},
				})
				continue
			}

			if (instr.Jal != nil) && (instr.Jal.Addr.IsRight()) {
				symb := instr.Jal.Addr.MustRight()
				objMeta := symbolsToObjects[symb]
				finalValue := objMeta.value

				if objMeta.symbolType == 1 {
					symbObjOffset := objectsToOffset[objMeta.module]
					finalValue += uint8(symbObjOffset)
				}

				linkedInstructions = append(linkedInstructions, model.MravInstruction{
					Jal: &model.MravJal{
						Rd:   instr.Jal.Rd,
						Addr: model.ImmOrSymbFromImm(uint8(finalValue)),
					},
				})
				continue
			}

			linkedInstructions = append(linkedInstructions, instr)
		}
	}

	return &model.MravModule{
		Instructions: linkedInstructions,
	}, nil
}
