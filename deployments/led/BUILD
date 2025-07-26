load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("//hardware/rtl/soc/build_defs:mrav_small.bzl", "mrav_small")
load("//software/build_defs:mrav.bzl", "mrav_binary")

mrav_binary(
    name = "sw",
    srcs = [
        "sw.mrav",
    ],
    out = "sw.bin",
    deps = [
        ":soc_gpio_lib",
    ],
)

mrav_small(
    name = "soc",
    gpio_verilog = "//hardware/rtl/soc:gpio.sv",
    soc_top = "//hardware/rtl/soc:soc.sv",
    software = ":sw.bin",
)

run_binary(
    name = "sw_run_simple",
    srcs = [":sw.bin"],
    outs = [":sw_mrav_state.txt"],
    args = [
        "--software=$(location :sw.bin)",
        "--instructions_to_sim=1000000",
        "--core_state_output=$(location :sw_mrav_state.txt)",
        "--verbose",
    ],
    tool = "//system/binaries/memonly",
)
