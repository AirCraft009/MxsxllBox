_read_line:
    MOVI R0 0       # R0 is the counter/len
    MOV R1 SL       # R1 has the len pos
    MOV R2 VS       # R2 has the size/cmp for end
    MOV R3 VC       # R3 is the pointer to the location
    MOV R4 SS       # R4 holds the  position of the stringsection
    MOV R5 R3
    ADDI R5 2       # R5 is the location of the offset
    STOREW SL R3
    JMP READ_LINE_LOOP


READ_LINE_LOOP:
    CALL _read_char
    JZ END_READ

    CMPI O1 13      # temporary check for newline
    JZ NEWLINE
    STOREB O1 R4
    STOREW R0 R1
    ADDI R4 1
    ADDI R0 1
    JMP READ_LINE_LOOP


END_READ:
    MOVI O1 2
    CALL _yield
    JMP READ_LINE_LOOP

NEWLINE:
    ADDI CY 8
    MOVI CX 0
    DIVI R0 16          # divide by 16 to calc in blocks
    ADDI R0 1
    MOV O2 R0
    CALL _alloc
    MOV O2 O1
    MOV O1 SL
    YIELD
    CALL _strcpy        # copy string from temp  to allcoated space
    STOREW SL 0        # store 0 at stringlen
    UNYIELD
    JMP _read_line