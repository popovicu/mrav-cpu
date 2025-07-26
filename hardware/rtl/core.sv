typedef enum logic [1:0] {
    CORE_READY    = 2'b00,
    CORE_LW_READ  = 2'b01,
    CORE_SW_WRITE = 2'b10
} state_t;

typedef enum logic [3:0] {
    MRAV_ADD = 4'h0,
    MRAV_SUB = 4'h1,
    MRAV_LW = 4'h2,
    MRAV_SW = 4'h3,
    MRAV_XOR = 4'h4,
    MRAV_AND = 4'h5,
    MRAV_OR = 4'h6,
    MRAV_ADDI = 4'h7,
    MRAV_LDHI = 4'h8,
    MRAV_BZ = 4'h9,
    MRAV_BNZ = 4'hA,
    MRAV_JAL = 4'hB,
    MRAV_JALR = 4'hC,
    MRAV_SHL = 4'hD,
    MRAV_SHR = 4'hE,
    MRAV_SHRA = 4'hF
} instruction_t;

function instruction_t decode_instruction(logic[MRAV_DATA_WIDTH-1:0] instruction_value);
    return instruction_t'(instruction_value[15:12]); // TODO: do not hardcode
endfunction

function logic[3:0] decode_rd(logic[MRAV_DATA_WIDTH-1:0] instruction_value);
    return instruction_value[11:8];
endfunction

function logic[3:0] decode_rs1(logic[MRAV_DATA_WIDTH-1:0] instruction_value);
    return instruction_value[7:4];
endfunction

function logic[3:0] decode_rs2(logic[MRAV_DATA_WIDTH-1:0] instruction_value);
    return instruction_value[3:0];
endfunction

function logic[7:0] decode_imm8(logic[MRAV_DATA_WIDTH-1:0] instruction_value);
    return instruction_value[7:0];
endfunction

function logic[3:0] decode_imm4(logic[MRAV_DATA_WIDTH-1:0] instruction_value);
    return instruction_value[7:4];
endfunction

module mrav_core(
    input logic clk,
    input logic rst_n,

    output logic read,
    output logic write,

    input logic read_done,
    input logic write_done,

    output logic[MRAV_ADDR_WIDTH-1:0] addr,
    output logic[MRAV_DATA_WIDTH-1:0] data_out,
    input logic[MRAV_DATA_WIDTH-1:0] data_in
);
    state_t state_q, state_d;

    logic [15:0] pc_q, pc_d;
    logic [15:0] r_q[MRAV_REG_NUM], r_d[MRAV_REG_NUM];
    logic [15:0] instruction_q, instruction_d;

    assign read = (state_q == CORE_READY) || (state_q == CORE_LW_READ);
    assign write = (state_q == CORE_SW_WRITE);
    assign addr = (state_q == CORE_READY) ? pc_q : ((state_q == CORE_LW_READ) ? r_q[int'(decode_rs1(instruction_q))] : ((state_q == CORE_SW_WRITE) ? r_q[int'(decode_rd(instruction_q))] : 16'hxxxx));

    always_ff @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state_q <= CORE_READY;
            pc_q <= 16'h0000;
            instruction_q <= 16'h0000;

            for (int i = 0; i < MRAV_REG_NUM; i++) begin
                r_q[i] <= 16'h0000;
            end
        end else begin
            state_q <= state_d;
            pc_q <= pc_d;
            instruction_q <= instruction_d;

            for (int i = 0; i < MRAV_REG_NUM; i++) begin
                r_q[i] <= r_d[i];
            end
        end
    end

    instruction_t current_instruction;
    logic[3:0] rd, rs1, rs2;
    logic[7:0] imm8;
    logic[3:0] imm4;

    assign rd = decode_rd(data_in);
    assign rs1 = decode_rs1(data_in);
    assign rs2 = decode_rs2(data_in);
    assign imm8 = decode_imm8(data_in);
    assign imm4 = decode_imm4(data_in); 

    always_comb begin
        pc_d = pc_q;
        instruction_d = instruction_q;
        state_d = state_q;

        for (int i = 0; i < MRAV_REG_NUM; i++) begin
            r_d[i] = r_q[i];
        end

        unique case (state_q)
            CORE_READY: begin
                if (read_done) begin
                    instruction_d = data_in;
                    current_instruction = decode_instruction(data_in);

                    unique case (current_instruction)
                        MRAV_ADD: begin
                            r_d[rd] = r_q[rs1] + r_q[rs2];
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_SUB: begin
                            r_d[rd] = r_q[rs1] - r_q[rs2];
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_LW: begin
                            state_d = CORE_LW_READ;
                        end
                        MRAV_SW: begin
                            state_d = CORE_SW_WRITE;
                        end
                        MRAV_XOR: begin
                            r_d[rd] = r_q[rs1] ^ r_q[rs2];
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_AND: begin
                            r_d[rd] = r_q[rs1] & r_q[rs2];
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_OR: begin
                            r_d[rd] = r_q[rs1] | r_q[rs2];
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_ADDI: begin
                            r_d[rd] = r_q[rd] + {8'b0, imm8};
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_LDHI: begin
                            r_d[int'(rd)] = {imm8, r_q[int'(rd)][7:0]};
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_BZ: begin
                            if (r_q[rd] == 0) begin
                                pc_d = {8'b0, imm8};
                            end else begin
                                pc_d = pc_q + 2;
                            end
                            state_d = CORE_READY;
                        end
                        MRAV_BNZ: begin
                            if (r_q[rd] != 0) begin
                                pc_d = {8'b0, imm8};
                            end else begin
                                pc_d = pc_q + 2;
                            end
                            state_d = CORE_READY;
                        end
                        MRAV_JAL: begin
                            pc_d = {8'b0, imm8};
                            r_d[rd] = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_JALR: begin
                            pc_d = r_q[rs1];
                            r_d[rd] = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_SHL: begin
                            r_d[rd] = r_q[rd] << imm4;
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_SHR: begin
                            r_d[rd] = r_q[rd] >> imm4;
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        MRAV_SHRA: begin
                            r_d[rd] = r_q[rd] >>> imm4;
                            pc_d = pc_q + 2;
                            state_d = CORE_READY;
                        end
                        default: begin
                            state_d = CORE_READY;
                            // TODO: report an error or something
                        end
                    endcase
                end
            end
            CORE_LW_READ: begin
                if (read_done) begin
                    r_d[int'(decode_rd(instruction_q))] = data_in;
                    pc_d = pc_q + 2;
                    state_d = CORE_READY;
                end
            end
            CORE_SW_WRITE: begin
                if (write_done) begin
                    pc_d = pc_q + 2;
                    state_d = CORE_READY;
                end
            end
            default: begin
                // TODO: report an error or something
            end
        endcase
    end

    // Bus data can only come out of the SW instruction.
    assign data_out = (state_q == CORE_SW_WRITE) ? r_q[int'(decode_rs1(instruction_q))] : 16'hxxxx;
endmodule
