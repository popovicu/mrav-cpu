load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "timer",
    srcs = [
        "timer.go",
    ],
    importpath = "mrav/system/easybus/device/timer",
    deps = [
        "//isa",
        "//system",
    ],
)
