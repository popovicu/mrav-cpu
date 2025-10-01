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
        "//system",
        "//system/easybus",
        "//system/easybus/device",
        "//system/easybus/device/memory",
    ],
)

go_cross_binary(
    name = "memonly_wasm",
    platform = "//platforms:wasm_js",
    target = ":memonly",
)