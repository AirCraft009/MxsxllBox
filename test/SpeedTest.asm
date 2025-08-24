YIELD
JMP LOOP

LOOP:
    ADDI R1 1
    JC END      # jumps if an overlow happens
    JMP LOOP

END:
    HALT