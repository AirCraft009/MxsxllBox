STINTI 00000000


MOVI O1 128
MOVI O2 128
MOVI O3 1

# first line
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

LOOP:
JMP LOOP



