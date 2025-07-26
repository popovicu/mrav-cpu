load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "easybus",
    srcs = [
        "easybus.go",
    ],
    importpath = "mrav/system/easybus",
    deps = [
        "//core",
        "//isa",
        "//system",
        "//system/easybus/device",
    ],
)
