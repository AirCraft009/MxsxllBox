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
    DIVI O2 8           # Divide by 8 to get the actual byte
    LS O3 O1            # leftshift the bits to the correct pos
    LOADB O1 O2         # get the current val to not override the other bits
    MOVI O6 255         # make O6 11111111 to invert the bits
    XOR O6 O3
    AND O1 O6           # Now use the mask to clear the bits at the currpos
    OR O3 O1            # Combine both bytes
    STOREB O3 O2        # Store the byte
    RET

_draw_rect:             # x(O1), y(O2), colorval(O3), len(O4), width (O5)   # the pos refers to the upper left corner
    RET

_draw_line:             # x(O1), y(O2), colorval(O3) len(O4) pos refers to the left side
    ADD O4 O1           # add to be able to check progress
    PUSH O2             # Push O2 on the stack to temp-save it
    PUSH O3
    JMP DRAW_LINE_LOOP

DRAW_LINE_LOOP:
    MOV O5 O1
    CMP O4 O1
    JG FINISH_DRAWING

    CALL _setP_i
    MOV O1 O5
    ADDI O1 1
    POP O3
    POP O2
    JMP DRAW_LINE_LOOP





