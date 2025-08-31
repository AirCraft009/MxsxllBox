_printstr:
    PRINTSTR O2
    RET

_print_char:
    LOADB O1 O2
    PRINT O1
    RET

_read_char:
    CLZ
    CALL _getReadPtr
    CALL _getWritePtr
    LOADB O2 O3      #Read ptr
    LOADB O1 O4      #Write ptr

    CMP O1 O2
    JZ END_READCHAR_BUF_EMPTY
    ADDI O4 1
    ADD O4 O2
    LOADB O1 O4     # buffer isn't empty so load char val. into O1
    ADDI O2 1
    MODI O2 30
    STOREB O2 O3
    RET


END_READCHAR_BUF_EMPTY:
    MOVI O1 257         #if buffer is empty load 256 into it (max byte val)
    STZ
    RET

_draw_string:               # O1(x) O2(y) O3(color) O4 (Stringlocation)
    PUSH O4
    MOV O5 O4
    LOADW O6 O5             # LOAD the len
    ADDI O5 2               # go to the first char
    ADD O6 O5               # make O6 be the end
    JMP DRAW_STRING_LOOP

DRAW_STRING_LOOP:
    CMP O5 O6
    JZ END_STRING_LOOP
    PUSH O5
    PUSH O6
    LOADB O4 O5
    CALL _draw_char
    POP O6
    POP O5
    ADDI O1 8
    ADDI O5 1
    JMP DRAW_STRING_LOOP

END_STRING_LOOP:
    POP O4
    RET

_draw_char:                 # O1(x) O2(y) O3(color) O4(charnum)
    MULI O4 8               # get the addr of the char
    MOV O5 O4               # O5 is now the ptr and O4 can be used as the mask
    MOVI O6 0
    JMP DRAW_CHAR_LOOP

DRAW_CHAR_LOOP:
    CMPI O6 8
    JZ END_CHAR

    LOADB O4 O5
    CALL _draw_mask_line

    ADDI O2 1
    ADDI O5 1
    ADDI O6 1
    JMP DRAW_CHAR_LOOP

END_CHAR:
    SUBI O2 8
    RET



_draw_mask_line:            # O1(x) O2(y) O3(colorval) O4(mask)
    PUSH O1
    PUSH O3
    PUSH O4
    PUSH O5
    PUSH O6
    MOVI O5 0
    JMP DRAW_MASK_LINE_LOOP

DRAW_MASK_LINE_LOOP:
    CMPI O5 8               # lenght of a byte/char-line
    JZ END_DRAW_MASK_LINE

    MOV O6 O4
    RS O6 O5                # left shift
    MODI O6 2
    CMPI O6 1               # see if it's !even /the last bit is a 1
    JNZ CONTINUE_LOOP
    DRAWPX O1 O2
    JMP CONTINUE_LOOP

CONTINUE_LOOP:
    ADDI O1 1
    ADDI O5 1
    JMP DRAW_MASK_LINE_LOOP

END_DRAW_MASK_LINE:
    POP O6
    POP O5
    POP O4
    POP O3
    POP O1
    RET

