load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("@rules_python//python:defs.bzl", "py_test")
load("//software/build_defs:mrav.bzl", "mrav_binary")

mrav_binary(
    name = "functions",
    srcs = [
        "functions.mrav",
    ],
    out = "functions.bin",
)

py_test(
    name = "functions_test",
    srcs = ["functions_test.py"],
    data = [
        ":functions.bin",
        ":functions_sim_proto.pb",
        "//hardware/rtl:mrav_core.sv",
    ],
    env = {
        "CORE_VERILOG": "$(location //hardware/rtl:mrav_core.sv)",
        "TEST_SOFTWARE": "$(location :functions.bin)",
        "SOFTWARE_CPU_STATE": "$(location :functions_sim_proto.pb)",
    },
    deps = [
        "//hardware/testbench/components:mrav_bus_memory",
        "//hardware/testbench/core",
        "//hardware/testbench/simulation",
        "//remote/cocotb",
        "//remote/pytest",
    ],
)

run_binary(
    name = "functions_sim_run",
    srcs = [":functions.bin"],
    outs = [
        ":functions_sim_proto.pb",
        ":functions_sim_state.txt",
    ],
    args = [
        "--software=$(location :functions.bin)",
        "--instructions_to_sim=150",
        "--core_state_output=$(location :functions_sim_state.txt)",
        "--core_state_proto_output=$(location :functions_sim_proto.pb)",
    ],
    tool = "//system/binaries/memonly",
)
