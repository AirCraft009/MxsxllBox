# MxsxllBox ABI

## Chapter 1 - Low Level System Information

### Chapter 1.1 Data  representation

- a byte refers to 8 bits of data
- a word refers to 2 bytes or 16 bits of data
- a doubleword or dw refers to 32 bits of data

### Chapter 1.2 General Calling Conventions

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

All Opcodes with descriptions can be found [here]()

## Std. Lib

- Arguments are always passed via Registers (expand to stack in future)
- The O - Registers (O1-O6)
- Return values are almost always in the O1 Register 
  - Exceptions are documented in the [std lib docs](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md)
- If there is more than one return value it is passed in O1 - ON
- Currently no std. lib functions exist with more than 2 returns

## User lib conventions

- It's encouraged to use the O1 - O6 Registers 
- 