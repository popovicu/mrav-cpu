load("@rules_go//go:def.bzl", "go_binary", "go_cross_binary")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_binary(
    name = "assembler",
    srcs = [
        "assembler.go",
    ],
    cgo = False,
    pure = "on",
    deps = [
        "//software/asm",
        "//software/format",
    ],
)

go_cross_binary(
    name = "assembler_wasm",
    platform = "//platforms:wasm_js",
    target = ":assembler",
)
