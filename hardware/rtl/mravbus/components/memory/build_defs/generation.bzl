def _mrav_imaged_memory_impl(ctx):
    output_file = ctx.outputs.out
    image_file = ctx.file.image

    ctx.actions.run(
        outputs = [output_file],
        inputs = [image_file],
        executable = ctx.executable.codegen,
        arguments = ["--rtl_file", output_file.path, "--mem_size", str(ctx.attr.memory_size), "--bin_payload", image_file.path, "--top_module", ctx.attr.top_module, "--internal_address", ctx.attr.internal_address_generation],
        progress_message = "Running memory codegen",
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

mrav_imaged_memory = rule(
    implementation = _mrav_imaged_memory_impl,
    attrs = {
        "memory_size": attr.int(
            mandatory = True,
        ),
        "internal_address_generation": attr.string(
            mandatory = True,
        ),
        "image": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "top_module": attr.string(
            mandatory = True,
        ),
        "out": attr.output(
            doc = "Output label for the RTL code",
            mandatory = True,
        ),
        "codegen": attr.label(
            default = Label("//hardware/rtl/mravbus/components/memory/generator/codegen"),
            allow_files = True,
            executable = True,
            cfg = "exec",
        ),
    },
)
