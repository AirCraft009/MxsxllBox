MOVI R0 100
MOVI R1 1
MOVI R2 192

main:
SUBI R0 10
JZ zeroloop
PRINT R0
JMP main

counter:
SUBI R2 1
PRINT R2
JZ returnLbl
JMP counter

returnLbl:
RET

zeroloop:
PRINT R1
CALL counter
JMP end

end:
PRINT R2
HALT