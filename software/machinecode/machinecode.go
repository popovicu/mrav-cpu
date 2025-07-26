package machinecode

import (
	"bytes"
	"fmt"

	"mrav/isa"
	"mrav/software/model"
)

func GenerateMachineCode(module model.MravModule, output *bytes.Buffer) error {
	for _, instr := range module.Instructions {
		if err := GenerateMachineCodeForInstruction(instr, output); err != nil {
			return err
		}
	}

	return nil
}

func GenerateMachineCodeForInstruction(instr model.MravInstruction, output *bytes.Buffer) error {
	if instr.Add != nil {
		return generateAdd(instr.Add, output)
	}

	if instr.Sub != nil {
		return generateSub(instr.Sub, output)
	}

	if instr.Lw != nil {
		return generateLw(instr.Lw, output)
	}

	if instr.Sw != nil {
		return generateSw(instr.Sw, output)
	}

	if instr.Xor != nil {
		return generateXor(instr.Xor, output)
	}

	if instr.And != nil {
		return generateAnd(instr.And, output)
	}

	if instr.Or != nil {
		return generateOr(instr.Or, output)
	}

	if instr.Addi != nil {
		return generateAddi(instr.Addi, output)
	}

	if instr.Ldhi != nil {
		return generateLdhi(instr.Ldhi, output)
	}

	if instr.Bz != nil {
		return generateBz(instr.Bz, output)
	}

	if instr.Bnz != nil {
		return generateBnz(instr.Bnz, output)
	}

	if instr.Jal != nil {
		return generateJal(instr.Jal, output)
	}

	if instr.Jalr != nil {
		return generateJalr(instr.Jalr, output)
	}

	if instr.Shl != nil {
		return generateShl(instr.Shl, output)
	}

	if instr.Shr != nil {
		return generateShr(instr.Shr, output)
	}

	if instr.Shra != nil {
		return generateShra(instr.Shra, output)
	}

	return fmt.Errorf("unexpected instruction: %v", instr)
}

func merge4BitVals(hi byte, lo byte) byte {
	// Assuming both are 4 bit values
	return (hi << 4) | lo
}

func generateAdd(add *model.MravAdd, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.ADD), byte(add.Rd)), merge4BitVals(byte(add.Rs1), byte(add.Rs2))})

	if err != nil {
		return fmt.Errorf("cannot generate code for ADD: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for ADD, wrote %d instead", written)
	}

	return nil
}

func generateSub(sub *model.MravSub, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.SUB), byte(sub.Rd)), merge4BitVals(byte(sub.Rs1), byte(sub.Rs2))})

	if err != nil {
		return fmt.Errorf("cannot generate code for SUB: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for SUB, wrote %d instead", written)
	}

	return nil
}

func generateLw(lw *model.MravLw, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.LW), byte(lw.Rd)), merge4BitVals(byte(lw.Rs1), byte(0x0))})

	if err != nil {
		return fmt.Errorf("cannot generate code for LW: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for LW, wrote %d instead", written)
	}

	return nil
}

func generateSw(sw *model.MravSw, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.SW), byte(sw.Rd)), merge4BitVals(byte(sw.Rs1), byte(0x0))})

	if err != nil {
		return fmt.Errorf("cannot generate code for SW: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for SW, wrote %d instead", written)
	}

	return nil
}

func generateXor(xor *model.MravXor, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.XOR), byte(xor.Rd)), merge4BitVals(byte(xor.Rs1), byte(xor.Rs2))})

	if err != nil {
		return fmt.Errorf("cannot generate code for XOR: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for XOR, wrote %d instead", written)
	}

	return nil
}

func generateAnd(and *model.MravAnd, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.AND), byte(and.Rd)), merge4BitVals(byte(and.Rs1), byte(and.Rs2))})

	if err != nil {
		return fmt.Errorf("cannot generate code for AND: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for AND, wrote %d instead", written)
	}

	return nil
}

func generateOr(or *model.MravOr, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.OR), byte(or.Rd)), merge4BitVals(byte(or.Rs1), byte(or.Rs2))})

	if err != nil {
		return fmt.Errorf("cannot generate code for OR: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for OR, wrote %d instead", written)
	}

	return nil
}

func generateAddi(addi *model.MravAddi, output *bytes.Buffer) error {
	if addi.Value.IsRight() {
		return fmt.Errorf("cannot generate machine code for ADDI, still pointing to a symbol '%s'", addi.Value.MustRight())
	}

	err := output.WriteByte(merge4BitVals(byte(isa.ADDI), byte(addi.Rd)))

	if err != nil {
		return fmt.Errorf("cannot write the first byte of ADDI: %w", err)
	}

	err = output.WriteByte(addi.Value.MustLeft())

	if err != nil {
		return fmt.Errorf("cannot write the immediate value of ADDI: %v", err)
	}

	return nil
}

func generateLdhi(ldhi *model.MravLdhi, output *bytes.Buffer) error {
	if ldhi.Value.IsRight() {
		return fmt.Errorf("cannot generate machine code for LDHI, still pointing to a symbol '%s'", ldhi.Value.MustRight())
	}

	err := output.WriteByte(merge4BitVals(byte(isa.LDHI), byte(ldhi.Rd)))

	if err != nil {
		return fmt.Errorf("cannot write the first byte of LDHI: %w", err)
	}

	err = output.WriteByte(ldhi.Value.MustLeft())

	if err != nil {
		return fmt.Errorf("cannot write the immediate value of LDHI: %v", err)
	}

	return nil
}

func generateBz(bz *model.MravBz, output *bytes.Buffer) error {
	if bz.Addr.IsRight() {
		return fmt.Errorf("cannot generate machine code for BZ, still pointing to a symbol '%s'", bz.Addr.MustRight())
	}

	err := output.WriteByte(merge4BitVals(byte(isa.BZ), byte(bz.Rd)))

	if err != nil {
		return fmt.Errorf("cannot write the first byte of BZ: %w", err)
	}

	err = output.WriteByte(bz.Addr.MustLeft())

	if err != nil {
		return fmt.Errorf("cannot write the immediate value of BZ: %v", err)
	}

	return nil
}

func generateBnz(bnz *model.MravBnz, output *bytes.Buffer) error {
	if bnz.Addr.IsRight() {
		return fmt.Errorf("cannot generate machine code for BNZ, still pointing to a symbol '%s'", bnz.Addr.MustRight())
	}

	err := output.WriteByte(merge4BitVals(byte(isa.BNZ), byte(bnz.Rd)))

	if err != nil {
		return fmt.Errorf("cannot write the first byte of BNZ: %w", err)
	}

	err = output.WriteByte(bnz.Addr.MustLeft())

	if err != nil {
		return fmt.Errorf("cannot write the immediate value of BNZ: %v", err)
	}

	return nil
}

func generateJal(jal *model.MravJal, output *bytes.Buffer) error {
	if jal.Addr.IsRight() {
		return fmt.Errorf("cannot generate machine code for JAL, still pointing to a symbol '%s'", jal.Addr.MustRight())
	}

	err := output.WriteByte(merge4BitVals(byte(isa.JAL), byte(jal.Rd)))

	if err != nil {
		return fmt.Errorf("cannot write the first byte of JAL: %w", err)
	}

	err = output.WriteByte(jal.Addr.MustLeft())

	if err != nil {
		return fmt.Errorf("cannot write the immediate value of JAL: %v", err)
	}

	return nil
}

func generateJalr(jalr *model.MravJalr, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.JALR), byte(jalr.Rd)), merge4BitVals(byte(jalr.Rs1), 0x0)})

	if err != nil {
		return fmt.Errorf("cannot generate code for JALR: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for JALR, wrote %d instead", written)
	}

	return nil
}

func generateShl(shl *model.MravShl, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.SHL), byte(shl.Rd)), merge4BitVals(byte(shl.Imm4), 0x0)})

	if err != nil {
		return fmt.Errorf("cannot generate code for SHL: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for SHL, wrote %d instead", written)
	}

	return nil
}

func generateShr(shr *model.MravShr, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.SHR), byte(shr.Rd)), merge4BitVals(byte(shr.Imm4), 0x0)})

	if err != nil {
		return fmt.Errorf("cannot generate code for SHR: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for SHR, wrote %d instead", written)
	}

	return nil
}

func generateShra(shra *model.MravShra, output *bytes.Buffer) error {
	written, err := output.Write([]byte{merge4BitVals(byte(isa.SHRA), byte(shra.Rd)), merge4BitVals(byte(shra.Imm4), 0x0)})

	if err != nil {
		return fmt.Errorf("cannot generate code for SHRA: %w", err)
	}

	if written != 2 {
		return fmt.Errorf("expected to write 2 bytes for SHRA, wrote %d instead", written)
	}

	return nil
}
