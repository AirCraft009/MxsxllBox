.program

# bootloader
STZ
MOVA O1 print_b
CALL _spawn
MOVA O1 print_c
CALL _spawn
CALL _init_scheduler



print_c:
    MOVI T2 2
    MOVI O1 1
    CALL _yield
    JMP print_c

print_b:
    PRINT O1
    CALL _readchar
    JNZ print_it
    JMP finish_print_b

finish_print_b:
    MOVI O1 2
    CALL _yield
    JMP print_b

print_it:
    PRINT O1
    JMP finish_print_b





