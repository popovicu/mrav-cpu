import os
import pathlib
import pytest
import sys

import cocotb
from cocotb import clock, triggers

from hardware.testbench.components import mrav_bus_memory
from hardware.testbench.core import core
from hardware.testbench.simulation import simulation


@cocotb.test()
async def core_tb(dut):
    with open(os.getenv('SOFTWARE_PATH'), 'rb') as f:
        software_payload = list(f.read())

    simulated_core = core.make_snapshot_from_proto_file(os.getenv('SOFTWARE_CPU_PROTO'))
    memory = mrav_bus_memory.make_memory(dut, 1024, software_payload)
    cocotb.start_soon(memory.work())

    clk = clock.Clock(dut.clk, 10)
    cocotb.start_soon(clk.start(start_high=False))

    await triggers.RisingEdge(dut.clk)
    await triggers.FallingEdge(dut.clk)
    dut.rst_n.value = 0
    await triggers.RisingEdge(dut.clk)
    await triggers.FallingEdge(dut.clk)
    dut.rst_n.value = 1

    for _ in range(40):
        await triggers.RisingEdge(dut.clk)
        await triggers.FallingEdge(dut.clk)
        await triggers.ReadOnly()
        snapshot = core.make_snapshot_from_dut(dut)
    
    assert simulated_core == snapshot


def test_equivalence():
    sim_runner, build_args, test_args = simulation.make_cocotb_runner(
        [os.getenv('CORE_VERILOG')],
        'mrav_core',
        'functions_test',
        {
            "SOFTWARE_PATH": pathlib.Path(os.getenv('TEST_SOFTWARE')).absolute(),
            "SOFTWARE_CPU_PROTO": pathlib.Path(os.getenv('SOFTWARE_CPU_STATE')).absolute(),
        },
    )
    sim_runner.build(**build_args)
    sim_runner.test(**test_args)

if __name__ == "__main__":
    sys.exit(pytest.main(['-v', '--tb=short', '-s'] + sys.argv[1:]))