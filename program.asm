.program

# bootloader
STZ
MOVA O1 print_1
CALL _spawn

MOVA O1 print_2
CALL _spawn

MOVA O1 print_3
CALL _spawn

JMP _init_scheduler
HALT

print_1:
    MOVI R1 10
    MOVI O1 1
    PRINT R1
    CALL _yield
    JMP print_1

print_2:
    MOVI R1 20
    MOVI O1 1
    PRINT R1
    CALL _yield
    JMP print_2

print_3:
    MOVI R1 30
    MOVI O1 1
    PRINT R1
    CALL _yield
    JMP print_3




