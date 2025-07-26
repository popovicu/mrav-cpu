MravLibInfo = provider(
    doc = "Provider for the Mrav assembly library",
    fields = {
        "modules": "[list[File]] ID name of the build environment",
        "deps": "[depset[Target]] depset of all the dependencies",
    },
)

def _mrav_library_impl(ctx):
    module_srcs = ctx.files.srcs
    deps = depset(ctx.attr.deps, transitive = [dep[MravLibInfo].deps for dep in ctx.attr.deps])
    return [
        MravLibInfo(modules = module_srcs, deps = deps),
        DefaultInfo(files = depset(module_srcs)),
    ]

mrav_library = rule(
    implementation = _mrav_library_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = [".mrav"],
        ),
        "deps": attr.label_list(
            providers = [MravLibInfo],
        ),
    },
    provides = [MravLibInfo],
)

def _mrav_binary_impl(ctx):
    output_image = ctx.outputs.out
    assembler = ctx.executable.assembler
    format = ctx.attr.format

    all_srcs = []
    seen_files = set()

    # Add direct srcs
    for s in ctx.files.srcs:
        if s.path in seen_files:
            fail("Duplicate source file found in mrav_binary srcs: %s" % s.path)
        else:
            all_srcs.append(s)
            seen_files.add(s.path)

    final_depset = depset(ctx.attr.deps, transitive = [dep[MravLibInfo].deps for dep in ctx.attr.deps])

    # Add srcs from library dependencies
    for dep in final_depset.to_list():
        if MravLibInfo in dep:
            for module in dep[MravLibInfo].modules:
                if module.path in seen_files:
                    fail("Duplicate source file found in mrav_binary deps: %s" % module.path)
                else:
                    all_srcs.append(module)
                    seen_files.add(module.path)

    ctx.actions.run(
        inputs = all_srcs,
        outputs = [output_image],
        arguments = ["--output", output_image.path, "--format", format] + [m.path for m in all_srcs],
        executable = assembler,
        progress_message = "Running Mrav assembler",
    )

    return [
        DefaultInfo(files = depset([output_image])),
    ]

mrav_binary = rule(
    implementation = _mrav_binary_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = [".mrav"],
        ),
        "deps": attr.label_list(
            providers = [MravLibInfo],
        ),
        "assembler": attr.label(
            default = Label("//software/asm/as"),
            allow_files = True,
            executable = True,
            cfg = "exec",
        ),
        "out": attr.output(
            doc = "Output label for the Mrav image",
            mandatory = True,
        ),
        "format": attr.string(
            default = "binary",
            values = ["human", "binary"],
        ),
    },
)
