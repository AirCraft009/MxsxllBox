_read_input:
    MOV R1 SS       # R1 will be the ptr to the next open spot
    MOV R2 SL       # R2 will hold the len location
    MOVI R0 0       # R0 will hold the current lenght of the string
    STOREW R0 R2    # init the string to zero so the screen doesn't display old data
    STOREW R2 CF    # store the location of the string
    MOV R4 CF
    ADDI R4 2       # R4 is the location of the next word where the offset is saved
    JMP READ_INPUT_LOOP

READ_INPUT_LOOP:
    CALL _read_char
    JZ END_TASK

    CMPI O1 10          # will later be set to 32 to check for any special chars 10 is \n
    JZ HANDLE_NEWLINE   # will be switched to a table that handles all special chars
    JNC READ_INPUT_LOOP # rn all other chars will be ignored
    YIELD
    STOREB O1 R1
    STOREW R4 R0
    ADDI R0 1
    STOREW R0 R2
    UNYIELD
    JMP READ_INPUT_LOOP


END_TASK:
    MOVI O1 2
    CALL _yield
    JMP READ_INPUT_LOOP

HANDLE_NEWLINE:
    DIVI R4 64
    ADDI R4 1
    MULI R4 64
    MOV CF R4
    JMP READ_INPUT_LOOP