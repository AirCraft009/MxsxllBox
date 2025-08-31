_console:
    MOVI K1 0
    MOV R5 VC       # R5 is the ptr
    JMP _read_line


_read_line:
    STOREW SL R5
    MOVI R0 0
    MOV R1 SL       # R1 has the len pos
    MOV R4 SS       # R4 holds the  position of the stringsection
    JMP READ_LOOP


READ_LOOP:
    CALL _read_char
    JZ END

    CMPI O1 13
    JZ NEWLINE
    PRINT O1
    ADDI R0 1
    STOREW R0 SL
    STOREB O1 R4
    ADDI R4 1
    JMP READ_LOOP

END:
    MOVI O1 1
    CALL _yield
    JMP READ_LOOP

NEWLINE:
    ADDI R0 2       # add 2 for the lenght
    DIVI R0 16      # Divide by the block size
    ADDI R0 1       # add 1 to make sure there's enough space
    MOV O2 R0
    CALL _alloc     # allocate the space
    MOV O2 O1
    MOV O1 SL
    CALL _strcpy    # copy string to new location
    STOREW O2 R5
    ADDI R5 64
    STOREW K1 SL
    JMP _read_line