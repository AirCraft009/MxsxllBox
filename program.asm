.program

# bootloader
STZ
MOVA O1 print_c
CALL _spawn
CALL _init_scheduler

print_c:
    MOV O1 2
    CALL _yield
    CALL _readchar
    JNZ printchar
    JMP print_c

printchar:
    PRINT O1
    JMP print_c




