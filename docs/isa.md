# ISA

*This page only contains the information that is additional to the ISA information from the README file.*

Most of the instructions should be self-explanatory. For example `add`/`sub`/`xor`/`and`/`or` will take two source registers, `rs1` and `rs2`, do the operation between the two, and store the result in the destination register, `rd`.

Some instructions like `addi` take the destination register and an 8-bit immediate value, so the operation is effectively `r[rd] += imm8`.

The immediate value in branching instructions like `bz` and `bnz` are absolute jump addresses. It's easy to implement, but a bit unwieldy from the software perspective, so jumping to addresses outside the 8-bit range is not straightforward and takes multiple instructions. This also applies to `jal`.

For `jal` and `jalr`, the destination register `rd` takes the address of the instruction right after the jump instruction, the "return" address. In `jalr` case, instead of the 8-bit absolute address, the jump address is read from the `rs1` source register.