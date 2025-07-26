load("@rules_python//python:defs.bzl", "py_library")

package(
    default_visibility = [
        "//hardware:__subpackages__",
        "//testing:__subpackages__",
    ],
)

py_library(
    name = "simulation",
    srcs = ["simulation.py"],
    deps = [
        "//core/proto:core_py_pb2",
        "//remote/py_protobuf",
    ],
)
