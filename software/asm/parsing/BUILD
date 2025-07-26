load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "parsing",
    srcs = [
        "parsing.go",
    ],
    importpath = "mrav/software/asm/parsing",
    deps = [
        "//isa",
        "//remote/mo",
        "//remote/spew",
    ],
)
