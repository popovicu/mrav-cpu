load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "device",
    srcs = [
        "device.go",
    ],
    importpath = "mrav/system/easybus/device",
    deps = [
        "//isa",
    ],
)
