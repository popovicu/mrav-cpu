load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "format",
    srcs = [
        "binary.go",
        "human.go",
    ],
    importpath = "mrav/software/format",
    deps = [
        "//isa",
        "//software/machinecode",
        "//software/model",
    ],
)
