load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "asm",
    srcs = [
        "asm.go",
    ],
    importpath = "mrav/software/asm",
    deps = [
        "//isa",
        "//software/asm/firstpass",
        "//software/asm/parsing",
        "//software/asm/secondpass",
        "//software/linker",
        "//software/model",
    ],
)
