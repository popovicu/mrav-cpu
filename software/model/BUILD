load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "model",
    srcs = [
        "instructions.go",
        "model.go",
    ],
    importpath = "mrav/software/model",
    deps = [
        "//isa",
        "//remote/mo",
    ],
)
