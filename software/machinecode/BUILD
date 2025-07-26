load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "machinecode",
    srcs = [
        "machinecode.go",
    ],
    importpath = "mrav/software/machinecode",
    deps = [
        "//isa",
        "//software/model",
    ],
)
