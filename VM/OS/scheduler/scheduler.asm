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

_init_scheduler:
    CALL _get_task_len
    MOV T5 T6
    ADDI T5 1
    CALL __get_active_task_location
    STOREB T5 T6                    # store the len + 1 so it starts at index len()-1
    JMP _scheduler

_setup_scheduler:
    CALL _get_active_task
    MOV T4 T6
    CALL _get_task_start
    MOV T3 T6
    CALL _get_task_size
    RET

_scheduler:
    CALL _setup_scheduler
    PUSH T4
    JMP ROUND_ROBIN

SETUP_INTERRUPT_HANDLER:
        ADDI I1 23965               # add the interrupt table location to the current interrupt ID
        GPC T1                      # get the current PC
        ADDI T1 9                   # add the offset of the next 3 instructions -5 because the normal RET expects a CALL which has an instruction len of 5 so it doesn't get stuck in an infinity loop
        PUSH T1                     # Push onto the stack so the next RET call returns to 'CALL _scheduler'
        SPC I1                      # JMP without lbl
        CALL _setup_scheduler
        JMP FOUND_TASK

_unblock_tasks:             # T2 now has the type of task to be unblocked
    CALL _get_task_len
    MOV T4 T6               # T4 == counter
    JMP UNBLOCK_LOOP

UNBLOCK_LOOP:
    CMPI T4 0
    JZ RETURN

    MOV T1 T4
    CALL _get_state

    CMP T1 T2
    JZ  UNBLOCK
    SUBI T4 1

    JMP UNBLOCK_LOOP

UNBLOCK:
    MOV T1 T4
    SUBI T4 1
    CALL _get_state_location

    MOVI T3 1
    STOREB T3 T1
    JMP UNBLOCK_LOOP

ROUND_ROBIN:
    POP T4
    SUBI T4 1
    CMPI T4 0
    JZ WRAP_ARROUND

    MOV T1 T4
    CALL _get_state

    CMPI T1 1       # check if state is ready
    JLE FOUND_TASK

    JMP TEMP_UNYIELD

TEMP_UNYIELD:
    PUSH T4         # pushes T4 so if an interrupt occured it won't overwrite it
    UNYIELD
    MOVI I2 1
    YIELD
    MOVI I2 0
    JMP ROUND_ROBIN

WRAP_ARROUND:
    CALL _get_task_len
    MOV T4 T6
    ADDI T4 1
    PUSH T4
    CALL _get_task_size
    JMP ROUND_ROBIN

FOUND_TASK:
    CALL __get_active_task_location
    STOREB T4 T6            # Set Active Task
    CALL CALC_PC_FOR_ACTIVE_TASK
    MOV T4 T1
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
    PUSH T2                 # save temp because flags are easily overwritten
    ADDI T4 1               # go to state byte
    MOVI T2 0
    STOREB T2 T4            #store state
    POP T2
    SF T2
    UNYIELD
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
    YIELD
    GF T4                   # save the flags because they're easily overwritten
    CMPI O1 9               # see if the yield-code is equal to termination(9)
    JZ MARK_TASK_AS_DELETED
    MOVI I2 5
    JMP SAVE_TASK_YIELD

MARK_TASK_AS_DELETED:
    CALL _get_active_task
    MOV T1 T6
    CALL _get_task_len      # see if len has to be shortened
    CMP T1 T6
    JZ  SHORTEN_TASK

    CALL _get_state_location
    STOREB O1 T1        # store the termination
    CALL CALC_PC_FOR_ACTIVE_TASK
    MOVI T6 0
    STOREW T1 T6        # also change the PC to 0 this can't be because the bootloader is always at 0
    MOVI I2 0
    JMP _scheduler

SHORTEN_TASK:
    MOV T2 T6
    CALL _get_task_len_pos
    SUBI T2 1
    STOREB T2 T6
    MOVI I2 0
    JMP _scheduler

CALC_PC_FOR_ACTIVE_TASK:
    CALL _get_active_task
    MOV T1 T6
    JMP CALC_PC_FOR_TASK    # will return because of the RET statement in CALC_PC_FOR_TASK

CALC_PC_FOR_TASK:     # T1 has the task
    SUBI T1 1
    CALL _get_task_size
    MUL T1 T6
    CALL _get_task_start
    ADD T1 T6
    RET


_interrupt:
    YIELD
    GF T4           # save because flags are easily overwritten
    CMPI I2 1       # if I2 is set to 1 it means the round robin loop was interrupted so Saving the task again isn't necesarry/ would crash the program
    JZ BLOCK_SAVE
    MOVI I2 0
    JMP SAVE_TASK_INTERRUPT


BLOCK_SAVE:
    CALL _get_active_task
    MOV T1 T6
    CALL _get_state_location
    MOVI T6 1
    STOREB T1 T6        # make sure that the task is saved as ready
    MOVI I2 0
    JMP SETUP_INTERRUPT_HANDLER




SAVE_TASK:
    CALL _get_active_task
    SUBI T6 1
    MOV T5 T6           # save activeTaskNum
    CALL _get_task_size
    MUL T5 T6           # get the offset
    CALL _get_task_start
    ADD T5 T6           # set to correct addr

    POP T2              # when save_task is called from yield or interrupt this saves the return addr of it
    POP T6              # pop the return addr/currPC
    ADD T6 I2           # add the offset of CALL instruction or nothing depending on if it was yield or interrupt
    STOREW T6 T5

    ADDI T5 2
    GSP T6              # get Stack Pointer
    
    STOREW T6 T5

    ADDI T5 2
    MOVI T1 0


    CALL SAVE_REGS_LOOP

    PUSH T2

    STOREB T4 T5    # save flags from earlier
    ADDI T5 1       # move to state
    RET

SAVE_TASK_INTERRUPT:
    CALL SAVE_TASK
    MOVI T1 1       # set state
    STOREB T1 T5
    JMP SETUP_INTERRUPT_HANDLER

SAVE_TASK_YIELD:
    CALL SAVE_TASK
    STOREB O1 T5
    MOVI I2 0
    JMP _scheduler


FIND_NEXT_EMPTY_TASK:
    MOVI T2 0
    CALL _get_task_len
    JMP FIND_NEXT_EMPTY_TASK_LOOP

FIND_NEXT_EMPTY_TASK_LOOP:
    CALL _get_task_len
    CMP T2 T6
    JZ UPDATE_LEN


    MOV T1 T2
    ADDI T1 1
    CALL CALC_PC_FOR_TASK
    LOADW T1 T1
    CMPI T1 0       # see if task is "emtpy"
    JZ RETURN_TO_SPAWN

    ADDI T2 1

    JMP FIND_NEXT_EMPTY_TASK_LOOP

UPDATE_LEN:
    ADDI T2 1
    CALL _get_task_len_pos
    STOREB T2 T6
    SUBI T2 1
    MOV T5 T2
    RET


RETURN_TO_SPAWN:
    MOV T5 T2
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

    CALL FIND_NEXT_EMPTY_TASK

    CALL _get_task_size          # set- up PC
    MUL T5 T6                   # where can we start to write offset
    CALL _get_task_start
    ADD T5 T6                   # actual start addr
    STOREW O1 T5                # store beginning of task

    ADDI T5 2                   # set- up Stack
    CALL _get_split_stack_size
    MOV T1 T6
    CALL _get_task_len
    SUBI T6 1
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