import os
import pathlib

from cocotb import runner


def make_cocotb_runner(sources, dut_top_level, testbench_python_module, cocotb_env_vars = {}):
    sim = os.getenv("SIM", "verilator")

    # Configure wave dumping based on simulator
    if sim == "icarus":
        extra_args = ["-fst"]  # FST wave format
    elif sim == "verilator":
        extra_args = ["--trace"]  # VCD format

    cocotb_sources = [runner.Verilog(source) for source in sources]

    sim_runner_build_args = {
        'sources': cocotb_sources,
        'hdl_toplevel': dut_top_level,
        'build_args': extra_args,
        'waves': True,
    }

    custom_dir = os.getenv('CUSTOM_BUILD_DIR')

    if custom_dir:
        sim_build_dir = pathlib.Path(custom_dir).absolute()
        print(f'Using custom build dir: {sim_build_dir}')
        sim_build_dir.mkdir(exist_ok=True)
        sim_runner_build_args['build_dir'] = custom_dir

    sim_runner_test_args = {
        'hdl_toplevel': dut_top_level,
        'test_module': testbench_python_module,
        'waves': True,
        'extra_env': cocotb_env_vars,
    }

    return runner.get_runner(sim), sim_runner_build_args, sim_runner_test_args