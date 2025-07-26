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
        "bus_tpl.sv",
    ],
    importpath = "mrav/hardware/rtl/mravbus/generator",
)
