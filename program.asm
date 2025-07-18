.program


main:
    CALL _readchar
    JNZ echo_char
    JMP main

echo_char:
    MOVI R1 1
    PRINT O1
    JMP main