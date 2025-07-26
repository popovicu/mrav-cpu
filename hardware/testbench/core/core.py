import dataclasses

from core.proto import core_pb2

@dataclasses.dataclass
class CoreSnapshot:
    pc: int
    r: list[int]
    
    def debug_string(self, regs_query = {}):
        reg_dump = ' '.join(f"{idx}:{r:04X}" for idx, r in enumerate(self.r) if idx in regs_query)
        output = f"(PC = {self.pc:04X}, r = [{reg_dump}])"
        return output
    

def make_snapshot_from_dut(dut):
    pc = int(dut.pc_q.value)
    r = [0] * 16

    for i in range(16):
        r[i] = int(dut.r_q[i].value)

    return CoreSnapshot(pc, r)


def make_snapshot_from_proto_file(file_path):
    with open(file_path, 'rb') as f:
        core_proto = core_pb2.CoreState()
        binary_data = f.read()
        core_proto.ParseFromString(binary_data)
        return CoreSnapshot(core_proto.pc, core_proto.r)