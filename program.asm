.program

# bootloader
STZ
MOVA O1 print_1
CALL _spawn
MOVA O1 print_2
CALL _spawn
CALL _init_scheduler

print_1:
    MOVI R1 1
    MOVI O1 2
    PRINT R1
    CALL _yield
    JMP print_1

print_2:
    MOVI R1 2
    MOVI O1 2
    CALL _yield
    JMP print_2








