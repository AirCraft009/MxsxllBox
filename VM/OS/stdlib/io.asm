_printstr:
    PRINTSTR O2
    RET

_printchar:
    LOADB O1 O2
    PRINT O1
    RET

_readchar:
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
    MOVI O1 255         #if buffer is empty load 256 into it (max byte val)
    STZ
    RET

_setP_xy:               # short for set Pixel x; y using x(O1) for col and y(O2) for row. val(O3)
    CALL _get_Dimension
    MUL O2 O6           # turn the row into an index
    ADD O2 O1           # now add the col number to get the final index

    JMP _setP_i         # JMP to index-based method

#TODO: REMEMBER IMPORTANT: THE RELATIVE OFFSET IS IN PIXEL NOT IN BYTE
_setP_i:                # short for set Pixel index(O2) using a relative index 0 - (16 * 1024)*4. val(O3)
    CALL _get_Ppb       # each byte has 4 pixels (2 bpp) so the relative index
    MOV O1 O2
    MOD O1 O6           # By moding you can see which bit in the byte has to be targeted
    CALL _get_Bpp       # now multiply that by the Bpp(0*2 = bit(0-1); 1*2 = bit(2-3); ...)
    MUL O1 O6
    DIVI O2 4           # Divide by 4 to get the actual byte
    CALL _get_video_start
    ADD O2 O6
    LS O3 O1            # leftshift the bits to the correct pos
    LOADB O1 O2         # get the current val to not override the other bits
    MOVI O6 255         # make O6 11111111 to invert the bits
    XOR O6 O3
    AND O1 O6           # Now use the mask to clear the bits at the currpos
    OR O3 O1            # Combine both bytes
    STOREB O3 O2        # Store the byte
    RET

_draw_rect:             # x(O1), y(O2), colorval(O3), len(O4), height (O5)   # the pos refers to the upper left corner
    CALL _get_Dimension
    MOV O6 O2
    ADD O6 O5           # to make it possible to check
    JMP DRAW_RECT_LOOP

DRAW_RECT_LOOP:
    PUSH O6
    CMP O2 O6
    JZ END_RECT

    CALL _draw_line_hori
    ADDI O2 1
    POP O6
    JMP DRAW_RECT_LOOP

END_RECT:
    RET

LINE_SETUP:
        POP O6
        PUSH O4             # push the original lenght to recover at the end
        PUSH O1             #-""-
        PUSH O2             # -""-
        PUSH O3             # -""-
        PUSH O6
        RET

_draw_line_vert:        # x(O1), y(O2), colorval(O3) len(O4) pos refers to the top
    CALL LINE_SETUP
    ADD O4 O2
    PUSH O1
    PUSH O3
    JMP DRAW_LINE_VERT_LOOP


DRAW_LINE_VERT_LOOP:
    MOV O5 O2
    CMP O2 O4
    JZ END_LINE

    CALL _setP_xy
    MOV O2 O5
    ADDI O2 1           # goto next line
    POP O3
    POP O1
    PUSH O1
    PUSH O3
    JMP DRAW_LINE_VERT_LOOP

_draw_line_hori:        # x(O1), y(O2), colorval(O3) len(O4) pos refers to the left side
    CALL LINE_SETUP
    ADD O4 O1           # add to be able to check progress
    PUSH O2             # Push O2 on the stack to temp-save it
    PUSH O3
    JMP DRAW_LINE_HORI_LOOP

DRAW_LINE_HORI_LOOP:
    MOV O5 O1
    CMP O1 O4
    JZ END_LINE

    CALL _setP_xy
    MOV O1 O5
    ADDI O1 1
    POP O3
    POP O2
    PUSH O2
    PUSH O3
    JMP DRAW_LINE_HORI_LOOP

END_LINE:
    POP O3      # POP TEMP
    POP O2      # POP TEMP
    POP O2      # POP TEMP
    POP O2      # POP TEMP
    POP O3      # Restore
    POP O2      # Restore
    POP O1
    POP O4      # Restore
    RET


_draw_char:    # 0 - 8192 KB is mapped to the char space  ; O1(x) O2(y) O3 (colorval) O4(charNum)
    MOVI O5 0
    MULI O4 8
    JMP DRAW_CHAR_LOOP

DRAW_CHAR_LOOP:
    PUSH O5
    CMPI O5 8       # char is 8x8
    JZ END_CHAR

    PUSH O4         # save temporarily
    ADD O4 O5

    LOADB O4 O4
    CALL _draw_line_mask
    POP O4
    POP O5
    ADDI O5 1       # counter ++
    ADDI O2 1       # y ++
    JMP DRAW_CHAR_LOOP

END_CHAR:
    DIVI O4 8
    SUBI O2 8
    POP O5
    RET

_draw_line_mask:    # draws one line 8 pixels wide and only draws the pixels that are 1; # O1(x), O2(y), O3(colorval), O4(mask)
    CALL LINE_SETUP
    PUSH O4
    PUSH O2             # Push O2 on the stack to temp-save it
    PUSH O3
    PUSH O1
    MOVI O5 0
    JMP DRAW_LINE_MASK_LOOP

DRAW_LINE_MASK_LOOP:
        CMPI O5 8
        JZ END_LINE
        RS O4 O5        # leftshift by x bytes
        TSTI O4 1       # check if the 0'th bit is active
        JNZ CONTINUE_DRAW_LINE_MASK_LOOP

        CALL _setP_xy
        JMP CONTINUE_DRAW_LINE_MASK_LOOP

CONTINUE_DRAW_LINE_MASK_LOOP:
        POP O1
        POP O3
        POP O2
        POP O4
        ADDI O1 1
        ADDI O5 1
        PUSH O4
        PUSH O2
        PUSH O3
        PUSH O1
        JMP DRAW_LINE_MASK_LOOP


_draw_string:       # O1(x) O2(y) O3(colorval) O4(stringpos)
    LOADW O6 O4     # Load len into O6
    ADDI O4 2       # Goto first char
    MOV O5 O4
    ADD O5 O6       # set to the end addr
    JMP DRAW_STRING_LOOP

DRAW_STRING_LOOP:
    CMP O4 O5
    JZ END_DRAWING_STRING
    PUSH O4
    PUSH O5
    LOADB O4 O4
    CALL _draw_char
    POP O5
    POP O4
    ADDI O4 1
    CALL _get_Dimension
    SUBI O6 16               # subtract 8(lenght of a char) so that if O1 is Greater then we can go to the next line
    CMP O1 O6
    JGE NEWLINE
    ADDI O1 8
    JMP DRAW_STRING_LOOP

NEWLINE:
    MOVI O1 8
    ADDI O2 8
    JMP DRAW_STRING_LOOP

END_DRAWING_STRING:
    RET

