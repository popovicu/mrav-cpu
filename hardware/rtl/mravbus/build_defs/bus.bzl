BusDeviceInfo = provider(
    doc = "Provider for a Mrav bus device",
    fields = {
        "device_id": "Device ID",
        "addr_lo": "Low address boundary on the bus",
        "addr_hi": "High address boundary on the bus",
        "top": "Top level module",
        "verilog": "Verilog RTL",
    },
)

def _mrav_bus_device_impl(ctx):
    if ctx.attr.addr_lo > ctx.attr.addr_hi:
        fail("Low address is higher than high address: lo - %d, hi - %d" % (ctx.attr.addr_lo, ctx.attr.addr_hi))

    return BusDeviceInfo(
        device_id = ctx.attr.device_id,
        addr_lo = ctx.attr.addr_lo,
        addr_hi = ctx.attr.addr_hi,
        top = ctx.attr.top,
        verilog = ctx.file.verilog,
    )

mrav_bus_device = rule(
    implementation = _mrav_bus_device_impl,
    attrs = {
        "device_id": attr.string(
            mandatory = True,
        ),
        "addr_lo": attr.int(
            doc = "Low address for the address range covered by this device.",
            mandatory = True,
        ),
        "addr_hi": attr.int(
            doc = "High address for the address range covered by this device.",
            mandatory = True,
        ),
        "top": attr.string(
            mandatory = True,
        ),
        "verilog": attr.label(
            mandatory = True,
            allow_single_file = [".sv"],
        ),
    },
)

def _mrav_bus_rtl_impl(ctx):
    dev_descriptor = ctx.actions.declare_file("%s_descriptor" % ctx.attr.name)
    devs = ctx.attr.devices

    # TODO: ensure there is no address overlap

    devices = []
    for dev in devs:
        dev_info = dev[BusDeviceInfo]
        devices.append({
            "device_id": dev_info.device_id,
            "addr_lo": dev_info.addr_lo,
            "addr_hi": dev_info.addr_hi,
        })

    ctx.actions.write(dev_descriptor, json.encode({"descriptor": devices}))

    output_file = ctx.outputs.out

    ctx.actions.run(
        outputs = [output_file],
        inputs = [dev_descriptor],
        executable = ctx.executable.codegen,
        arguments = ["--devices_file", dev_descriptor.path, "--rtl_file", output_file.path],
        progress_message = "Running bus codegen",
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

mrav_bus_rtl = rule(
    implementation = _mrav_bus_rtl_impl,
    attrs = {
        "devices": attr.label_list(
            mandatory = True,
            providers = [BusDeviceInfo],
        ),
        "out": attr.output(
            doc = "Output label for the bus RTL code",
            mandatory = True,
        ),
        "codegen": attr.label(
            default = Label("//hardware/rtl/mravbus/generator/codegen"),
            allow_files = True,
            executable = True,
            cfg = "exec",
        ),
    },
)

BusDeviceStitchInfo = provider(
    doc = "Provider for a Mrav bus device stitching data",
    fields = {
        "device_info": "Device ID",
        "stitch_rtl": "Stitching code",
        "bus_connections": "'Instructions' for stitching the bus module with the generated lines.",
    },
)

def _mrav_bus_stitch_rtl_impl(ctx):
    dev = ctx.attr.device
    dev_info = dev[BusDeviceInfo]

    additional_connections_rtl = "\n".join(["  .{port}({net}),".format(port = key, net = val) for (key, val) in ctx.attr.additional_connections.items()])

    if len(additional_connections_rtl) > 0:
        additional_connections_rtl = additional_connections_rtl[:-1]  # Remove the trailing comma

    maybe_comma = "," if len(ctx.attr.additional_connections) > 0 else ""

    stitch_rtl = """
  // Stitching wires for '{dev_id}'
  logic dev_{dev_id}_read;
  logic dev_{dev_id}_write;
  logic dev_{dev_id}_read_done;
  logic dev_{dev_id}_write_done;
  logic[MRAV_ADDR_WIDTH-1:0] dev_{dev_id}_addr;
  logic[MRAV_DATA_WIDTH-1:0] dev_{dev_id}_cpu_data_out;
  logic[MRAV_DATA_WIDTH-1:0] dev_{dev_id}_cpu_data_in;

  {device_module} {device_module}_i (
    .read(dev_{dev_id}_read),
    .write(dev_{dev_id}_write),
    .read_done(dev_{dev_id}_read_done),
    .write_done(dev_{dev_id}_write_done),
    .addr(dev_{dev_id}_addr),
    .cpu_data_out(dev_{dev_id}_cpu_data_out),
    .cpu_data_in(dev_{dev_id}_cpu_data_in){maybe_comma}
    {additional_connections}
  );
""".format(dev_id = dev_info.device_id, device_module = dev_info.top, maybe_comma = maybe_comma, additional_connections = additional_connections_rtl)

    # Maps bus module's ports to the stitching connections.
    # At the moment it basically stitches with key/value pairs where key matches the value, but doing it like this
    # for future flexibility.
    bus_connections = {
        "dev_{dev_id}_read".format(dev_id = dev_info.device_id): "dev_{dev_id}_read".format(dev_id = dev_info.device_id),
        "dev_{dev_id}_write".format(dev_id = dev_info.device_id): "dev_{dev_id}_write".format(dev_id = dev_info.device_id),
        "dev_{dev_id}_read_done".format(dev_id = dev_info.device_id): "dev_{dev_id}_read_done".format(dev_id = dev_info.device_id),
        "dev_{dev_id}_write_done".format(dev_id = dev_info.device_id): "dev_{dev_id}_write_done".format(dev_id = dev_info.device_id),
        "dev_{dev_id}_addr".format(dev_id = dev_info.device_id): "dev_{dev_id}_addr".format(dev_id = dev_info.device_id),
        "dev_{dev_id}_cpu_data_out".format(dev_id = dev_info.device_id): "dev_{dev_id}_cpu_data_out".format(dev_id = dev_info.device_id),
        "dev_{dev_id}_cpu_data_in".format(dev_id = dev_info.device_id): "dev_{dev_id}_cpu_data_in".format(dev_id = dev_info.device_id),
    }

    output_file = ctx.outputs.out
    ctx.actions.write(output_file, stitch_rtl)

    return [
        DefaultInfo(files = depset([output_file])),
        BusDeviceStitchInfo(device_info = dev_info, stitch_rtl = stitch_rtl, bus_connections = bus_connections),
    ]

mrav_bus_stitch_rtl = rule(
    implementation = _mrav_bus_stitch_rtl_impl,
    attrs = {
        "device": attr.label(
            mandatory = True,
            providers = [BusDeviceInfo],
        ),
        "out": attr.output(
            doc = "Output label for the stitch RTL code",
            mandatory = True,
        ),
        "additional_connections": attr.string_dict(),
    },
)
