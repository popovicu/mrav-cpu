load("//hardware/rtl/build:bundling.bzl", "generate_bundling_shell")

def _system_verilog_bundle_impl(ctx):
    output_bundle = ctx.outputs.out
    cmd = generate_bundling_shell(ctx.files.srcs, output_bundle)

    ctx.actions.run_shell(
        outputs = [output_bundle],
        inputs = ctx.files.srcs,
        command = cmd,
    )

    return [DefaultInfo(files = depset([output_bundle]))]

system_verilog_bundle = rule(
    implementation = _system_verilog_bundle_impl,
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = [".sv"],
        ),
        "out": attr.output(
            doc = "Output label for the stitch RTL code",
            mandatory = True,
        ),
    },
)
