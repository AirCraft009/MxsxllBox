.program


main:
    CALL _readchar
    JNZ echo_char
    JMP main

echo_char:
    PRINT O1
    JMP main