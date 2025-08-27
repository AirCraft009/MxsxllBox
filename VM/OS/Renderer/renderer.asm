_render_screen:
    MOVI O1 8
    CALL _yield
    MOV R2 VM       # get the video-mode
    MULI R2 5
    MOVA R1 render_table
    ADD R1 R2
    SPC R1

render_table:
JMP render_grid
# JMP render_pixels
# JMP render_optimized
# JMP render_dirty

render_grid:
    MOV R1 VC
    MOVI O3 1
    MOVI R2 0
    MOVI OR 0
    JMP RENDER_TABLE_LOOP

RENDER_TABLE_LOOP:
    CMP VS R2
    JZ END_RENDER

    LOADW R3 R1     # get the pos of the string
    CMPI R3 0
    JZ EMPTY_CELL
    CMPI OR 1
    JZ SKIP_OFFSET
    MOV O1 R1
    MOV O2 R1
    MODI O1 64      # x cord/8
    DIVI O2 64      # y cord/8
    MULI O1 8
    MULI O2 8
    MOV R4 R3
    MOVI R5 15
    RS R4 R5       # 16th bit is the offset bit
    MOV OR R4      # MOV it into the offset register
    CALL _draw_string
    ADDI R1 2
    JMP RENDER_TABLE_LOOP

SKIP_OFFSET:
    MULI O3 2
    ADD R1 O3       # skip
    JMP RENDER_TABLE_LOOP