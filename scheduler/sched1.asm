#ActiveTaskPtr = 9088
#Task Start =  9089
#Task len = 9149 buffer bit of 1'st task
#Task End = 9628
#Tasksize = 60 bytes
#Per-task stack size = 910 bytes
#Max tasks = 9
#Interrupt Table = 	23965

# TODO : add deleting of tasks
# TODO : remove bloat out off scheduler


_get_stack_start:
    MOVI T6 32768
    RET

_get_split_stack_size:
    MOVI T6 910
    RET

_get_active_task:
    MOVI T6 9088
    LOADB T6 T6
    RET

_get_active_task_location:
    MOVI T6 9088
    RET

_get_task_start:
    MOVI T6 9089
    RET

_get_task_len:
    CALL _get_task_len_pos
    LOADB T6 T6
    RET

_get_task_size:
    MOVI T6 61      # task-size is actually 60 but 61 is returned to make the calculations easier
    RET

_get_task_len_pos:
    MOVI T6 9149
    RET

_init_scheduler:
    CALL _get_task_len
    MOV T5 T6
    ADDI T5 1
    CALL _get_active_task_location
    STOREB T5 T6                    # store the len + 1 so it starts at index len()-1
    JMP _scheduler_change

_scheduler:
    CALL _get_active_task
    MOV T4 T6
    CALL _get_task_start
    MOV T3 T6
    CALL _get_task_size
    RET

_scheduler_change:
    CALL _scheduler
    JMP ROUND_ROBIN

_scheduler_interrupt:
    ADDI I1 23965               # add the interrupt table location to the current interrupt ID
    GPC T1                      # get the current PC
    ADDI T1 9                   # add the offset of the next 3 instructions -5 because the normal RET expects a CALL which has an instruction len of 5 so it doesn't get stuck in an infinity loop
    PUSH T1                     # Push onto the stack so the next RET call returns to 'CALL _scheduler'
    SPC I1                      # JMP without lbl
    MOVI I2 25
    JMP POP_TASK_REGS           # pop them after saving the task

CONTINUE_SCHEDULER_INTERRUPT:
    CALL _scheduler
    JMP FOUND_TASK

GET_OFFSETS_FROM_TASK:
    MOV T1 T4
    MUL T1 T6       # get Offsets
    ADD T1 T3       # get location

    SUBI T1 2
    RET

ROUND_ROBIN:
    SUBI T4 1
    CMPI T4 0
    JZ WRAP_ARROUND

    CALL GET_OFFSETS_FROM_TASK
    LOADB T1 T1
    
    CMPI T1 1       # check if state is ready
    JZ  FOUND_TASK
    JNC FOUND_TASK  # will only happen if an alr. running task is the only option available so if all others are blocked

    JMP ROUND_ROBIN

WRAP_ARROUND:
    CALL _get_task_len
    MOV T4 T6
    ADDI T4 1
    CALL _get_task_size
    JMP ROUND_ROBIN

FOUND_TASK:
    CALL _get_active_task_location

    STOREB T4 T6            # Set Active Task
    CALL _get_task_size


    MUL T4 T6
    MOVI R1 1
    MUL R1 T6
    SUB T4 R1              # make sure to go to end of task - 1
    CALL _get_task_start
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

_yield:                     # cooperative yield( willingly from the current lbl)
    CALL SAVE_TASK_YIELD    # task is changed

    ADDI T5 1               # move to state
    MOV T1 O1               # set state
    STOREB T1 T5

    JMP _scheduler_change


_interrupt:                    # when this is called the interrupt id should aleready be loaded into I1 (interrupt 1)
    MOVI I2 32
    JMP PUSH_TASK_REGS         # push all task registers so if a yield get's interrupted the task registers stay


CONTINUE_INTERRUPT:
    CALL SAVE_TASK_INTERRUPT    # task is interrupted continues after handling
    JMP _scheduler_interrupt

PUSH_TASK_REGS:
    SUBI I2 1
    CMPI I2 25
    JZ CONTINUE_INTERRUPT

    GSP T6
    GRFN I2 T6

    PUSH T6

    JMP PUSH_TASK_REGS


POP_TASK_REGS:
    ADDI I2 1
    CMPI I2 32
    JZ CONTINUE_SCHEDULER_INTERRUPT

    POP T6
    SRFN I2 T6

    JMP POP_TASK_REGS




SAVE_TASK_YIELD:
    CALL SET_ADDR
    POP T2              # pop the return addr of the previous lbl either _yield or _interrupt
    POP T6              # pop the return addr/currPC of the lbl that called the actual _yield/_interrupt
    CALL SAVE_CONTEXT
    PUSH T2             # add the return addr of the previous lbl so it can be returned to
    RET

SAVE_TASK_INTERRUPT:
    CALL SET_ADDR
    SUBI T2 5           # offset the CALL instructionthat will be added on the interrupt is only possible after every step
    CALL SAVE_CONTEXT
    RET

SET_ADDR:
    CALL _get_active_task
    SUBI T6 1
    MOV T5 T6           # save activeTaskNum
    CALL _get_task_size
    MUL T5 T6           # get the offset
    CALL _get_task_start
    ADD T5 T6           # set to correct addr
    RET

SAVE_CONTEXT:
    GF T4
    ADDI T6 5           # add the offset of CALL instruction
    STOREW T6 T5
    ADDI T5 2
    GSP T6              # get Stack Pointer
    STOREW T6 T5
    ADDI T5 2
    MOVI T1 0

    CALL SAVE_REGS_LOOP

    STOREB T4 T5        # save flags from earlier
    RET








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

    CALL _get_task_len
    MOV T5 T6
    ADDI T6 1

    CMPI T6 9
    JC TASKS_FULL

    CALL _get_task_size          # set- up PC
    MUL T5 T6                   # where can we start to write offset
    CALL _get_task_start
    ADD T5 T6                   # actual start addr
    STOREW O1 T5                # store beginning of task

    ADDI T5 2                   #set- up Stack
    CALL _get_split_stack_size
    MOV T1 T6
    CALL _get_task_len
    MUL T1 T6
    CALL _get_stack_start
    SUB T6 T1
    STOREW T6 T5

    ADDI T5 2                   # goto regs. location
    MOVI T1 0                   # set loop counter 0
    CALL SAVE_REGS_LOOP

    STOREB T4 T5    # save flags from earlier
    ADDI T5 1       # move to state
    MOVI T1 1       # set base state to ready maybe change
    STOREB T1 T5
    CALL _get_task_len
    ADDI T6 1
    MOV T5 T6
    CALL _get_task_len_pos
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