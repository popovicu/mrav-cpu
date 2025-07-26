load("@rules_go//go:def.bzl", "go_binary", "go_cross_binary")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_binary(
    name = "codegen",
    srcs = [
        "codegen.go",
    ],
    cgo = False,
    pure = "on",
    deps = [
        "//hardware/rtl/mravbus/generator",
    ],
)

go_cross_binary(
    name = "codegen_x86_64",
    platform = "//platforms:x86_64_linux",
    target = ":codegen",
)
