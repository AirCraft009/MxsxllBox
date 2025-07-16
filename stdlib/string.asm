strcpy:
    LOADB O3 O1        # O3 = length
    ADDI O1 1          # O1 = O1 + 1  (back to first char)

    STOREB O3 O2      # store length byte at dst - 1
    ADDI O2 1          # O2 = O2 + 1

    # Now copy string bytes using O3 as counter
STRCPY_LOOP:
    TSTI O3 0         # test if length == 0
    JZ END_STRCPY

    LOADB O4 O1       # load byte from src
    STOREB O4 O2      # store byte to dst

    ADDI O1 1         # src++
    ADDI O2 1         # dst++
    SUBI O3 1         # length--

    JMP COPY_LOOP

END_STRCPY:
    RET
