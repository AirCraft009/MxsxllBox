_add:
    ADD O1 O2
    RET


_sub:
    SUB O1 O2
    RET

_mul:
    MUL O1 O2
    RET

_div:
    DIV O1 O2
    RET

_mod:
    MOD O1 O2
    RET

_max:
    CMP O1 O2
    JC 1_LARGER
    JNZ 2_LARGER
    RET

1_LARGER:
    RET

2_LARGER:
    MOV O1 O2
    RET

_min:
    CALL _max
    MOV O1 O2
    RET

_pow:

POWER_LOOP:

    CMPI O2 0
    JZ END_POWER

    MUL O1 O1
    SUBI O2 1

    JMP POWER_LOOP

END_POWER:
    RET

_inc:
    ADDI O1 1

_dec:
    SUBI O1 1

_clamp:         #Value is between two values input low O1 high O3 value O2 0 flag is set
    MOV O4 O1
    MOV O5 O2
    CALL _min
    CMP O1 O4
    JNZ END_CLAMP_LOW
    JMP CONTINUE_CLAMP

CONTINUE_CLAMP:
    MOV O1 O3
    MOV O2 O5
    CALL _max
    CMP O3 O1
    JNZ END_CLAMP_HIGH
    STZ
    CLC
    MOV O2 O5
    RET

END_CLAMP_HIGH:
    MOV O2 O3
    STC
    CLZ
    RET

END_CLAMP_LOW:
    MOV O2 O1
    CLC
    CLC
    RET


