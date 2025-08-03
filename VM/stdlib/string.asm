_strcpy:
    LOADB O3 O1        # O3 = length
    ADDI O3 1


STRCPY_LOOP:          # Now copy string bytes using O3 as counter
    TSTI O3 0         # test if length == 0
    JZ END_STRCPY

    LOADB O4 O1       # load byte from src
    STOREB O4 O2      # store byte to dst

    SUBI O3 1         # length--
    ADDI O1 1         # src++
    ADDI O2 1         # dst++

    JMP STRCPY_LOOP

END_STRCPY:
    RET

_strlen:
    LOADB O1 O2      # load lenght into O2
    JMP END_STRLEN

END_STRLEN:
    RET



_strcmp:            # sets the 0 - Flag if they're equal, Carry - Flag if higher
    LOADB O3 O1      # load strlen 1
    LOADB O4 O2      # load strlen 2
    CMP O3 O4
    JNZ END_STRCMP
    JMP STRCMP_LOOP


STRCMP_LOOP:
    TSTI O3 0
    JZ END_STRCMP

    ADDI O1 1
    ADDI O2 1

    LOADB O4 O1
    LOADB O5 O2
    CMP O4 O5
    JNZ END_STRCMP

    SUBI O3 1
    JMP STRCMP_LOOP


END_STRCMP:
    RET



_strcat:
    LOADB O4 O1          # O4 = firstString lenght
    LOADB O5 O2          # O5 = secondStringlenght
    ADD O4 O5           # 04 = combinedStringLenght
    STOREB O4 O1         # Write combinedLenght to start of string
    ADD O1 O5           # Add len(firstString) to start of firststring; O1 = ptr(to last byte)
    ADDI O1 1           # O1 = ptr(first free byte)
    CALL _strcpy             # _strcpy copies str2 to O1
    SUBI O1 1
    SUB O1 O5           # return O1 to start of string
    JMP STRCAT_END


STRCAT_END:
    RET