load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("@rules_python//python:defs.bzl", "py_test")
load("//software/build_defs:mrav.bzl", "mrav_binary")

mrav_binary(
    name = "program",
    srcs = [
        "program.mrav",
    ],
    out = "program.bin",
)

py_test(
    name = "equivalence_test",
    srcs = ["equivalence_test.py"],
    data = [
        ":program.bin",
        ":program_sim_proto.pb",
        "//hardware/rtl:mrav_core.sv",
    ],
    env = {
        "CORE_VERILOG": "$(location //hardware/rtl:mrav_core.sv)",
        "TEST_SOFTWARE": "$(location :program.bin)",
        "SOFTWARE_CPU_STATE": "$(location :program_sim_proto.pb)",
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
    name = "equivalence_sim_run",
    srcs = [":program.bin"],
    outs = [
        ":program_sim_proto.pb",
        ":program_sim_state.txt",
    ],
    args = [
        "--software=$(location :program.bin)",
        "--instructions_to_sim=45",
        "--core_state_output=$(location :program_sim_state.txt)",
        "--core_state_proto_output=$(location :program_sim_proto.pb)",
    ],
    tool = "//system/binaries/memonly",
)
