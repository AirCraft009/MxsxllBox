_render_screen:
JMP RENDER_TABLE

RENDER_TABLE:
    JMP render_chars

render_chars:
    MOVI R1 0           # counter
    MOVI O1 0           # x cord
    MOVI O2 0           # y cord
    JMP RENDER_CHARS

RENDER_CHARS:
    MOVI R2 1           # ammount of cells to skip
    MOV R0 VC           # ptr to the current place in the Char-table
    CMP R1 VS           # cmp-to video-char table size
    JZ render_chars     # restart the proccess
    ADD R0 R1           # add R1 to get the current location
    LOADW O4 R0
    CMPI O4 0           # if the cell is zero it can be skipped
    JZ UPDATE_LOCATION  # automatically skips 1 cell because R2 is 1
    CALL _draw_string
    LOADW R2 O4         # get the lenght
    JMP UPDATE_LOCATION # now update by this lenght


UPDATE_LOCATION:
     MULI R2 2          # multiply by 2 because each cell is a word
     ADD R1 R2          # add to the counter
     MULI R2 4          # multiply by 2*4 = 8 to skip a cell in pixel
     ADD O1 R2
     MODI O1 256        # mod to make sure it's not overflowing
     JC UPDATE_LINE     # if O1 is bigger than 256
     JMP RENDER_CHARS

UPDATE_LINE:
    ADDI O2 1
    JMP RENDER_CHARS