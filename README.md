# MxsxllBox
    
## Overview

MxsxllBox is an emulator of a 16 bit microcontroller. Designed with a 64 KB memory space divided into several segments, supporting a simple instruction set and features like labels, branching, function calls, and dynamic memory management.\
It has 32 Registers 18 general purpose(R), 6 Registers for syscalls(O), 6 for the Scheduler(T) and 2 for an interrupt(I)

## Usage
    ./MxsxllBox.exe InputFile.bin [--debug] 
    !! Always put flags after the positional inputFile arg 

## Dependencies
### Dependency: SDL2

- download [SDL2-devel-2.32.0-os_specific_format](https://github.com/libsdl-org/SDL/releases/tag/release-2.32.0)
- uncompress it 
- place it at the project root
- Rename it to SDL2
- place SDL2.dll into project root or add it to PATH

## Build Windows - Main
- set CGO_ENABLED=1
- set CC=gcc
- set CXX=g++
- set CGO_CFLAGS=-I%cd%/SDL2/.../include   (... replace with your sdl implementation dir )
- set CGO_LDFLAGS=-L%cd%/SDL2/.../ -lSDL2
- optionally
  - go clean -cache -modcache
- go build -v -o MxsxllBox.exe ./VM-main

## How to run

- MxsxllBox
- 

---




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
