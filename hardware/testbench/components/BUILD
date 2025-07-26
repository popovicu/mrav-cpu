load("@rules_python//python:defs.bzl", "py_library")

package(
    default_visibility = [
        "//hardware/rtl:__subpackages__",
        "//testing:__subpackages__",
    ],
)

py_library(
    name = "mrav_bus_memory",
    srcs = ["mrav_bus_memory.py"],
    deps = [
        "@pip//cocotb:pkg",
    ],
)
