#ActiveTaskPtr = 8256
#Task Start =  8257
#Task len = 8317 buffer bit of 1'st task
#Task End = 8795
#Tasksize = 60 bytes
#Max tasks = 9

GET_ACTIVE_TASK:
    MOVI T6 8256
    LOADB T6 T6
    RET

GET_TASK_START:
    MOVI T6 8257
    RET

GET_TASK_LEN:
    CALL GET_TASK_LEN_POS
    LOADB T6 T6
    RET

GET_TASK_SIZE:
    MOVI T6 61      # task-size is actually 60 but 61 is returned to make the calculations easier
    RET

GET_TASK_LEN_POS:
    MOVI T6 8317 
    RET

_scheduler:
    CALL GET_ACTIVE_TASK
    MOV T5 T6
    CALL GET_TASK_LEN
    MOV T4 T6
    CALL GET_TASK_START
    MOV T3 T6
    CALL GET_TASK_SIZE
    CALL ROUND_ROBIN


ROUND_ROBIN:
    TSTI T4 0
    JZ NO_TASK_FOUND

    CMP T4 T5
    JZ SKIP_TASK

    MOV T1 T4
    MUL T1 T6       # get Offsets
    ADD T1 T3       # get location

    SUB T1 T4       # remove extra byte
    SUB T1 T4       # go to state byte

    LOADB T1 T1
    CMPI T1 1       # check if state is ready
    JZ  FOUND_TASK

    SUBI T4 1
    JMP ROUND_ROBIN



SKIP_TASK:
    SUBI T4 1
    JMP ROUND_ROBIN

NO_TASK_FOUND:      # stay in infinity loop untill next  interrupt
    JMP NO_TASK_FOUND

FOUND_TASK:         # TODO: implement loading the task and jumping to it
    HALT






# going to optimize by converging _spawn / _yield

_yield:                 # cooperative yield( willingly from the current lbl)
    CALL GET_ACTIVE_TASK
    MOV T5 T6           # save activeTaskNum
    CALL GET_TASK_SIZE
    MUL T5 T6           # get the offset
    CALL GET_TASK_START
    ADD T5 T6           # set to correct addr
    POP T6              # pop the return addr/currPC

    GF T4
    ADDI T6 5           # add the offset of CALL instruction
    STOREW T6 T5
    ADDI T5 2
    GSP T6              # get Stack Pointer
    STOREW T6 T5
    ADDI T5 2
    MOVI T1 0

    CALL SAVE_REGS_LOOP

    STOREB T4 T5    # save flags from earlier
    ADDI T5 1       # move to state
    MOV T1 O1       # set state
    STOREB T1 T5

    CALL _scheduler






_spawn:         # creates a task and saves it
    # O1 = addr
    # ---TASK_LAYOUT---
    #
    #   PC    : uint16
    #   SP    : uint16
    #   R1    : uint16
    #   R2    : uint16
    #   .     :
    #   .     :
    #   T1    : uint16
    #   FLags : word  bit 0 = ZeroF 1 = CarryF
    #   State : byte  see __states
    #   len   : byte  only 1'st task others have nothing here

    GF T4

    CALL GET_TASK_LEN
    MOV T5 T6
    ADDI T6 1
    CMP T6 9
    JC TASKS_FULL
    CALL GET_TASK_SIZE
    MUL T5 T6       # where can we start to write offset
    CALL GET_TASK_START
    ADD T5 T6       # actual start addr
    STOREW O1 T5    # store beginning of task
    ADDI T5 4       # move to next location - ignore SP when task created the first task will always have SP 0 the later ones will save it when switched to
    MOVI T1 0       # set loop counter 0
    CALL SAVE_REGS_LOOP

    STOREB T4 T5    # save flags from earlier
    ADDI T5 1       # move to state
    MOVI T1 1       # set base state to ready maybe change
    STOREB T1 T5
    CALL GET_TASK_LEN
    ADDI T6 1
    MOV T5 T6
    CALL GET_TASK_LEN_POS
    PRINT T6
    STOREB T5 T6    # update lenght += 1
    RET


    SAVE_REGS_LOOP:

        CMPI T1 27  # number of regs + 1
        JZ RETURN

        GRFN T1 T3      #TODO: implement GREG Get Register From Number
        STOREW T3 T5

        ADDI T5 2
        ADDI T1 1
        JMP SAVE_REGS_LOOP

    TASKS_FULL:
        MOVI O2 1000
        PRINT O2
        STZ
        RET


    RETURN:
        RET



#   __states
#   running == 0
#   ready == 1
#   blocked == 2 # just general blockage if it's not implemented
#   KeyBoardBlocked == 3
#   timerBlocked == 4
#   to be...