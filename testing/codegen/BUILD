load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("//hardware/rtl/mravbus/build_defs:bus.bzl", "mrav_bus_device")
load("//hardware/rtl/mravbus/build_defs:software.bzl", "mrav_bus_addr_module")
load("//software/build_defs:mrav.bzl", "mrav_binary")

mrav_bus_device(
    name = "nop",
    addr_hi = 0x00AA,
    addr_lo = 0x00AA,
    device_id = "nop",
    top = "nop",
    verilog = "nop_device.sv",
)

mrav_bus_addr_module(
    name = "nop_lib",
    out = ":nop_lib.mrav",
    bus_device = ":nop",
    symbol_prefix = "NOP",
)

mrav_binary(
    name = "codegen",
    srcs = [
        "codegen.mrav",
    ],
    out = "codegen.bin",
    deps = [
        ":nop_lib",
    ],
)

run_binary(
    name = "codegen_run_simple",
    srcs = [":codegen.bin"],
    outs = [":codegen_mrav_state.txt"],
    args = [
        "--software=$(location :codegen.bin)",
        "--instructions_to_sim=45",
        "--core_state_output=$(location :codegen_mrav_state.txt)",
    ],
    tool = "//system/binaries/memonly",
)
