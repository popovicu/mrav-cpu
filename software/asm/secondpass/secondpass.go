package secondpass

import (
	"fmt"
	"mrav/software/model"
)

type MravObject struct {
	Module            *model.MravModule
	UnresolvedSymbols []model.MravSymbol
}

func BuildObject(m *model.MravModule) (MravObject, error) {
	moduleSymbols := make(map[model.MravSymbol]struct{})

	for _, label := range m.Labels {
		if _, exists := moduleSymbols[label.Symbol]; exists {
			return MravObject{}, fmt.Errorf("duplicate symbol '%s'", label.Symbol)
		}

		moduleSymbols[label.Symbol] = struct{}{}
	}

	for _, def := range m.AssignedSymbols {
		if _, exists := moduleSymbols[def.Symbol]; exists {
			return MravObject{}, fmt.Errorf("duplicate symbol '%s'", def.Symbol)
		}

		moduleSymbols[def.Symbol] = struct{}{}
	}

	referencedSymbols := make([]model.MravSymbol, 0)

	for _, instr := range m.Instructions {
		if (instr.Addi != nil) && (instr.Addi.Value.IsRight()) {
			referencedSymbols = append(referencedSymbols, instr.Addi.Value.MustRight())
			continue
		}

		if (instr.Ldhi != nil) && (instr.Ldhi.Value.IsRight()) {
			referencedSymbols = append(referencedSymbols, instr.Ldhi.Value.MustRight())
			continue
		}

		if (instr.Bz != nil) && (instr.Bz.Addr.IsRight()) {
			referencedSymbols = append(referencedSymbols, instr.Bz.Addr.MustRight())
			continue
		}

		if (instr.Bnz != nil) && (instr.Bnz.Addr.IsRight()) {
			referencedSymbols = append(referencedSymbols, instr.Bnz.Addr.MustRight())
			continue
		}

		if (instr.Jal != nil) && (instr.Jal.Addr.IsRight()) {
			referencedSymbols = append(referencedSymbols, instr.Jal.Addr.MustRight())
			continue
		}
	}

	unresolvedSymbols := make([]model.MravSymbol, 0)

	for _, refSymb := range referencedSymbols {
		if _, found := moduleSymbols[refSymb]; !found {
			unresolvedSymbols = append(unresolvedSymbols, refSymb)
		}
	}

	return MravObject{
		Module:            m,
		UnresolvedSymbols: unresolvedSymbols,
	}, nil
}
