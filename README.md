# MxsxllBox
    
## Overview

github.com/AirCraft009/MxsxllBox is an emulator of a 16 bit microcontroller. Designed with a 64 KB memory space divided into several segments, supporting a simple instruction set and features like labels, branching, function calls, and dynamic memory management.\
It has 32 Registers 18 general purpose(R), 6 Registers for syscalls(O), 6 for the Scheduler(T) and 2 for an interrupt(I)

## Dependencies
### Dependency: SDL2

- download [SDL2-devel-2.32.0-os_specific_format](https://github.com/libsdl-org/SDL/releases/tag/release-2.32.0)
- uncompress it 
- place it at the project root or somewhere else and add it to PATH
- Rename it to SDL2
- place SDL2.dll into project root or add it to PATH

## How to run

- ./VM-main.exe inputfile
- ./Debugger-main.exe inputfile

## Dynamic Memory Allocation

- Heap size: 16 KB
- writeable Heap size: 14 KB
- Uses a **bitmap allocator** with block size of 16 bytes
- Metadata stored in the first word of an allocation block
- `alloc`: requests block counts (multiples of 16 bytes)
- `free`: returns blocks to the heap
- Bitmap is stored at the beginning of Heap after the tasks

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
> the scheduler can save the current context/state of the cpu meaning(regs, PC, SP and flags) to memory.\
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

## Special Characters

- the strings in the controller are ASCII-chars
- I extended them to 256 to represent boxes
- the chars 0 - 31 are special and unprintable
- they have special actions like '\n' or '\t'

### Chararcter Table
| charNum | description                                                                                |
|---------|--------------------------------------------------------------------------------------------|
| 0       | null-byte: often used for null terminated strings<br/> tho it has no use here.             |
| 1       | up-arrow: used in consoles and files, <br/> to navigate the cursors                        |
| 2       | left-arrow: -""-                                                                           |
| 3       | down-arrow: -""-                                                                           |
| 4       | right-arrow: -""-                                                                          |
| 5       |                                                                                            |
| 6       |                                                                                            |
| 7       |                                                                                            |
| 8       | back-space: used when deleting                                                             |
| 9       | tab: enters the number of bytes in the TB reg <br/> that's can be changed in the settings. |
| 10      |                                                                                            |
| 11      |                                                                                            |
| 12      |                                                                                            |
| 13      | enter: used to create newlines and finish commands.                                        |
| 14      |                                                                                            |
| 15      |                                                                                            |
| 16      |                                                                                            |
| 17      |                                                                                            |
| 18      |                                                                                            |
| 19      |                                                                                            |
| 20      |                                                                                            |
| 21      |                                                                                            |
| 22      |                                                                                            |
| 23      |                                                                                            |
| 24      |                                                                                            |
| 25      |                                                                                            |
| 26      |                                                                                            |
| 27      | escape: used to go back in menues and close apps                                           |
| 28      |                                                                                            |
| 29      |                                                                                            |
| 30      |                                                                                            |
| 31      |                                                                                            |

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
- Flat system (no real protected kernel) programs can interface directly with hardware, but abstractions are provided
- Labels aren't assembled into byte code 
- The pre-Assembled code is outputted in a .obj file
- The .obj file is then linked together with other .obj files into a  final .bin file
- _The linker is given a `map[string]int`: that contains [.obj filepaths] location
- Heap allocator uses first-fit policy for simplicity and speed

---

## Current Status

- MxsxllBox running with working jumps, calls, arithmetic
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
