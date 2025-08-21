BOOTLOADER:
YIELD
MOVA O1 SimOtherStuff
CALL _spawn
MOVA O1 DrawInput
CALL _spawn
UNYIELD
JMP _init_scheduler

SimOtherStuff:
    JMP SimOtherStuff

DrawInput:
    MOVI R1 0
    MOVI O1 0
    MOVI O2 0
    MOVI O4 97
    MOVI O3 1
    JMP DRAW_NEW_CHAR

DRAW_NEW_CHAR:
    CMPI R1 32
    JZ NewLine
    JMP FINISH_DRAWING

FINISH_DRAWING:
    CALL _draw_char
    ADDI R1 1
    ADDI O1 8
    PUSH O1
    MOVI O1 2
    CALL _yield
    JMP DRAW_NEW_CHAR

NewLine:
    ADDI O2 8
    MOVI O1 0
    MOVI R1 0
    JMP  FINISH_DRAWING







