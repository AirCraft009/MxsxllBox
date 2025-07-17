.program
main:
    CALL _readchar
    JNZ echo_char
    JMP main

echo_char:
    MOVI R1 1
    ALLOC R1 O2
    STOREB O1 O2
    PRINT O2
    FREE O2