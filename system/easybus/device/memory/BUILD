load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "memory",
    srcs = [
        "mem.go",
    ],
    importpath = "mrav/system/easybus/device/memory",
    deps = [
        "//isa",
        "//system",
    ],
)
