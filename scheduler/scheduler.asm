#ActiveTaskPtr = 8256
#Task Start =  8257
#Task len = 8317 buffer bit of 1'st task
#Task End = 8795
#Tasksize = 60 bytes

GET_ACTIVE_TASK:
    MOVI T6 8256
    RET

GET_TASK_START:
    MOVI T6 8257
    RET

GET_TASK_LEN:
    MOVI T6 8317
    LOADB T6 T6
    RET

GET_TASK_SIZE:
    MOVI T6 60
    RET

_yield:


_spawn:         # creates a task and saves it at 8321 ++
    # O1 = addr
    # ---TASK_LAYOUT---
    #
    #   PC    : uint16
    #   SP    : uint16
    #   R1    : uint16
    #   R2    : uint16
    #   .     :
    #   .     :
    #   O6    : uint16
    #   FLags : word  bit 0 = ZeroF 1 = CarryF
    #   State : byte  see __states
    #   len   : byte  only 1'st task others have nothing here


    CALL GET_TASK_LEN
    MOV T5 T6
    CALL GET_TASK_SIZE
    MUL T5 T6       # where can we start to write offset
    CALL GET_TASK_START
    ADD T5 T6       # actual start addr
    STOREW O1 T5    # store beginning of task
    ADDI T5 4       # move to next location - ignore SP when task created the first task will always have SP 0 the later ones will save it when switched to
    MOVI T1 0       # set loop counter 0
    CALL SAVE_REGS_LOOP



    SAVE_REGS_LOOP:

        CMP T1 27  # number of regs + 1
        JZ RETURN

        GRFN T1 T3      #TODO: implement GREG Get Register From Number
        STOREW T3 T5

        ADDI T5 2
        ADDI T1 1
        JMP SAVE_REGS_LOOP


    RETURN:
        RET



#   __states
#   running == 0
#   ready == 1
#   blocked == 2 # just general blockage if it's not implemented
#   KeyBoardBlocked == 3
#   timerBlocked == 4
#   to be...