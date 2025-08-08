# MxsxllBox VM

## Overview

MxsxllBox is a custom virtual machine (VM) designed with a 64 KB memory space divided into several segments, supporting a simple instruction set and features like labels, branching, function calls, and dynamic memory management.\
It has 32 Registers 18 general purpose, 6 Registers for syscalls, 6 for the Scheduler and 2 for an interrupt

## How to run

### run precompiled

- run Debugger in cmd `Debugger-main.exe`
- run normal in cmd `VM-main.exe`

### compile & run your own program
- all stdlib/scheduler files are compiled automatically
- all binaries are in VM-bin


Compile & Run
- edit the program.asm


- build debugger `go build ./Debugger-main`
- run debugger in cmd `Debugger-main.exe`
-  build main `go build ./VM-main`
- run main in cmd `VM-main.exe`

Run
- edit the program.asm


- run debugger directly `go run ./Debugger-main`
- run main directly ` go run ./VM-main`

## Memory Layout

## Memory Layout (64 KB)

| Segment             | Size  | Address Range       | Description                                |
|---------------------|-------|---------------------|--------------------------------------------|
| **Program**         | 8 KB  | `0x0000` – `0x1FFF` | Code and instructions (user + stdlib)      |
| ├─ User Code        | 6 KB  | `0x0000` – `0x17FF` | User program code                          |
| └─ Std. Library     | 2 KB  | `0x1800` – `0x1FFF` | Standard library functions                 |
| **Heap**            | 16 KB | `0x2000` – `0x5FFF` |                                            |
| ├─ Tasks            | 540 B | `0x2381` - `0x259C` | Tasks for Scheduling                       |
| ├─ Bitmap           | 895 B | `0x2000` - `0x237F` | Bitmap with 16 B blocks                    |
| └─ Writeable Heap   | 14 KB | `0x259C` - `0x5D9C` | Dynamic memory allocation (heap)           |
| **Interrupt-Table** | 611 B | `0x5D9D` - `0x....` | Interrupt first jump to here               |
| **Stack**           | 8 KB  | `0x6000` – `0x7FFF` | Stack for function calls, grows downward   |
| **Video RAM**       | 16 KB | `0x8000` – `0xBFFF` | Framebuffer for visual output              |
| **Reserved**        | 8 KB  | `0xC000` – `0xDFFF` | Reserved for I/O, buffers, MMIO            |
| └─ Keyboard I/O     | ~30 B | `0xC000` – `0xC020` | Ring buffer, read/write pointers           |
| **Extra / Future**  | 8 KB  | `0xE000` – `0xFFFF` | Expansion, paging tables, filesystem, etc. |

---

## Instructions and Parsing

- Instructions vary in size, mostly 3 or 5 bytes
- Labels support for jumps and calls, with addresses resolved in two passes
- Opcodes for arithmetic, load/store, jumps, calls, push/pop, print, and halt
- Support for string-related instructions and memory operations 

### Instruction Set

| Byte-code | Instruction | description                                                         |
|-----------|-------------|---------------------------------------------------------------------|
| 0x00      | NOP         | Nothing                                                             |
| 0x01      | LOADB       | `LOADB Rx(return) Ry(addr)`: loads 1 B                              |
| 0x02      | LOADW       | `LOADW Rx(return) Ry(addr)`:  loads a word 2 B                      |
| 0x03      | STOREB      | `STOREB Rx(val) Ry(addr)`: Stores 1 B to an addr                    |
| 0x04      | STOREW      | `STOREW Rx(val) Ry(addr)`: Stores a word 2B                         |
| 0x05      | ADD         | `ADD Rx(return, val1) Ry(val2)`                                     |
| 0x06      | SUB         | `SUB Rx(return, val1) Ry(val2)`                                     |
| 0x07      | MUL         | `MUL Rx(return, val1) Ry(val2)`                                     |
| 0x08      | DIV         | `DIV Rx(return, val1) Ry(val2)`                                     |
| 0x09      | JMP         | `JMP Lbl`: jumps to a label                                         |
| 0x0A      | JZ          | `JZ lbl`: jumps if 0-flag is triggered                              |
| 0x0B      | JC          | `JC lbl`: jumps if carry-flag is triggered                          |
| 0x0C      | PRINT       | `PRINT Rx`: Prints the val of a Register                            |
| 0x0D      | MOVI        | `MOVI Rx imm`: Load immediate value to Reg                          |
| 0x0E      | ADDI        | `ADDI Rx imm`: Adds immediate value to Reg                          |
| 0x0F      | SUBI        | `SUBI Rx imm`: Subs immediate value from Reg                        |
| 0x10      | MULI        | `MULI Rx imm`: Multiplies Reg with immediate                        |
| 0x11      | DIVI        | `DIVI Rx imm`: Divides Reg with immediate                           |
| 0x12      | LOAD        | `LOAD Rx Ry`: Loads either byte or word depending on loacation      |
| 0x13      | STORE       | `STORE Rx Ry`: Stores --""--                                        |
| 0x14      | PUSH        | `PUSH Rx`: Pushes Reg(val) to the Stack                             |
| 0x15      | POP         | `POP Rx`: Pops(removes) val into a Reg                              |
| 0x16      | CALL        | `CALL Lbl`: Jmps to a Lbl and pushes the current PC to the Stack    |
| 0x17      | RET         | `RET`: Returns to the PC on the top of the stack                    |
| 0x18      | ALLOC       | <del>`ALLOC`</del>Deprecated                                        |
| 0x19      | FREE        | <del>`FREE`</del>Deprecated                                         |
| 0x1A      | PRINTSTR    | `PRINTSTR Rx`: prints a string starting at addr Rx                  |
| 0x1B      | JNZ         | `JNZ Lbl`: Jumps to a Lbl if the 0-flag isn't set                   |
| 0x1C      | JNC         | `JNC Lbl`: Jumps to a Lbl if the carry-flag isn't set               |
| 0x1D      | CMP         | `CMP Rx Ry`: Compares regs. Sets 0-flag if equal carry if bigger |
| 0x1E      | CMPI        | `CMPI Rx imm`: Compares reg & imm --""--                            |
| 0x1F      | TEST        | `TEST Rx Ry`: Sees if any bytes overlap                             |
| 0x20      | TSTI        | `TEST Rx imm`: --""--                                               |
| 0x21      | JL          | `JL Lbl`: Jump if less(!0-flag & !C-flag)                           |
| 0x22      | JLE         | `JLE Lbl`: Jump if less equal(!C-flag)                              |
| 0x23      | JG          | `JG Lbl`: Jump if greater(!0-flag & C-flag)                         |
| 0x24      | JGE         | `JGE Lbl`: Jump if greater equal(C-flag)                            |
| 0x25      | STZ         | `STZ`: Sets 0-flag                                                  |
| 0x26      | STC         | `STC`: Sets carry-flag                                              |
| 0x27      | CLZ         | `CLZ`: Clears 0-flag                                                |
| 0x28      | CLC         | `CLC`: Clears carry-flag                                            |
| 0x29      | MOD         | `MOD Rx Ry`: Rx = Rx % Ry                                           |
| 0x2A      | MOV         | `MOV Rx Ry`: Moves value from Ry to Rx                              |
| 0x2B      | MODI        | `MODI Rx imm`: Rx = Rx % imm                                        |
| 0x2C      | RS          | `RS Rx Ry`: Rightshift Rx << Ry                                     |
| 0x2D      | LS          | `LS Rx Ry`: Leftshift Rx << Ry                                      |
| 0x2E      | OR          | `OR Rx Ry`: Bitwise or Rx \| Ry                                     |
| 0x2F      | AND         | `AND Rx Ry`: Bitwise and Rx & Ry                                    |
| 0x30      | MOVA        | `MOVA Rx Lbl`: Move Lbladdr into Rx                                 |
| 0x31      | GPC         | `GPC Rx`: Get PC into Rx                                            |
| 0x32      | SPC         | `SPC Rx`: Sets PC to Rx                                             |
| 0x33      | GSP         | `GSP Rx`: Get SP into Rx                                            |
| 0x34      | SSP         | `SSP Rx`: Sets SP to Rx                                             |
| 0x35      | GRFN        | `GRFN Rx Ry`: Get Register from Number Rx = registernum Ry(return)  |
| 0x36      | GF          | `GF Rx`: Gets flags byte 0 = 0-flag byte 1 = carry-flag             |
| 0x37      | SF          | `SF Rx`: Sets flags --""--                                          |
| 0x38      | SRFN        | `SRFN`: Set Register from Number Rx = registernum Ry = val          |
| 0x39      | YIELD       | `YIELD`: Doesn't yield but dissables interrupt will change          |
| 0x3A      | UNYIELD     | `UNYIELD`: Enables all interrupts could be masked                   |
| 0xFF      | HALT        | `HALT`: Stops the Program                                           |


---

## Dynamic Memory Allocation

- Heap size: 16 KB
- writeable Heap size: 14 KB
- Uses a **bitmap allocator** with block size of 16 bytes
- Metadata stored in the first word of an allocation block
- `alloc`: requests block counts (multiples of 16 bytes)
- `free`: returns blocks to the heap
- Bitmap is stored at the beginning of Heap after the tasks


## Std. Lib. 

- Labels that are saved into the library Region of the Program Space
- ProgramStdLibStart = 0x0C01 
- ProgramEnd         = 0x0FFF
- In Syscalls O1 is always the return addr if there's only 1 return value

### String functions

- `_strcpy`: copies a String to O1 from O2
- `_strlen`: Loads len(O2) into O1
- `_strcmp`: compares two strings sets 0 flag if they are equal carry-flag if a byte is higher
- `_strcat`: concacts two strings O1 = O1+O2 `"a", "b" = "ab"`

### Utility funtctions

- `_memset`: sets a region of memory(lenght O3) starting at addr(O1) to val(O2)
- `_memcpy`: copies from addr(O1) to addr(O2) for ammount(O3) bytes 

### sys functions

- `_alloc`: allocates ammount(O2) blocks(16 B) the start is returned in O1
- `_free`: frees a block of memory O1 is the start of that Memory 


### math functions

- basic arithmetic:
  - `_add`/`_sub`/`_mul`/`_div`/`_mod`/`_inc`/`_dec`
- `_max`: returns the larger val of O1/O2 in O1
- `_min`: return the smaller val of O1/O2 in O1
- `_pow`: returns O1**O2 in O1
- `_clamp`: if Val(O2) is between two values(low O1, high O3) the 0 flag is set

### io functions

- `_printstr`: prints a string from O2 is also in instruction set
- `_printchar`: prints a char from O2 O1 contains char
- `_readchar`: reads a char from the Keyboard Buffer into O1

---

## Keyboard Buffer

- The Keyboard-buffer is a [ring buffer]("https://en.wikipedia.org/wiki/Circular_buffer") N = 30
- It's write ptr is at 0x0C000 the read prt is at 0x0C001

## Custom Binary

- A custom Binary set for the linker 
- Contains code, Symbol, Relocation- lenghts under "MXOB" header
- Globals are identified as labels with an underscore

## Scheduler

- 9 Task slots
- `_spawn`: creates task at addr(O1)
- `_yield`: willingly gives up control
- `_init_scheduler`: gives control to the scheduler only used once at the beginning

> To avoid a single  action like redrawing the screen | reading the keyboard-buffer "hogging" all resources \
> the scheduler can save the current context/state of the vm meaning(regs, PC, SP and flags) to memory.\
> To then give another Task the oportunity to continue. \
> This can occur with the help of `_yield`: which willingly gives up control,\
> a code can be moved into O2(see `Yield Table`) this code confirms if the task is IO,\
> for example Keyboard blocke(waiting on input) or smth else code 1 is ready,\
> meaning that it can be chosen again if no other task is currently available.\
> The other possibility for changing the current task is an `interrupt`.\
> Either IO | hardware timer these are forced and interrupt the program wherever it is at the moment.\
> A task is created by first moving its addr into O1 using `MOVA REG LBL`: and then calling `_spawn`\
> After creating all tasks control can be given to the scheduler by calling `_init_scheduler`.\
> The scheduler will start with the last task added and work it's way down before returning to the last.\

### Yield Table

| Code | description                                             |
|------|---------------------------------------------------------|
| 0    | `running`(only internal)                                |
| 1    | `ready`(is ready)                                       |
| 2    | `keyboard-blocked`: (is set to ready after input)       |
| 3    | `timer/blocked`: (waiting on the next timer interrupt)  |
| 4    | `unused`:                                               |
| 5    | `unused`:                                               |
| 6    | `unused`:                                               |
| 7    | `terminated`: (the tasks isn't supposed to run anymore) |

## Debugger

- opens a window where a decompiled version of the script is shown
- the currently active line is highlighted
- Step mode:
  - progress forward with right arrow key
- Run mode:
  - Runs until mode is switched back to step, or it hits a breakpoint
- Breakpoints can be set by left-clicking on any line
- When jumping to lbls the debugger jmps to the first line with actual content
- Set a break point there and not on the lbl name
- Any Lbls that aren't jmped, called etc. to can't be decompiled
- Regs, Pc and Stack-Top are visible on the right side of the screen

## Design Decisions

- Memory layout optimized for simplicity and performance
- Labels aren't assembled into byte code 
- The pre-Assembled code is outputted in a .obj file
- The .obj file is then linked together with other .obj files into a  final .bin file
- _The linker is given a `map[string]int`: that contains [.obj filepaths] location
- Heap allocator uses first-fit policy for simplicity and speed

---

## Current Status

- VM running with working jumps, calls, arithmetic
- RET instruction fixed to correctly return from calls
- String Support lenght based indexing
- String functions like `_strcpy`
- Basic Memory functions like `_memset`
- Bitmap allocator
- Custom Binary Format 
- Keyboard input handling via circular buffer
- Planning to add:
    - Visual Output
    - File system
    - custom mini tcp/ip stack
