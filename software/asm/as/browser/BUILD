load("@bazel_skylib//rules:copy_file.bzl", "copy_file")
load("@rules_go//go:def.bzl", "go_binary")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

copy_file(
    name = "wasm_exec_copy",
    src = "//third_party:wasm_exec.js",
    out = "wasm_exec.js",
)

copy_file(
    name = "assembler_copy",
    src = "//software/asm/as/browser/assembler:assembler_wasm",
    out = "assembler.wasm",
)

go_binary(
    name = "browser",
    srcs = [
        "browser.go",
    ],
    cgo = False,
    embedsrcs = [
        "as.html",
        ":assembler.wasm",
        ":wasm_exec.js",
    ],
    pure = "on",
)
