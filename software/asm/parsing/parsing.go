package parsing

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/davecgh/go-spew/spew"
	"github.com/samber/mo"

	"mrav/isa"
)

type Module struct {
	Lines []Line
}

type Line struct {
	Number  int
	Content AssemblyLine
}

type AssemblyLine = mo.Either5[BlankLine, AssignmentLine, LabelLine, InstructionLine, LabeledInstructionLine]

func AssemblyBlankLine() AssemblyLine {
	return AssemblyLine(mo.NewEither5Arg1[BlankLine, AssignmentLine, LabelLine, InstructionLine, LabeledInstructionLine](BlankLine{}))
}

func AssemblyAssignmentLine(sym Symbol, value HardcodedValue) AssemblyLine {
	return AssemblyLine(mo.NewEither5Arg2[BlankLine, AssignmentLine, LabelLine, InstructionLine, LabeledInstructionLine](AssignmentLine{
		Sym:   sym,
		Value: value,
	}))
}

func AssemblyLabelLine(label Label) AssemblyLine {
	return AssemblyLine(mo.NewEither5Arg3[BlankLine, AssignmentLine, LabelLine, InstructionLine, LabeledInstructionLine](LabelLine{Label: label}))
}

func AssemblyInstructionLine(instruction Instruction) AssemblyLine {
	return AssemblyLine(mo.NewEither5Arg4[BlankLine, AssignmentLine, LabelLine, InstructionLine, LabeledInstructionLine](InstructionLine{Instruction: instruction}))
}

func AssemblyLabeledInstructionLine(label Label, instruction Instruction) AssemblyLine {
	return AssemblyLine(mo.NewEither5Arg5[BlankLine, AssignmentLine, LabelLine, InstructionLine, LabeledInstructionLine](LabeledInstructionLine{
		Label:       label,
		Instruction: instruction,
	}))
}

type BlankLine struct{}

type AssignmentLine struct {
	Sym   Symbol
	Value HardcodedValue
}

type LabelLine struct {
	Label Label
}

type InstructionLine struct {
	Instruction Instruction
}

type LabeledInstructionLine struct {
	Label       Label
	Instruction Instruction
}

type Instruction struct {
	CpuInstruction isa.InstructionCode
	Rd             isa.RegisterId   // First arg is always rd, a register
	Args           []InstructionArg // These are yet unprocessed in this first phase of parsing
}

type Symbol string

type HardcodedValue string

type Comment string

type Label string

type InstructionArgType int

const (
	INSTRUCTION_ARG_TYPE_UNKNOWN InstructionArgType = iota
	INSTRUCTION_ARG_TYPE_NUMBER
	INSTRUCTION_ARG_TYPE_IDENTIFIER
)

type InstructionArg struct {
	UnprocessedValue string
	ArgType          InstructionArgType
}

func NumberValue(stringVal string) (uint16, error) {
	if strings.HasPrefix(stringVal, "0x") {
		val, err := strconv.ParseUint(stringVal[2:], 16, 16)

		if err != nil {
			return 0, fmt.Errorf("cannot extract a hex constant: %w", err)
		}

		return uint16(val), err
	}

	if strings.HasPrefix(stringVal, "0b") {
		val, err := strconv.ParseUint(stringVal[2:], 2, 16)

		if err != nil {
			return 0, fmt.Errorf("cannot extract a binary constant: %w", err)
		}

		return uint16(val), err
	}

	val, err := strconv.ParseUint(stringVal, 10, 16)

	if err != nil {
		return 0, fmt.Errorf("cannot extract a decimal constant (after trying hex and binary): %w", err)
	}

	return uint16(val), err
}

func ParseModuleString(module string) (*Module, error) {
	moduleScanner := bufio.NewScanner(strings.NewReader(module))

	lines := make([]Line, 0)

	for lineNum := 1; moduleScanner.Scan(); lineNum++ {
		line := moduleScanner.Text()
		parsedLine, err := parseAsmLine(lineNum, line)

		if err != nil {
			return nil, err
		}

		lines = append(lines, parsedLine)
	}

	asmModule := &Module{
		Lines: lines,
	}

	_ = spew.Sdump(asmModule)
	//spew.Dump(asmModule)

	return asmModule, nil
}

type lineToken struct {
	text      string
	tokenType rune
	position  scanner.Position
}

func parseAsmLine(lineNum int, line string) (Line, error) {
	var s scanner.Scanner
	lineReader := strings.NewReader(line)
	s.Init(lineReader)
	tokens := make([]lineToken, 0)

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokens = append(tokens, lineToken{
			text:      s.TokenText(),
			tokenType: tok,
			position:  s.Position,
		})
	}

	// fmt.Printf("Line %d, tokens (cnt: %d) %v\n", lineNum, len(tokens), tokens)

	lineMaker := func(asmLine AssemblyLine) Line {
		return Line{
			Number:  lineNum,
			Content: asmLine,
		}
	}

	if len(tokens) == 0 {
		return lineMaker(AssemblyBlankLine()), nil
	}

	if len(tokens) == 1 {
		return Line{}, fmt.Errorf("incomplete and possibly malformed line %d", lineNum)
	}

	if tokens[0].tokenType != scanner.Ident {
		return Line{}, fmt.Errorf("expected an identifier at line %d, column %d", lineNum, tokens[0].position.Column)
	}

	if tokens[1].text == "=" {
		if len(tokens) != 3 {
			return Line{}, fmt.Errorf("line %d looks like assignment, but expected it in 'symbol = value' format, got excess content on the line", lineNum)
		}

		if tokens[2].tokenType != scanner.Int {
			return Line{}, fmt.Errorf("line %d column %d, expected an int value", lineNum, tokens[2].position.Column)
		}

		return lineMaker(AssemblyAssignmentLine(
			Symbol(tokens[0].text),
			HardcodedValue(tokens[2].text),
		)), nil
	}

	if tokens[1].text == ":" {
		lineLabel := Label(tokens[0].text)

		if len(tokens) == 2 {
			return lineMaker(AssemblyLabelLine(lineLabel)), nil
		}

		instr, err := parseInstructionTokens(tokens[2:])

		if err != nil {
			return Line{}, fmt.Errorf("error with instruction on line %d: %w", lineNum, err)
		}

		return lineMaker(AssemblyLabeledInstructionLine(lineLabel, instr)), nil
	}

	instr, err := parseInstructionTokens(tokens)

	if err != nil {
		return Line{}, fmt.Errorf("error with instruction on line %d: %w", lineNum, err)
	}

	return lineMaker(AssemblyInstructionLine(instr)), nil
}

func parseInstructionTokens(tokens []lineToken) (Instruction, error) {
	instr, err := isa.StringToInstruction(tokens[0].text)

	if err != nil {
		return Instruction{}, fmt.Errorf("column %d, instruction parsing error: %w", tokens[0].position.Column, err)
	}

	rdToken := tokens[1]

	if rdToken.tokenType != scanner.Ident {
		return Instruction{}, fmt.Errorf("column %d, expected a register identifier", rdToken.position.Column)
	}

	rd, err := ParseRegister(rdToken.text)

	if err != nil {
		return Instruction{}, fmt.Errorf("column %d, cannot parse destination register: %w", rdToken.position.Column, err)
	}

	// Maps instructions to argument token types (excluding first argument which is always rd, a register; it has already been checked).
	// It's a list of lists: each element of this list is one of the options that the parser can accept.
	instrToArgTokens := map[isa.InstructionCode][][]rune{
		isa.ADD:  {{scanner.Ident, scanner.Ident}},
		isa.SUB:  {{scanner.Ident, scanner.Ident}},
		isa.LW:   {{scanner.Ident}},
		isa.SW:   {{scanner.Ident}},
		isa.XOR:  {{scanner.Ident, scanner.Ident}},
		isa.AND:  {{scanner.Ident, scanner.Ident}},
		isa.OR:   {{scanner.Ident, scanner.Ident}},
		isa.ADDI: {{scanner.Ident}, {scanner.Int}},
		isa.LDHI: {{scanner.Ident}, {scanner.Int}},
		isa.BZ:   {{scanner.Ident}, {scanner.Int}},
		isa.BNZ:  {{scanner.Ident}, {scanner.Int}},
		isa.JAL:  {{scanner.Ident}, {scanner.Int}},
		isa.JALR: {{scanner.Ident}},
		isa.SHL:  {{scanner.Int}},
		isa.SHR:  {{scanner.Int}},
		isa.SHRA: {{scanner.Int}},
	}

	remainingTokens := tokens[2:]
	remainingTokenOptions, ok := instrToArgTokens[instr]

	if !ok {
		// This should really never happen.
		return Instruction{}, fmt.Errorf("unable to parse different options for the instruction")
	}

	foundCompat := false

	for _, option := range remainingTokenOptions {
		compatible := len(option) == len(remainingTokens)

		if !compatible {
			continue
		}

		for i := range option {
			compatible = option[i] == remainingTokens[i].tokenType

			if !compatible {
				break
			}
		}

		if !compatible {
			continue
		}

		foundCompat = true
		break
	}

	if !foundCompat {
		// TODO: add more details, improve the error message
		return Instruction{}, fmt.Errorf("malformed instruction")
	}

	args := make([]InstructionArg, 0, len(remainingTokens))

	for _, token := range remainingTokens {
		var argType InstructionArgType

		switch token.tokenType {
		case scanner.Ident:
			argType = INSTRUCTION_ARG_TYPE_IDENTIFIER
		case scanner.Int:
			argType = INSTRUCTION_ARG_TYPE_NUMBER
		default:
			argType = INSTRUCTION_ARG_TYPE_UNKNOWN
		}

		args = append(args, InstructionArg{
			ArgType:          argType,
			UnprocessedValue: token.text,
		})
	}

	return Instruction{
		CpuInstruction: instr,
		Rd:             rd,
		Args:           args,
	}, nil
}

func ParseRegister(r string) (isa.RegisterId, error) {
	if (r[0] != 'r') && (r[0] != 'R') {
		return 0, fmt.Errorf("register reference must begin with r or R")
	}

	regId, err := strconv.Atoi(r[1:])

	if err != nil {
		return 0, fmt.Errorf("unable to parse the register number from '%s' (possibly illegal suffix and/or register number)", r)
	}

	if (regId < int(isa.MinRegId)) || (regId > int(isa.MaxRegId)) {
		return 0, fmt.Errorf("register IDs should be between %d and %d, got %d instead", isa.MinRegId, isa.MaxRegId, regId)
	}

	return isa.RegisterId(regId), nil
}
