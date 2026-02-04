## Instructions and Parsing

### Important --- ReadMe

- Due to undecisiveness STORE instructions differ from the rest
- While all other instruction have the receiving operand left (MOV R0 R1; moves the val of R1 into R0)
- The STORE instruction stores Register RX into Ry / an immediate addr
- 

### Instruction Set

| Byte-code | Instruction | description                                                        |
|-----------|-------------|--------------------------------------------------------------------|
| 0x00      | NOP         | Nothing                                                            |
| 0x01      | LOADB       | `LOADB Rx(return) Ry(addr)`: loads 1B into Rx                      |
| 0x02      | LOADW       | `LOADW Rx(return) Ry(addr)`:  loads a word into Rx                 |
| 0x03      | STOREB      | `STOREB Rx(val) Ry(addr)`: Stores 1 B to  Ry/immediate             |
| 0x04      | STOREW      | `STOREW Rx(val) Ry(addr)`: Stores a word to Ry/immediate           |
| 0x05      | ADD         | `ADD Rx(return, val1) Ry(val2)`                                    |
| 0x06      | SUB         | `SUB Rx(return, val1) Ry(val2)`                                    |
| 0x07      | MUL         | `MUL Rx(return, val1) Ry(val2)`                                    |
| 0x08      | DIV         | `DIV Rx(return, val1) Ry(val2)`                                    |
| 0x09      | JMP         | `JMP Lbl`: jumps to a label                                        |
| 0x0A      | JZ          | `JZ lbl`: jumps if 0-flag is triggered                             |
| 0x0B      | JC          | `JC lbl`: jumps if carry-flag is triggered                         |
| 0x0C      | PRINT       | `PRINT Rx`: Prints the val of a Register                           |
| 0x0D      | MOVI        | `MOVI Rx imm`: Load immediate value to Reg                         |
| 0x0E      | ADDI        | `ADDI Rx imm`: Adds immediate value to Reg                         |
| 0x0F      | SUBI        | `SUBI Rx imm`: Subs immediate value from Reg                       |
| 0x10      | MULI        | `MULI Rx imm`: Multiplies Reg with immediate                       |
| 0x11      | DIVI        | `DIVI Rx imm`: Divides Reg with immediate                          |
| 0x12      | LOAD        | `LOAD Rx Ry`: Loads either byte or word depending on loacation     |
| 0x13      | STORE       | `STORE Rx Ry`: Stores --""--                                       |
| 0x14      | PUSH        | `PUSH Rx`: Pushes Reg(val) to the Stack                            |
| 0x15      | POP         | `POP Rx`: Pops(removes) val into a Reg                             |
| 0x16      | CALL        | `CALL Lbl`: Jmps to a Lbl and pushes the current PC to the Stack   |
| 0x17      | RET         | `RET`: Returns to the PC on the top of the stack                   |
| 0x18      | ALLOC       | <del>`ALLOC`</del>Deprecated                                       |
| 0x19      | FREE        | <del>`FREE`</del>Deprecated                                        |
| 0x1A      | PRINTSTR    | `PRINTSTR Rx`: prints a string starting at addr Rx                 |
| 0x1B      | JNZ         | `JNZ Lbl`: Jumps to a Lbl if the 0-flag isn't set                  |
| 0x1C      | JNC         | `JNC Lbl`: Jumps to a Lbl if the carry-flag isn't set              |
| 0x1D      | CMP         | `CMP Rx Ry`: Compares regs. Sets 0-flag if equal carry if bigger   |
| 0x1E      | CMPI        | `CMPI Rx imm`: Compares reg & imm --""--                           |
| 0x1F      | TEST        | `TEST Rx Ry`: Sees if any bytes overlap                            |
| 0x20      | TSTI        | `TEST Rx imm`: --""--                                              |
| 0x21      | JL          | `JL Lbl`: Jump if less(!0-flag & !C-flag)                          |
| 0x22      | JLE         | `JLE Lbl`: Jump if less equal(!C-flag)                             |
| 0x23      | JG          | `JG Lbl`: Jump if greater(!0-flag & C-flag)                        |
| 0x24      | JGE         | `JGE Lbl`: Jump if greater equal(C-flag)                           |
| 0x25      | STZ         | `STZ`: Sets 0-flag                                                 |
| 0x26      | STC         | `STC`: Sets carry-flag                                             |
| 0x27      | CLZ         | `CLZ`: Clears 0-flag                                               |
| 0x28      | CLC         | `CLC`: Clears carry-flag                                           |
| 0x29      | MOD         | `MOD Rx Ry`: Rx = Rx % Ry                                          |
| 0x2A      | MOV         | `MOV Rx Ry`: Moves value from Ry to Rx                             |
| 0x2B      | MODI        | `MODI Rx imm`: Rx = Rx % imm                                       |
| 0x2C      | RS          | `RS Rx Ry`: Rightshift Rx << Ry                                    |
| 0x2D      | LS          | `LS Rx Ry`: Leftshift Rx << Ry                                     |
| 0x2E      | OR          | `OR Rx Ry`: Bitwise or Rx \| Ry                                    |
| 0x2F      | AND         | `AND Rx Ry`: Bitwise and Rx & Ry                                   |
| 0x30      | MOVA        | `MOVA Rx Lbl`: Move Lbladdr into Rx                                |
| 0x31      | GPC         | `GPC Rx`: Get PC into Rx                                           |
| 0x32      | SPC         | `SPC Rx`: Sets PC to Rx                                            |
| 0x33      | GSP         | `GSP Rx`: Get SP into Rx                                           |
| 0x34      | SSP         | `SSP Rx`: Sets SP to Rx                                            |
| 0x35      | GRFN        | `GRFN Rx Ry`: Get Register from Number Rx = registernum Ry(return) |
| 0x36      | GF          | `GF Rx`: Gets flags byte 0 = 0-flag byte 1 = carry-flag            |
| 0x37      | SF          | `SF Rx`: Sets flags --""--                                         |
| 0x38      | SRFN        | `SRFN`: Set Register from Number Rx = registernum Ry = val         |
| 0x39      | YIELD       | `YIELD`: Doesn't yield but dissables interrupt will change         |
| 0x3A      | UNYIELD     | `UNYIELD`: Enables all interrupts could be masked                  |
| 0xFF      | HALT        | `HALT`: Stops the Program                                          |