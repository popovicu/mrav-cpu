module nop(
    input logic clk,
    input logic rst_n,

    input logic read,
    input logic write,

    output logic read_done,
    output logic write_done,

    input logic[MRAV_ADDR_WIDTH-1:0] addr,
    input logic[MRAV_DATA_WIDTH-1:0] cpu_data_out,
    output logic[MRAV_DATA_WIDTH-1:0] cpu_data_in,

    output logic[7:0] external_output
);
    // Nop
endmodule