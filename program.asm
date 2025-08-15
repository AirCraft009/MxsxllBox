.entry
YIELD
MOVA O1 print_1
PRINT O1
CALL _spawn
MOVA O1 print_2
PRINT O1
CALL _spawn
JMP _init_scheduler


print_1:
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
    JMP print_2_loop