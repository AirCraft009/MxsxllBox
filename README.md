# MxsxllBox VM

## Overview

MxsxllBox is a custom virtual machine (VM) designed with a 64 KB memory space divided into several segments, supporting a simple instruction set and features like labels, branching, function calls, and dynamic memory management.\
It has 11 general purpose Registers and 5 Registers for syscalls 

---

## Memory Layout

## 🧠 Memory Layout (64 KB)

| Segment            | Size  | Address Range       | Description                                |
|--------------------|-------|---------------------|--------------------------------------------|
| **Program**        | 8 KB  | `0x0000` – `0x1FFF` | Code and instructions (user + stdlib)      |
| ├─ User Code       | 6 KB  | `0x0000` – `0x17FF` | User program code                          |
| └─ Std. Library    | 2 KB  | `0x1800` – `0x1FFF` | Standard library functions                 |
| **Heap**           | 16 KB | `0x2000` – `0x5FFF` | Dynamic memory allocation (heap)           |
| **Stack**          | 8 KB  | `0x6000` – `0x7FFF` | Stack for function calls, grows downward   |
| **Video RAM**      | 16 KB | `0x8000` – `0xBFFF` | Framebuffer for visual output              |
| **Reserved**       | 8 KB  | `0xC000` – `0xDFFF` | Reserved for I/O, buffers, MMIO            |
| └─ Keyboard I/O    | ~30 B | `0xC000` – `0xC020` | Ring buffer, read/write pointers           |
| **Extra / Future** | 8 KB  | `0xE000` – `0xFFFF` | Expansion, paging tables, filesystem, etc. |

---

## Instructions and Parsing

- Instructions vary in size, mostly 3 or 5 bytes
- Labels support for jumps and calls, with addresses resolved in two passes
- Opcodes for arithmetic, load/store, jumps, calls, push/pop, print, and halt
- Support for string-related instructions and memory operations planned

---

## Dynamic Memory Allocation

- Heap size: 4 KB
- Uses a **bitmap allocator** with block size of 16 bytes
- Metadata stored in the first word of an allocation block
- `alloc` requests block counts (multiples of 16 bytes)
- `free` returns blocks to the heap
- Bitmap is stored outside the heap in VM internal structures


## Std. Lib. 

- Labels that are saved into the library Region of the Program Space
- ProgramStdLibStart = 0x0C01 
- ProgramEnd         = 0x0FFF
- In Syscalls O1 is always the return addr
- When I mention copy from Reg to Reg I always mean the Addr inside of the Reg

### String functions

- `_strcpy` copies a String to O1 from O2
- `_strlen` Loads len(O2) into O1
- `_strcmp` compares two strings sets 0 flag if they are equal C if a byte is higher
- `_strcat` concacts two strings O1 = O1+O2 `"a", "b" = "ab"`

### io functions

- `_printstr` prints a string from O2 is also in instruction set
- `_printchar` prints a char from O2 O1 contains char
- `_readchar` reads a char from the Keyboard Buffer into O1

---

## Keyboard Buffer

- The Keyboard-buffer is a [ring buffer]("https://en.wikipedia.org/wiki/Circular_buffer") N = 30
- It's write ptr is at 0x0C000 the read prt is at 0x0C001

## Custom Binary

- A custom Binary set for the linker 
- Contains code, Symbol, Relocation- lenghts under "MXOB" header
- Globals are identified as labels with an underscore

## Design Decisions

- Memory layout optimized for simplicity and performance
- Labels arent assembled into byte code 
- The pre-Assembled code is outputted in a .obj file
- The .obj file is then linked together with other .obj files into a  final .bin file
- The linker is given a `map[string]int` that contains .obj files 
- It also contains the location of where the code is meant to be stored
- Heap allocator uses first-fit policy for simplicity and speed

---

## Current Status

- Basic assembler and VM running with working jumps, calls, arithmetic
- RET instruction fixed to correctly return from calls
- Debugging tools print label addresses and instruction assembly details
- String Support lenght based indexing
- Custom Binary Format
- Planning to add:
    - Memory copy, clear, and comparison instructions
    - Linker to compile Std-Lib functions together with user/Os code
    - Std-Lib functions like strcpy etc
    - Keyboard input handling
    - Enhanced memory management features
    - Visual Output
