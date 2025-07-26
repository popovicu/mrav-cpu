load("@rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_python//python:proto.bzl", "py_proto_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

proto_library(
    name = "core_proto",
    srcs = ["core.proto"],
)

go_proto_library(
    name = "core_go_proto",
    importpath = "mrav/core/proto",
    proto = ":core_proto",
)

py_proto_library(
    name = "core_py_pb2",
    deps = [":core_proto"],
)
