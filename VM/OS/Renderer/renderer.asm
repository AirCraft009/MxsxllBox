_render_screen:
JMP RENDER_TABLE

RENDER_TABLE:
    JMP render_chars


render_chars:
    MOVI O1 1           # will change to the video-intrr.
    CALL _yield
    MOVI R3 15          # offset var to check if the offset flag is check
    MOV R4 VC           # the location of the current pos
    MOVI R5 0           # counter
    MOV R6 VS           # size
    DIVI R6 2
    MOVI O1 0           # x
    MOVI O2 0           # y
    MOVI O3 1           # color
    JMP RENDER_CHAR_LOOP

RENDER_CHAR_LOOP:
    MOVI R7 1
    CMP R5 R6
    JZ render_chars
    LOADW O4 R4         # load the val
    CMPI O4 0           # if zero the cell can be skipped always
    JZ UPDATE
    CALL _draw_string
    LOADW R7 O4         # get the len
    JMP RENDER_CHAR_LOOP

UPDATE:
    ADD R5 R7
    MULI R7 2
    ADD R4 R7
    MULI R7 4
    ADD O1 R7
    MODI O1 256
    CMPI O1 0
    JZ UPDATE_LINE
    JMP RENDER_CHAR_LOOP

UPDATE_LINE:
    ADDI O2 8
    JMP RENDER_CHAR_LOOP



