_memset:        # sets a region of mem O1 = start addr, O2 = what to set, O3 = ammount of bytes to set

    JMP MEMSET_LOOP

MEMSET_LOOP:
    CMPI O3 0
    JZ END_MEMSET

    STOREB O2 O1
    ADDI O1 1
    SUBI O3 1
    JMP MEMSET_LOOP


END_MEMSET:
    RET


_memcpy:        # copies a region from O1 to O2 for O3 bytes

    JMP MEMCPY_LOOP

MEMCPY_LOOP:
    CMPI O3 0
    JZ END_MEMCPY

    LOADB O4 O1
    STOREB O4 O2

    ADDI O1 1
    ADDI O2 1
    SUBI O3 1

    JMP MEMCPY_LOOP

END_MEMCPY:
    RET
