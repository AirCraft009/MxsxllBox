STINTI 00000000


MOVI R0 32
MOVI R1 0       # counter x
MOVI R2 0       # counter y
MOVI R3 0       # x var
MOVI R4 0       # y var
MOVI O3 1       # color
JMP LOOP_y


RENDER_A:
    ADDI O2 1
    MOVI O4 6
    CALL _draw_line_vert

    SUBI O2 1
    ADDI O1 1
    MOVI O4 4
    CALL _draw_line_hori

    ADDI O2 3
    CALL _draw_line_hori

    ADDI O1 4
    SUBI O2 2
    MOVI O4 6
    CALL _draw_line_vert
    RET

LOOP_y:
    MOVI R1 0
    CMP R2 R0
    JZ FULL_SCREEN
    CALL LOOP_x
    ADDI R4 7
    ADDI R2 1
    JMP LOOP_y

LOOP_x:
    MOV O1 R3
    MOV O2 R4
    CMP R1 R0
    JZ RETURN
    CALL RENDER_A
    ADDI R1 1
    ADDI R3 8
    JMP LOOP_x

RETURN:
    RET

FULL_SCREEN:
    JMP FULL_SCREEN



