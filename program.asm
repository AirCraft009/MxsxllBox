.program

print_a:
    MOVI R1 1
    SPAWN
    PRINT R1
    MOVI O1 0
    CALL print_b
    YIELD
    JMP print_a

print_b:
    MOVI R1 2
    SPAWN
    PRINT R1
    YIELD
    JMP print_b