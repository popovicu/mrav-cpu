load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_library(
    name = "core",
    srcs = [
        "core.go",
        "serialization.go",
    ],
    importpath = "mrav/core",
    deps = [
        "//core/proto:core_go_proto",
        "//isa",
        "//remote/protobuf",
    ],
)
