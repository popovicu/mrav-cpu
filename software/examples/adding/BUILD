load("@bazel_skylib//rules:run_binary.bzl", "run_binary")
load("//software/build_defs:mrav.bzl", "mrav_binary")

mrav_binary(
    name = "adding",
    srcs = [
        "adding.mrav",
        "memory.mrav",
    ],
    out = "adding.bin",
)

run_binary(
    name = "adding_run_simple",
    srcs = [":adding.bin"],
    outs = [":adding_mrav_state.txt"],
    args = [
        "--software=$(location :adding.bin)",
        "--instructions_to_sim=45",
        "--core_state_output=$(location :adding_mrav_state.txt)",
    ],
    tool = "//system/binaries/memonly",
)
