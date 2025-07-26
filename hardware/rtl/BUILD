load("@rules_python//python:defs.bzl", "py_test")
load("//hardware/rtl/verilog/build_defs:system_verilog.bzl", "system_verilog_bundle")
load("//software/build_defs:mrav.bzl", "mrav_binary")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

py_test(
    name = "core_test",
    srcs = ["core_test.py"],
    data = [
        ":mrav_core.sv",
        ":program.bin",
    ],
    env = {
        "CORE_VERILOG": "$(location :mrav_core.sv)",
        "TEST_SOFTWARE": "$(location :program.bin)",
    },
    deps = [
        "//hardware/testbench/components:mrav_bus_memory",
        "//hardware/testbench/core",
        "//hardware/testbench/simulation",
        "//remote/cocotb",
        "//remote/pytest",
    ],
)

mrav_binary(
    name = "program",
    srcs = [
        "program.mrav",
    ],
    out = "program.bin",
)

system_verilog_bundle(
    name = "mrav_core_bundle",
    srcs = [
        "constants.sv",
        "core.sv",
    ],
    out = ":mrav_core.sv",
)
