load("//hardware/rtl/build:bundling.bzl", "generate_bundling_shell")
load("//hardware/rtl/mravbus/build_defs:bus.bzl", "BusDeviceStitchInfo")

MravTopInfo = provider(
    doc = "Provider for a Mrav top",
    fields = {
        "stitch_stage_output": "The preliminary output where the stitching is first added.",
    },
)

def gen_bus_stitch_rtl(dev):
    return "\n".join(["    .{bus_port}({top_connection}),".format(bus_port = key, top_connection = val) for (key, val) in dev[BusDeviceStitchInfo].bus_connections.items()])

def _mrav_top_impl(ctx):
    stitched_devices = ctx.attr.stitched_devices
    stitch_rtl = "\n\n".join([dev[BusDeviceStitchInfo].stitch_rtl for dev in stitched_devices])

    bus_rtl = "\n".join([gen_bus_stitch_rtl(dev) for dev in ctx.attr.stitched_devices])

    if len(bus_rtl) > 0:
        bus_rtl = bus_rtl[:-1]  # Remove the trailing comma
        bus_rtl = ",\n%s" % bus_rtl  # Deal with the preceeding comma

    top_template = ctx.file.top_rtl
    stitch_stage_file = ctx.actions.declare_file("%s_stitch_stage.sv" % ctx.attr.name)
    ctx.actions.expand_template(
        template = top_template,
        output = stitch_stage_file,
        substitutions = {
            "{stitchings}": stitch_rtl,
            "{bus_connections}": bus_rtl,
        },
    )

    output_file = ctx.outputs.out
    bundle_files = [ctx.file.core_rtl_bundle, ctx.file.bus_rtl] + [dev[BusDeviceStitchInfo].device_info.verilog for dev in ctx.attr.stitched_devices] + [stitch_stage_file]
    cmd = generate_bundling_shell(bundle_files, output_file)

    ctx.actions.run_shell(
        outputs = [output_file],
        inputs = bundle_files,
        command = cmd,
    )

    return [
        DefaultInfo(files = depset([stitch_stage_file, output_file])),
        MravTopInfo(stitch_stage_output = stitch_stage_file),
    ]

mrav_top = rule(
    implementation = _mrav_top_impl,
    attrs = {
        "bus_rtl": attr.label(
            mandatory = True,
            allow_single_file = [".sv"],
        ),
        "stitched_devices": attr.label_list(
            mandatory = True,
            providers = [BusDeviceStitchInfo],
        ),
        "core_rtl_bundle": attr.label(
            mandatory = True,
            allow_single_file = [".sv"],
        ),
        "top_rtl": attr.label(
            mandatory = True,
            allow_single_file = [".sv"],
        ),
        "out": attr.output(
            doc = "Output label for the stitch RTL code",
            mandatory = True,
        ),
    },
)
