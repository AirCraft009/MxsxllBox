_console:
    MOVI K1 0
    MOV R5 VC       # R5 is the ptr
    JMP _read_line


_read_line:
    STOREW SL R5
    MOVI R0 0
    MOV R1 SL       # R1 has the len pos
    MOV R4 SS       # R4 holds the  position of the stringsection
    JMP READ_LOOP


READ_LOOP:
    CALL _read_char
    JZ END

    CMPI O1 32
    JL SPECIAL_CHAR_CALLER
    ADDI R0 1
    STOREW R0 SL
    STOREB O1 R4
    ADDI R4 1
    JMP READ_LOOP

END:
    MOVI O1 1
    CALL _yield
    JMP READ_LOOP

SPECIAL_CHAR_CALLER:
    CALL _handle_special_char
    JZ READ_LOOP
    JMP  _read_line


_strinsert:         # insert a char val in O1 into string val O2 at pos O3
    LOADW O4 O2     # get the lenght into O4
    PUSH O4
    PUSH O3
    PUSH O2
    SUB O4 O3
    JC ERROR_OVERFLOW   # if it overflows (O3 >  O4)
    MOV O6 O1
    MOV O2 O1           # setup for memcpy
    ADDI O2 1
    MOV O3 O4
    CALL _memcpy
    POP O2
    POP O3              # restore position
    POP O4
    ADD O3 O2
    ADDI O4 1           # add 1 to prev len
    STOREW O4 O2
    STOREB O6 O3        # insert the char
    RET


ERROR_OVERFLOW:
    RET