BOOTLOADER:
MOVA O1 SimOtherStuff
CALL _spawn
MOVA O1 DrawInput
CALL _spawn
STINT 01111111
JMP _init_scheduler

SimOtherStuff:
    PRINT O1
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
    JMP DRAW_NEW_CHAR

NewLine:
    ADDI O2 8
    MOVI O1 0
    MOVI R1 0
    JMP  FINISH_DRAWING







