load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("//software/build_defs:mrav.bzl", "mrav_binary")

mrav_binary(
    name = "timing",
    srcs = [
        "timing.mrav",
    ],
    out = "timing.bin",
    #format = "human",
)

run_binary(
    name = "timing_run",
    srcs = [":timing.bin"],
    outs = [":timing_state.txt"],
    args = [
        "--software=$(location :timing.bin)",
        "--instructions_to_sim=45",
        "--core_state_output=$(location :timing_state.txt)",
    ],
    tool = "//system/binaries/memonly",
)
