module mrav_bus(
    input logic core_read,
    input logic core_write,

    input logic[MRAV_ADDR_WIDTH-1:0] mrav_addr,
    input logic[MRAV_DATA_WIDTH-1:0] mrav_data_out,
    output logic[MRAV_DATA_WIDTH-1:0] mrav_data_in,
    output logic mrav_read_done,
    output logic mrav_write_done,

{{- range $index, $peripheral := .Peripherals}}
    output logic dev_{{$peripheral.DeviceId}}_read,
    output logic dev_{{$peripheral.DeviceId}}_write,
    input logic dev_{{$peripheral.DeviceId}}_read_done,
    input logic dev_{{$peripheral.DeviceId}}_write_done,
    output logic[MRAV_ADDR_WIDTH-1:0] dev_{{$peripheral.DeviceId}}_addr,
    output logic[MRAV_DATA_WIDTH-1:0] dev_{{$peripheral.DeviceId}}_cpu_data_out,
    input logic[MRAV_DATA_WIDTH-1:0] dev_{{$peripheral.DeviceId}}_cpu_data_in{{if not (isLast $index $.Peripherals)}},{{end}}
{{- end}}
);

{{- range $index, $peripheral := .Peripherals}}
    logic dev_{{$peripheral.DeviceId}}_hit;
    // Disabling the unsigned lint because unsigned semantics are indeed needed here and codegen can generate a redundant check like >= 0.
    /* verilator lint_off UNSIGNED */
    assign dev_{{$peripheral.DeviceId}}_hit = (mrav_addr >= {{$peripheral.AddrLo}}) && (mrav_addr <= {{$peripheral.AddrHi}});
    /* verilator lint_off UNSIGNED */
    assign dev_{{$peripheral.DeviceId}}_read = dev_{{$peripheral.DeviceId}}_hit && core_read;
    assign dev_{{$peripheral.DeviceId}}_write = dev_{{$peripheral.DeviceId}}_hit && core_write;
    assign dev_{{$peripheral.DeviceId}}_addr = mrav_addr;
    assign dev_{{$peripheral.DeviceId}}_cpu_data_out = mrav_data_out;
{{- end}}

    always_comb begin
        mrav_data_in = 0;
        mrav_read_done = 0;
        mrav_write_done = 0;

{{- range $index, $peripheral := .Peripherals}}
        {{ if eq $index 0 }}if{{ else }}else if{{ end }} (dev_{{$peripheral.DeviceId}}_hit) begin mrav_data_in = dev_{{$peripheral.DeviceId}}_cpu_data_in; mrav_read_done = dev_{{$peripheral.DeviceId}}_read_done; mrav_write_done = dev_{{$peripheral.DeviceId}}_write_done; end
{{- end}}
    end

endmodule