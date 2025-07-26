module mrav_soc(
    input logic clk,
    input logic rst_n,

    output logic[7:0] gpio
);
    {stitchings}

    logic core_read;
    logic core_write;

    logic mrav_core_read_done;
    logic mrav_core_write_done;

    logic[MRAV_ADDR_WIDTH-1:0] mrav_core_addr;
    logic[MRAV_DATA_WIDTH-1:0] mrav_core_data_out;
    logic[MRAV_DATA_WIDTH-1:0] mrav_core_data_in;

    mrav_bus bus_i(
        // CPU connections first
        .core_read(core_read),
        .core_write(core_write),
        .mrav_addr(mrav_core_addr),
        .mrav_data_out(mrav_core_data_out),
        .mrav_data_in(mrav_core_data_in),
        .mrav_read_done(mrav_core_read_done),
        .mrav_write_done(mrav_core_write_done){bus_connections}
    );

    mrav_core core_i(
        .clk(clk),
        .rst_n(rst_n),
        .read(core_read),
        .write(core_write),
        .read_done(mrav_core_read_done),
        .write_done(mrav_core_write_done),
        .addr(mrav_core_addr),
        .data_out(mrav_core_data_out),
        .data_in(mrav_core_data_in)
    );
endmodule