#Bitmap_Start = 8192
#Bitmap_End = 8255 incl
#Writeable_Heap = 8796

GET_BITMAP_START:
    MOVI O6 8192
    RET

GET_BITMAP_END:
    MOVI O6 8255
    RET

GET_WRITEABLE_HEAP:
    MOVI O6 8796
    RET


_alloc:                     # O2 is the ammount O1 will be the start
                            # Allocates number of bytes*blocksize(16)
    CALL GET_BITMAP_START   # location of the bitmap first 128 bytes of heap
    MOV O3 O6
    CALL GET_BITMAP_END
    MOV O5 O2               # Store O2 for RESET

    ALLOC_BITMAP_LOOP:
        TSTI O2 0
        JZ ALLOCATE

        LOADB O4 O3      # see if block is alr. set
        TSTI O4 0       # if it's 0 then the space is free
        JNZ RESET_BITMAP_LOOP

        SUBI O2 1       # counter --
        ADDI O3 1       # src ++
        JMP ALLOC_BITMAP_LOOP

    RESET:
        MOV O2 O5
        RET

    ALLOCATE:
        CALL RESET
        JMP ALLOCATE_LOOP

    ALLOCATE_LOOP:      # set all Bitmap Entries to 1 (full)
        TSTI O2 0
        JZ SUCCES_ALLOC

        MOVI O4 1
        STOREB O4 O3    # set  to 1

        ADDI O3 1
        SUBI O2 1
        JMP ALLOCATE_LOOP


    RESET_BITMAP_LOOP:
        CALL RESET
        ADD O3 O2       # check the final location if all next blocks are free
        CMPI O3 O6      # If it's bigger than the bitmap end it fails Set 0 flag
        JC FAILED_ALLOC
        CLC             # make sure carry isn't set next time
        SUB O3 O2
        JMP ALLOC_BITMAP_LOOP # try again until fail

    SUCCES_ALLOC:
        CLC
        CLZ
        CALL RESET
        CALL GET_BITMAP_START
        SUB O3 O2       # subtract ammount from O3 to get the start of the region in bitmap
        SUB O3 O6       # subtract the addr to get the offset in blocks
        MULI O3 16      # multiply with blocksize to get offset from heapstart in bytes
        CALL GET_WRITEABLE_HEAP
        ADD O3 O6       # Add start of actually writeable heap. to get a ptr to the start
        STOREW O2 O3    # write the size to the first to bytes
        MOV O1 O3
        ADDI O1 2       # Add the offset two bytes to the final addr
        RET


    FAILED_ALLOC:
        STZ
        CLC
        MOVI O1 0       # Return 0  if nothing can be allocated
        RET


_free:  # frees a block of mem; O1 is the addr
    CALL GET_WRITEABLE_HEAP
    SUBI O1 2           # go to len
    LOADW O2 O1         # get len
    SUB O1 O6        # get offset in bytes from heapStart
    DIVI O1 16          # divide by 16 to get offset in blocks
    CALL GET_BITMAP_START
    ADD O1 O6        # go to addr overlayed in Bitmap
    MOVI O3 0
    CALL FREE_LOOP

    FREE_LOOP:
        TSTI O2 0
        JZ END_FREE

        STOREB O3 O1   # set bitmap to 0

        ADDI O1 1
        SUBI O2 1

        JMP FREE_LOOP

    END_FREE:
        RET






_exit:          #will call the END instruction that'll shut down the machine
    RET

_halt:
    HALT