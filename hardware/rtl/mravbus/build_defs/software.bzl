load("//software/build_defs:mrav.bzl", "MravLibInfo")
load(":bus.bzl", "BusDeviceInfo")

def _mrav_bus_addr_module_impl(ctx):
    dev = ctx.attr.bus_device[BusDeviceInfo]
    addr_lo = dev.addr_lo
    addr_hi = dev.addr_hi

    if addr_lo != addr_hi:
        fail("currently, only single register devices are supported")

    # mrav_module = "{prefix}_LO = {lo}\n{prefix}_HI = {hi}\n".format(prefix = ctx.attr.symbol_prefix, lo = "foo", hi = "bar")
    mrav_module = "{prefix}_ADDR = {addr}\n".format(prefix = ctx.attr.symbol_prefix, addr = addr_lo)

    ctx.actions.write(ctx.outputs.out, mrav_module)

    return [
        MravLibInfo(modules = [ctx.outputs.out], deps = depset()),
        DefaultInfo(files = depset([ctx.outputs.out])),
    ]

mrav_bus_addr_module = rule(
    implementation = _mrav_bus_addr_module_impl,
    attrs = {
        "bus_device": attr.label(
            mandatory = True,
            providers = [BusDeviceInfo],
        ),
        "symbol_prefix": attr.string(
            mandatory = True,
        ),
        "out": attr.output(
            mandatory = True,
        ),
    },
    provides = [MravLibInfo],
)
