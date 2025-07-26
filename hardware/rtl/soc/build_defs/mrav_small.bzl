load("//hardware/rtl/mravbus/build_defs:bus.bzl", "mrav_bus_device", "mrav_bus_rtl", "mrav_bus_stitch_rtl")
load("//hardware/rtl/mravbus/build_defs:software.bzl", "mrav_bus_addr_module")
load("//hardware/rtl/mravbus/components/memory/build_defs:generation.bzl", "mrav_imaged_memory")
load("//hardware/rtl/soc/build_defs:top.bzl", "mrav_top")

def mrav_small(name, software, soc_top, gpio_verilog):
    memory_label = "%s_ram_programmed" % name
    ram_programmed_rtl = "%s_ram_programmed.sv" % name

    mrav_imaged_memory(
        name = memory_label,
        memory_size = 64,
        image = software,
        top_module = "memory",
        out = ram_programmed_rtl,
        internal_address_generation = "5:0",
    )

    mem_device_label = "%s_mem" % name

    mrav_bus_device(
        name = mem_device_label,
        device_id = "mem",
        addr_lo = 0x0000,
        addr_hi = 0x005F,
        top = "memory",
        verilog = ram_programmed_rtl,
    )

    gpio_device_label = "%s_gpio" % name

    mrav_bus_device(
        name = gpio_device_label,
        device_id = "gpio",
        addr_lo = 0x0070,
        addr_hi = 0x0070,
        top = "gpio",
        verilog = gpio_verilog,
    )

    gpio_lib_label = "%s_lib" % gpio_device_label
    gpio_lib_module = "%s_gpio_lib.mrav" % name

    mrav_bus_addr_module(
        name = gpio_lib_label,
        bus_device = gpio_device_label,
        symbol_prefix = "GPIO",
        out = gpio_lib_module,
    )

    simple_bus_label = "%s_simple_bus" % name
    simple_bus_rtl = "%s_simple_bus.sv" % name

    mrav_bus_rtl(
        name = simple_bus_label,
        devices = [
            gpio_device_label,
            mem_device_label,
        ],
        out = simple_bus_rtl,
    )

    mem_stitch_label = "%s_mem_stitch" % name
    mem_stitch_rtl = "%s_mem_stitch.sv" % name

    mrav_bus_stitch_rtl(
        name = mem_stitch_label,
        device = mem_device_label,
        out = mem_stitch_rtl,
        additional_connections = {
            "clk": "clk",
        },
    )

    gpio_stitch_label = "%s_gpio_stitch" % name
    gpio_stitch_rtl = "%s_gpio_stitch.sv" % name

    mrav_bus_stitch_rtl(
        name = gpio_stitch_label,
        device = gpio_device_label,
        out = gpio_stitch_rtl,
        additional_connections = {
            "clk": "clk",
            "rst_n": "rst_n",
            "external_output": "gpio",
        },
    )

    bundle_rtl = "%s_bundle.sv" % name

    mrav_top(
        name = name,
        bus_rtl = simple_bus_rtl,
        stitched_devices = [
            gpio_stitch_label,
            mem_stitch_label,
        ],
        core_rtl_bundle = "//hardware/rtl:mrav_core.sv",
        top_rtl = soc_top,
        out = bundle_rtl,
    )
