.program

# bootloader
STZ
MOVA O1 print_c
CALL _spawn
MOVA O1 print_b
CALL _spawn
CALL _init_scheduler


print_c:
    MOVI T2 2
    PRINT T2
    MOVI O1 1
    CALL _yield
    JMP print_c

print_b:
    MOVI T1 1
    PRINT T1
    MOVI O1 1
    CALL _yield
    JMP print_b




