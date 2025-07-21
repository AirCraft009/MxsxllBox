#ActiveTaskPtr = 9088
#Task Start =  9089
#Task len = 9149 buffer bit of 1'st task
#Task End = 9628
#Tasksize = 60 bytes
#Per-task stack size = 910 bytes
#Max tasks = 9

#going to reformat the whole system but focusing on getting it to work first
# TODO : add deleting of tasks


GET_STACK_START:
    MOVI T6 32768
    RET

GET_SPLIT_STACK_SIZE:
    MOVI T6 910
    RET

GET_ACTIVE_TASK:
    MOVI T6 9088
    LOADB T6 T6
    RET

GET_ACTIVE_TASK_LOCATION:
    MOVI T6 9088
    RET

GET_TASK_START:
    MOVI T6 9089
    RET

GET_TASK_LEN:
    CALL GET_TASK_LEN_POS
    LOADB T6 T6
    RET

GET_TASK_SIZE:
    MOVI T6 61      # task-size is actually 60 but 61 is returned to make the calculations easier
    RET

GET_TASK_LEN_POS:
    MOVI T6 9149
    RET

_init_scheduler:
    CALL GET_TASK_LEN
    MOV T5 T6
    ADDI T5 1
    CALL GET_ACTIVE_TASK_LOCATION
    STOREB T5 T6                    # store the len + 1 so it starts at index len()-1
    JMP _scheduler


_scheduler:
    CALL GET_ACTIVE_TASK
    MOV T4 T6
    CALL GET_TASK_START
    MOV T3 T6
    CALL GET_TASK_SIZE
    CALL ROUND_ROBIN


ROUND_ROBIN:
    SUBI T4 1
    CMPI T4 0
    JZ WRAP_ARROUND

    MOV T1 T4
    MUL T1 T6       # get Offsets
    ADD T1 T3       # get location

    SUBI T1 2


    LOADB T1 T1


    CMPI T1 1       # check if state is ready
    JZ  FOUND_TASK
    JNC FOUND_TASK  # will only happen if an alr. running task is the only option available so if all others are blocked

    JMP ROUND_ROBIN

WRAP_ARROUND:
    CALL GET_TASK_LEN
    MOV T4 T6
    ADDI T4 1
    CALL GET_TASK_SIZE
    JMP ROUND_ROBIN

FOUND_TASK:                 # TODO: implement loading the task and jumping to it
    CALL GET_ACTIVE_TASK_LOCATION

    STOREB T4 T6            # Set Active Task
    CALL GET_TASK_SIZE


    MUL T4 T6
    MOVI R1 1
    MUL R1 T6
    SUB T4 R1              # make sure to go to end of task - 1
    CALL GET_TASK_START
    ADD T4 T6

    JMP LOAD_TASK



LOAD_TASK:                  # Return state of the program to last task
    LOADW T3 T4             # get PC
    ADDI T4 2               # go to SP-byte
    LOADW T2 T4

    SSP T2                  # set SP
    ADDI T4 2               # go to first register
    MOVI T1 0
    CALL RESTORE_REGS_LOOP

    LOADB T2 T4
    SF  T2                  # set flags
    ADDI T4 1               # go to state byte
    MOVI T2 0
    STOREB T2 T4            #store state
    SPC T3                  # finally set the PC/should jump

    MOVI T2 1000            # if smth somehow went wrong
    PRINT T2                # print status code 1000
    HALT


RESTORE_REGS_LOOP:

        CMPI T1 26  # number of regs + 1
        JC RETURN

        LOADW T6 T4

        SRFN T1 T6      #SRFN Set Register From Number

        ADDI T4 2
        ADDI T1 1

        JMP RESTORE_REGS_LOOP


# going to optimize by converging _spawn / _yield

_yield:                 # cooperative yield( willingly from the current lbl)
    CALL GET_ACTIVE_TASK
    SUBI T6 1
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

    JMP _scheduler






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
    #   FLags : byte  bit 0 = ZeroF 1 = CarryF
    #   State : byte  see __states
    #   len   : byte  only 1'st task others have nothing here

    GF T4

    CALL GET_TASK_LEN
    MOV T5 T6
    ADDI T6 1

    CMPI T6 9
    JC TASKS_FULL

    CALL GET_TASK_SIZE          # set- up PC
    MUL T5 T6                   # where can we start to write offset
    CALL GET_TASK_START
    ADD T5 T6                   # actual start addr
    STOREW O1 T5                # store beginning of task

    ADDI T5 2                   #set- up Stack
    CALL GET_SPLIT_STACK_SIZE
    MOV T1 T6
    CALL GET_TASK_LEN
    MUL T1 T6
    CALL GET_STACK_START
    SUB T6 T1
    STOREW T6 T5

    ADDI T5 2                   # goto regs. location
    MOVI T1 0                   # set loop counter 0
    CALL SAVE_REGS_LOOP

    STOREB T4 T5    # save flags from earlier
    ADDI T5 1       # move to state
    MOVI T1 1       # set base state to ready maybe change
    STOREB T1 T5
    CALL GET_TASK_LEN
    ADDI T6 1
    MOV T5 T6
    CALL GET_TASK_LEN_POS
    STOREB T5 T6    # update lenght += 1
    RET


    SAVE_REGS_LOOP:

        CMPI T1 26  # number of regs + 1
        JC RETURN

        GRFN T1 T3      #GREG Get Register From Number
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