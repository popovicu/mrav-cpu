load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "secondpass",
    srcs = [
        "secondpass.go",
    ],
    importpath = "mrav/software/asm/secondpass",
    deps = [
        "//isa",
        "//software/model",
    ],
)
