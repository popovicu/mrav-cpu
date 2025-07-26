load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "linker",
    srcs = [
        "linker.go",
    ],
    importpath = "mrav/software/linker",
    deps = [
        "//isa",
        "//software/asm/secondpass",
        "//software/model",
    ],
)
