module gpio(
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
    logic [7:0] gpio_q, gpio_d;

    always_ff @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            gpio_q <= 8'h00;
        end else begin
            gpio_q <= gpio_d;
        end
    end

    assign external_output = gpio_q;

    // TODO: maybe don't depend on the bus only to generate correct 'read' and 'write', and check the addr too.

    assign gpio_d = (write) ? cpu_data_out[7:0] : gpio_q;
    assign cpu_data_in = { 8'h00, gpio_q };
    assign read_done = read; // It's always a one clock operation for this module.
    assign write_done = write;
endmodule