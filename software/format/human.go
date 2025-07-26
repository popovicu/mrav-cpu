package format

import (
	"bytes"
	"fmt"

	"mrav/software/machinecode"
	"mrav/software/model"
)

func HumanReadable(m *model.MravModule) ([]string, error) {
	output := make([]string, 0, len(m.Instructions))

	for _, instr := range m.Instructions {
		var buf bytes.Buffer

		if err := machinecode.GenerateMachineCodeForInstruction(instr, &buf); err != nil {
			return nil, err
		}

		output = append(output, fmt.Sprintf("%X", buf.Bytes()))
	}

	return output, nil
}
