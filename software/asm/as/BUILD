load("@rules_go//go:def.bzl", "go_binary", "go_cross_binary")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_binary(
    name = "as",
    srcs = [
        "as.go",
    ],
    cgo = False,
    pure = "on",
    deps = [
        "//remote/spew",
        "//software/asm",
        "//software/format",
    ],
)

go_cross_binary(
    name = "as_x86_64",
    platform = "//platforms:x86_64_linux",
    target = ":as",
)

go_cross_binary(
    name = "as_js_wasm",
    platform = "//platforms:wasm_js",
    target = ":as",
)
