.entry
STINTI 11111111
MOVA O1 print_1
CALL _spawn
MOVA O1 print_2
CALL _spawn
JMP _init_scheduler


print_1:
    MOVA O1 print_3
    CALL _spawn
    MOVI O1 10
    JMP print_1_loop

print_1_loop:
    PRINT O1
    JMP print_1_loop

print_2:
    MOVI O1 11
    JMP print_2_loop

print_2_loop:
    PRINT O1
    MOVI O1 9
    CALL _yield
    JMP print_2_loop

print_3:
    MOVI O1 12
    JMP print_3_loop

print_3_loop:
    PRINT O1
    JMP print_3_loop