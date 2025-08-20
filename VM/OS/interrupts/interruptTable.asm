_interruptTable:
    JMP _interrupt          # handles saving the current task and goes back to the right offset
    JMP _keyboard_handler   # always jump so that the next RET call goes back to the scheduler
    JMP _timer_handler


_keyboard_handler:
    MOVI T2 2
    CALL _unblock_tasks
    RET

_timer_handler:
    POP T3
    JMP _scheduler
 