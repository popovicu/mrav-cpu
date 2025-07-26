load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "generator",
    srcs = [
        "generator.go",
    ],
    embedsrcs = [
        "memory.sv",
    ],
    importpath = "mrav/hardware/rtl/mravbus/components/memory/generator",
)
