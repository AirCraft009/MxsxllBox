.program

MOVI R1 30
ALLOC R1 R2
PRINT R2
STRING R1 R2 "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
PRINTSTR R2

main:
    CALL _readchar
    JNZ echo_char
    JMP main

echo_char:
    MOVI R1 1
    ALLOC R1 O2
    STOREB O1 O2
    PRINT O1
    FREE O2
    JMP main