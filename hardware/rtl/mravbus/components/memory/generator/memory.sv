module {{.TopModule}} (
    input logic clk,

    input logic read,
    input logic write,

    output logic read_done,
    output logic write_done,

    input logic[MRAV_ADDR_WIDTH-1:0] addr,
    input logic[MRAV_DATA_WIDTH-1:0] cpu_data_out,
    output logic[MRAV_DATA_WIDTH-1:0] cpu_data_in
);
    logic[{{.InternalAddressRtl}}] first_byte_internal, second_byte_internal;
    assign first_byte_internal = addr[{{.InternalAddressRtl}}];
    assign second_byte_internal = first_byte_internal + 1;

    logic[7:0] mem [{{.MemSize}}];

    logic[MRAV_ADDR_WIDTH-1:0] second_byte_addr;
    logic first_byte_valid, second_byte_valid;

    assign first_byte_valid = first_byte_internal <= {{sub .MemSize 1}};
    assign second_byte_valid = second_byte_internal <= {{sub .MemSize 1}};

    logic[7:0] first_byte, second_byte;

    assign first_byte = (first_byte_valid) ? mem[first_byte_internal] : 0;
    assign second_byte = (second_byte_valid) ? mem[second_byte_internal] : 0;
    
    initial begin
{{- range $index, $dataByte := .Payload}}
        mem[{{$index}}] = 8'{{$dataByte}};
{{- end}}
    end

    always_ff @(posedge clk) begin
        if (write) begin
            if (first_byte_valid) begin
                mem[first_byte_internal] <= cpu_data_out[15:8];
            end

            if (second_byte_valid) begin
                mem[second_byte_internal] <= cpu_data_out[7:0];
            end
        end
    end

    assign cpu_data_in = { first_byte, second_byte };

    assign read_done = read; // It's always a one clock operation for this module.
    assign write_done = write;

endmodule