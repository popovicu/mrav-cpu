load("//software/build_defs:mrav.bzl", "mrav_binary", "mrav_library")

mrav_library(
    name = "zero",
    srcs = [
        "zero.mrav",
    ],
)

mrav_library(
    name = "base",
    srcs = [
        "base.mrav",
    ],
    deps = [
        ":zero",
    ],
)

mrav_library(
    name = "lib",
    srcs = [
        "lib.mrav",
    ],
    deps = [
        ":base",
    ],
)

mrav_binary(
    name = "program",
    srcs = [
        "program.mrav",
    ],
    out = "program.bin",
    deps = [
        ":base",
        ":lib",
        ":zero",
    ],
)
