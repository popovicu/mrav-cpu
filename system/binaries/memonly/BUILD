load("@rules_go//go:def.bzl", "go_binary", "go_cross_binary")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_binary(
    name = "memonly",
    srcs = [
        "memonly.go",
    ],
    cgo = False,
    pure = "on",
    deps = [
        "//isa",
        "//software/asm",
        "//system",
        "//system/easybus",
        "//system/easybus/device",
        "//system/easybus/device/memory",
        "//system/easybus/device/timer",
    ],
)

go_cross_binary(
    name = "memonly_x86_64",
    platform = "//platforms:x86_64_linux",
    target = ":memonly",
)
