_getReadPtr:
    MOVI O3 49152
    RET

_getWritePtr:
    MOVI O4 49153
    RET


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

    ADD O4 O2
    LOADB O1 O4
    RET

    
END_READCHAR_BUF_EMPTY:
    MOVI O1 255        #if buffer is empty load 256 into it (max byte val)
    STZ
    RET


