# MxsxllBox VM

## Overview

MxsxllBox is a custom virtual machine (VM) designed with a 16 KB memory space divided into several segments, supporting a simple instruction set and features like labels, branching, function calls, and dynamic memory management.

---

## Memory Layout

| Segment      | Size | Address Range  | Description                    |
|--------------|------|----------------|--------------------------------|
| Program      | 4 KB | 0x0000 - 0x0FFF| Code and instructions          |
| Heap         | 4 KB | 0x1000 - 0x1FFF| Dynamic memory allocation      |
| Video Memory | 5 KB | 0x2800 - 0x3FFF| Video buffer                   |
| Stack        | 2 KB | 0x2000 - 0x27FF| Stack for function calls, etc. |
| I0-Reserved  | 1 KB | 0x4000 - 0x43FF| Reserved for I/O Keybord et.   |
Stack initializes at 0x1000

---

## Instructions and Parsing

- Instructions vary in size, mostly 2 or 4 bytes
- Labels support for jumps and calls, with addresses resolved in two passes
- Opcodes for arithmetic, load/store, jumps, calls, push/pop, print, and halt
- Support for string-related instructions and memory operations planned

---

## Dynamic Memory Allocation

- Heap size: 4 KB
- Uses a **bitmap allocator** with block size of 16 bytes
- Metadata stored in the first word of an allocation block
- `alloc`/`malloc` requests block counts (multiples of 16 bytes)
- `free` returns blocks to the heap
- Bitmap is stored outside the heap in VM internal structures


## Std. Lib. 

- Labels that are saved into the library Region of the Program Space
- ProgramStdLibStart = 0x0C01 
- ProgramEnd         = 0x0FFF

### String functions

- strcpy O1 = loc O2 = dest
- strlen O1 = loc O2 = lenght
- strcmp O1 = loc O2 = val.
- strcat O1 = loc1 O2 = loc2 O1 = conc.

### 

---

## Design Decisions

- Memory layout optimized for simplicity and performance
- Labels are stored separate from instructions during assembly for quick address resolution
- Two-pass assembler: first collects labels, second generates machine code
- Jump and call instructions jump to label addresses, verified at runtime
- Heap allocator uses first-fit policy for simplicity and speed

---

## Current Status

- Basic assembler and VM running with working jumps, calls, arithmetic
- RET instruction fixed to correctly return from calls
- Debugging tools print label addresses and instruction assembly details
- String Support lenght based indexing
- Planning to add:
    - Memory copy, clear, and comparison instructions
    - Std-Lib functions like strcpy etc
    - Keyboard input handling
    - Enhanced memory management features
    - Visual Output
