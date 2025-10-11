# mrav-cpu

Mrav is a set of tools and libraries for deploying a tiny, minimal CPU core in a microcontroller-like setting. This monorepo enables you to generate the RTL code, run software simulations, and produce software for Mrav with maximum automation.

> :warning: This project is not ready for any sort of production use at this point.
> Use with caution for any non-experimental settings, and please file feature requests for any features you'd like to see for production use.
> At the moment, consider this codebase a proof of concept only. ISA is still subject to change.

The build system used in this project is Bazel, and it covers everything from RTL generation, to running tests, to building the assembler and assembling the Mrav software itself.

> :bulb: Concepts and implementation from this project will be documented in detail at [this page](https://popovicu.com/mrav-cpu/). This will likely never be a production-ready design, and it is more meant to illustrate some of the computer engineering ideas the author has. For a better ecosystem, consider using lightweight RISC-V cores, such as `rv32e`-based designs.

The name Mrav comes from the author's native language, Serbian, and it means 'ant'. This CPU is tiny like an ant, but capable of doing work. :)

## Try it in browser first!

> :bulb: NEW: Browser-based playground is available at [mrav-playground.popovicu.com](https://mrav-playground.popovicu.com)! Try Mrav in the browser first to see if you like it, no downloads required.

## Remote dependency

If you don't want to vendor this code into your codebase, you can depend on this repository as a remote dependency in Bazel by using `bazel_dep` in your `MODULE.bazel`. However, because `mrav-cpu` is currently not tracked in the Bazel registry system, you also need to add an explicit Git pointer to this repository, like this:

```
git_override(
    module_name = "mrav_cpu",
    commit = "COMMIT_HASH",
    remote = "https://github.com/popovicu/mrav-cpu",
)
```

## ISA

There are 16 instructions in the Mrav ISA, and each instruction is 16 bits. The core accesses 16-bit memory words only. Address alignment is not enforced for simplicity.

The highest 4 bits encode the instruction. The next 4 bits encode the destination register. The final 8 bits may either be two source register IDs, an 8-bit immediate value, a source register ID + 4 unused bits, or an immediate 4-bit value + 4 unused bits.

```
╔════════════════════════════════════════════════════════════════════╗
║                      MRAV CPU INSTRUCTION SET                      ║
║                         16 Instructions                            ║
╠════════════════════════════════════════════════════════════════════╣
║ ARITHMETIC & LOGIC                                                 ║
║  add  rd rs1 rs2    Add rs1 + rs2 → rd                             ║
║  sub  rd rs1 rs2    Subtract rs1 - rs2 → rd                        ║
║  xor  rd rs1 rs2    XOR rs1 ^ rs2 → rd                             ║
║  and  rd rs1 rs2    AND rs1 & rs2 → rd                             ║
║  or   rd rs1 rs2    OR rs1 | rs2 → rd                              ║
║  addi rd imm8       Add immediate rd + imm8 → rd                   ║
╠════════════════════════════════════════════════════════════════════╣
║ MEMORY ACCESS                                                      ║
║  lw   rd rs1        Load word from [rs1] → rd                      ║
║  sw   rd rs1        Store word rd → [rs1]                          ║
╠════════════════════════════════════════════════════════════════════╣
║ IMMEDIATE LOAD                                                     ║
║  ldhi rd imm8       Load high immediate imm8 → rd[15:8]            ║
╠════════════════════════════════════════════════════════════════════╣
║ BRANCHES (absolute addressing)                                     ║
║  bz   rd imm8       Branch if rd == 0 to address imm8              ║
║  bnz  rd imm8       Branch if rd != 0 to address imm8              ║
╠════════════════════════════════════════════════════════════════════╣
║ JUMPS                                                              ║
║  jal  rd imm8       Jump to imm8, save return addr → rd            ║
║  jalr rd rs1        Jump to [rs1], save return addr → rd           ║
╠════════════════════════════════════════════════════════════════════╣
║ SHIFTS                                                             ║
║  shl  rd imm4       Shift left rd << imm4 → rd                     ║
║  shr  rd imm4       Shift right (logical) rd >> imm4 → rd          ║
║  shra rd imm4       Shift right (arithmetic) rd >> imm4 → rd       ║
╠════════════════════════════════════════════════════════════════════╣
║ NOTES:                                                             ║
║  • 16-bit instruction width, 16-bit data width                     ║
║  • 16 general purpose registers (r0-r15)                           ║
║  • imm8 = 8-bit immediate, imm4 = 4-bit immediate                  ║
║  • Branch/jump immediates are absolute addresses                   ║
╚════════════════════════════════════════════════════════════════════╝
```

Additional description is at [the ISA doc](/docs/isa.md)

> :warning: Reiterating from the top: the ISA is subject to change. The author is well aware that this ISA is suboptimal in many ways.

## Interrupts

At the moment, no interrupts are supported for simplicity, and interaction with the world outside of the core should be done via polling.

## RTL

Currently, there is only one 'SoC'-like setup, simply called 'small'. To instantiate it, use the `mrav_small` macro:

```
load("//hardware/rtl/soc/build_defs:mrav_small.bzl", "mrav_small")

mrav_small(
    name = "soc",
    software = ":software.bin",
    gpio_verilog = "//hardware/rtl/soc:gpio.sv",
    soc_top = "//hardware/rtl/soc:soc.sv",
)
```

This will provide a SystemVerilog bundle file which has no dependencies on anything else. Additionally, the bundle will contain the RAM memory RTL preloaded with the software machine code (`software.bin` file). Therefore, this bundle will be ready to deploy as a soft core in FPGA.

Check `//deployments/led/BUILD` file for a concrete example. You can build that RTL bundle by running:

```
bazel build //deployments/led:soc_bundle.sv
```

For this concrete reference, this is how the bundle was used in the final SystemVerilog module:

```verilog
module btn_led(
    input sysclk,
    input [1:0]btn,   // Button inputs
    output [1:0]led  // Led outputs
    );
    
    wire [7:0] gpio;

    mrav_soc mrav_soc_i(
        .clk(sysclk),
        .rst_n(!btn[1]),
        .gpio(gpio)
    );

    assign led[1] = gpio[1];
    assign led[0] = gpio[0];
    
endmodule
```

## Software

Building software for Mrav does not require anything outside the Bazel build system. The tooling is based on Go. When you build a Mrav software target, Bazel will dynamically fetch the Go toolchain, build the necessary tools from this repository (e.g. the assembler) and then use the newly built assembler to produce the final software binary. This should all be transparent to the user.

For example, you can build a simple example like this:

```
bazel build //software/examples/adding
```

For that particular example, you can also build:

```
bazel build //software/examples/adding:adding_mrav_state.txt
```

which will build the binary from above, run it in the simulator and then dump the simulated core's state into a text file.

### Libraries

Mrav software also supports a simple form of a library system for the assembly files, and the example can be found here:

```
//software/examples/libraries/BUILD
```

An implicit benefit here is that you can easily import someone else's library using Bazel's remote dependency system, like this:

```
mrav_binary(
    name = "program",
    srcs = [
        "program.mrav",
    ],
    deps = [
        "@remote_dep_module//something/foo:bar_lib",
    ],
    out = "program.bin",
)
```

## Software simulation

The Go code for Mrav's simulation is in the `//core` Bazel package.

That package also includes the logic for binary serialization of the core's state. This comes in useful for comparing the core state between different representations, such as Go and Python (used for RTL simulation).

Some simple SoC-like systems are emulated in the `//system` Bazel package and subpackages. `//system/binaries` contains Go binaries for running full system simulations. An example is `memonly` which simply consists of a Mrav core and RAM memory attached to the virtual bus.

Go implementation is very portable, and can run in many contexts, including simply running the core inside a browser simulation, which is explained in more detail below.

## RTL simulation & equivalence tests

The RTL is heavily tested, and the simulation is done through Python `cocotb` library which drives `verilator`.

> :warning: You need Verilator in one of the standard locations like `/usr/bin`, `/usr/local/bin`, etc. in order to be able to run Bazel tests.

The main simulation test is in `hardware/rtl/core_test.py`. To run this test, you can run the following:

```
bazel test --sandbox_writable_path=$HOME/.cache/ccache --test_output=all //hardware/rtl:core_test
```

The writable path flag for `ccache` is needed because of the `cocotb` behavior.

Running the command above will transparently fetch the Go toolchain, build the assembler, assemble the software, create a simple system with the core + memory attached to the bus, preload the memory with the software, and run the simulation for the system using `verilator`.

A related question arises about confirming the equivalence between the Go simulation and the RTL cycle-for-cycle simulation. There is a dedicated test for that in `testing/equivalence/equivalence_test.py`. To run this, execute the following:

```
bazel test --sandbox_writable_path=$HOME/.cache/ccache --test_output=all //testing/equivalence:equivalence_test
```

This test is similar to the core test described above, but additionally builds and runs the Go simulation, dumping the core state to a binary protobuf representation, which is then diff'ed with the Python representation obtained through cycle-level simulation of the RTL itself.

## Portability & web browser environment

The Mrav components and tools are designed to be as portable as possible, and one of the objectives was to enable running in many contexts, including the browser.

This is one of the major reasons why Go was chosen for the software stack and tooling implementation. Go is extremely well supported by the Bazel build system and enables trivial cross-compilation to different platforms.

As an example, it is possible to run the Mrav assembler in browser, by building a WASM binary. To run a sample web server for this, execute the following:

```
bazel run //software/asm/as/browser
```

This will start a web server on `localhost:9876`. There should be a page available on http://localhost:9876/as.html with a text field to write a Mrav program and assemble it to machine code (the output will be in hex).

Other tools from the Mrav suite can similarly be ported to the browser, like the system simulation. This can be particularly useful for educational purposes.