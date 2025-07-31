interruptTable:
    JMP _interrupt          # handles saving the current task and goes back to the right offset
    PRINT T1
    JMP _keyboard_handler   # always jump so that the next RET call goes back to the scheduler
   #JMP _timer_handler


_keyboard_handler:
    MOVI T2 2
    PRINT T1
    CALL _unblock_tasks
    RET