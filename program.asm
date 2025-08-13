.program
# Dies ist ein kleiner Test für meinen Scheduler der Ouput sollte 2 \n 1 unendlich wiederholt sein

# bootloader
MOVA O1 print_1 # MOVA = Move addr
CALL _spawn
MOVA O1 print_2
CALL _spawn
CALL _init_scheduler

print_1:
    MOVI R1 1
    PRINT R1
    MOVI O1 1
    CALL _yield
    JMP print_1

print_2:
    MOVI R1 2
    PRINT R1
    MOVI O1 1
    CALL _yield
    JMP print_2








