load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "firstpass",
    srcs = [
        "firstpass.go",
    ],
    importpath = "mrav/software/asm/firstpass",
    deps = [
        "//isa",
        "//software/asm/parsing",
        "//software/model",
    ],
)
