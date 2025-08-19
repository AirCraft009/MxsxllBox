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

_setP_xy:               # short for set Pixel x; y using x(O1) for col and y(O2) for row
    CALL _get_Dimension
    MUL O2 O6           # turn the row into an index
    ADD O2 O1           # now add the col number to get the final index
    JMP _setP_i         # JMP to index-based method


_setP_i:                # short for set Pixel index(O2) using a relative index 0 - (16 * 1024)*4
    CALL _get_Ppb       # each byte has 4 pixels (2 bpp) so the relative index
    MOV O1 O2
    MOD O1 O6           # By moding you can see which bit in the byte has to be targeted
    CALL _get_Bpp       # now multiply that by the Bpp(0*2 = bit(0-1); 1*2 = bit(2-3); ...)
    MUL O1 O6




