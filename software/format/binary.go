package format

import (
	"bytes"

	"mrav/software/machinecode"
	"mrav/software/model"
)

func Binary(m *model.MravModule) ([]byte, error) {
	var binaryOutput bytes.Buffer

	for _, instr := range m.Instructions {
		if err := machinecode.GenerateMachineCodeForInstruction(instr, &binaryOutput); err != nil {
			return nil, err
		}
	}

	return binaryOutput.Bytes(), nil
}
