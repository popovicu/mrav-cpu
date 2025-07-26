from cocotb import triggers


class Memory:
    def __init__(self, dut, memory_bytes):
        self.dut = dut
        self.active = True
        self._memory = memory_bytes
    
    async def work(self):
        while self.active:
            await triggers.RisingEdge(self.dut.clk)

            await triggers.ReadOnly()
            mem_addr = int(self.dut.addr.value)
            await triggers.ReadWrite()

            self.dut.read_done.value = 0
            self.dut.write_done.value = 0

            if (self.dut.read.value == 1) and (self.dut.write.value == 1):
                raise ValueError('Core both reading and writing')
            

            if (mem_addr >= len(self._memory)) or ((mem_addr + 1) >= len(self._memory)):
                raise ValueError(f'Address {mem_addr} out of bounds')
            
            if (self.dut.read.value == 1):
                # print(f'[Memory] Reading: {mem_addr:04X}')
                first_byte = self._memory[mem_addr]
                second_byte = self._memory[mem_addr + 1]
                value = (first_byte << 8) | second_byte
                # print(f'  Value: {value:04X}')
                self.dut.data_in.value = value
                self.dut.read_done.value = 1
                continue

            if (self.dut.write.value == 1):
                # print(f'[Memory] Writing: {mem_addr:04X}')
                core_data = int(self.dut.data_out.value)
                # print(f'  [Memory] Value: {core_data:04X}')
                first_byte = (core_data >> 8) & 0xFF
                second_byte = core_data & 0xFF
                self._memory[mem_addr] = first_byte
                self._memory[mem_addr + 1] = second_byte
                self.dut.write_done.value = 1
                continue
            

def make_memory(dut, size, image):
    if len(image) > size:
        raise ValueError(f'Unable to create a memory of size {size} because the given image is of length {len(image)}')
    
    memory = [0] * size
    memory[0:len(image)] = image

    return Memory(dut, memory)