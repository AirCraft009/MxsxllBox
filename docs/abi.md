# MxsxllBox ABI


## Chapter 1 Low Level System Information

### Chapter 1.1 Data  representation

- A byte refers to 8 bits of data
- A word refers to 2 bytes or 16 bits of data
- A doubleword or dw refers to 32 bits of data
- All registers are one word wide
- FB refers to [Frame-Buffer](https://de.wikipedia.org/wiki/Framebuffer)

### Chapter 1.2 Strings

- Strings are byte sequences mapped to the ASCII code
- Strings are lenght prefixed
  - The lenghts are a word
- Strings can both by printed to console and Screen

### Chapter 1.3 Asm Operand Structure

- Generally The Mxsxll asmebly standard follows the Intel syntax
  - Operand destination source
  - MOV R0 R1 (moves the val of R1 into R0)
  - An exception can be found with [STORE](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/instruction-set.md#important-----readme)

## Chapter 2 General Calling Conventions

### Chapter 2.1 Instruction derivs

The format of an instruction for the MxsxllBox cpu is structured as following

- Double Register operation

    ``Operation Reg1(modify) Reg2 (value)``\
    ``MOV R1 R2``

- Single Register operation

    ``Operation Reg``\
    ``GF R0`` - get flags to R0


- Register and Immediate

    ``Operation Reg1(modify) imm(value)``\
    ``ADDI R1 12``


- Label op

    ``Operation label``\
    ``JMP _start``
- Single OP

    ``Operation``\
    ``STC`` - set carry

- Special OP's
    
    ``Operation Reg label``- ``MOVA R0 _start`` - get label addr

All Opcodes with descriptions can be found [here](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/instruction-set.md)

### Chapter 2.2 Std. Lib Calling Conventions

- Arguments are always passed via Registers (expand to stack in future)
- The O - Registers (O1-O6)
- Return values are almost always in the O1 Register 
  - Exceptions are documented in the [std lib docs](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md)
- If there is more than one return value it is passed in O1 - ON
- Currently no std. lib functions exist with more than 2 returns

### Chapter 2.3 User lib conventions

- It's encouraged to use the O1 - O6 Registers 
- There is scratch space on registers K1 - K15
  - These Registers are not saved on context switches
  - They can be used to communicate between tasks
  - more details [here](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/scheduling.md)
- T1 - T6 should never be used in User code they're reserved for Scheduling


## Chapter 3 Memory

### Chapter 3.1 - Interacting with the Stack

- The Stack's address-space 
  - start: 0x6000
  - end: 0x7FFF
  
- Segmented in 9 equally sized blocks
    - each task has its own stack
    - currently accessing memory outside the correct stack segment is possible
    - If no scheduling is used accessing memory outside the 0th segment is safe behaviour


- PUSH & POP
  - these instructions get/put word sized data off/on the stack
  - these instructions modify the stackptr

    
- normal READ-BYTE/WORD & WRITE-BYTE/WORD
  - these addresses are mapped to [EEPROM]()
  - EEPROM should contain the bitmap font(bytes 0 - 255)

### Chapter 3.2 - Interacting with the Heap

- The heaps address-space 
  - start: 0x2000
  - end: 0x5FFF
- The heaps writeable address-space
  - start : 0x259C
  - end: 0x5D9C


- The Heap is managed dynamically
- use ``_alloc`` from stdlib/sys.obj 
  - allocates memory on the heap [definition](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md#sys-functions)
  - only write to the address returned in O1
  - the two bytes preceding O1 are the ammount of blocks this alloc entailed.
    - if this information gets lost freeing is impossible
  - !! Allocation is only theoretical not enforced !!
    - You are able to write to any addr allocated or not though it is strongly discouraged


- use ``_free`` from stdlib/sys.obj
  - frees memory from the heap [definition](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md#sys-functions)
  - double frees are undefined
    - free reads uninitialized mem (undefined behaviour)
  - always free the Address that was returned to you by ``_alloc``
  

### Chapter 3.3 - Interacting with the Framebuffer

- The FB's address-space
    - start: 0x8000
    - end: 0xBFFF

- The Frame Buffer is directly connected to the visual output

- Each Byte represents 4 Bytes
  - 1 Pixel is represented by 2 bits 
  - each pixel can have 4 different colors
    - 0/00 : Black
    - 1/01 : White
    - 2/10 : Red
    - 3/11 : Blue

- Singular pixels can be written to with the DRAWPX instruction
- A section (4 pixels / 1Byte in FB) can be set with the STOREBLOCK instruction

- stdlib/io.obj provides [``_draw_char``](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md#io-functions) & [``_draw_string``](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md#io-functions)
  - They allow easy visualisation of chars
  - They require EEPROM to be loaded with a 8x8 bitmap font (bytes 0 - 255) 




