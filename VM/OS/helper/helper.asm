_get_video_char_table_start:
    MOVI O6 49185
    RET

_get_video_char_table_size:
    MOVI O6 2048
    RET

_get_video_start:
    MOVI O6 32768
    RET

_get_video_end:
    MOVI O6 49151
    RET

_get_Dimension:     # the native resolution is 256 x 256 but is upscaled by 4
    MOVI O6 256     # add 1 to make offset calc easier
    RET

_get_Bpp:
    MOVI O6 2
    RET

_get_Ppb:   # Pixel per byte
    MOVI O6 4
    RET


_get_stack_start:
    MOVI T6 32767
    RET

_get_split_stack_size:
    MOVI T6 910
    RET

_get_task_size:
    MOVI T6 61      # task-size is actually 60 but 61 is returned to make the calculations easier
    RET


_get_active_task:
    MOVI T6 9088
    LOADB T6 T6
    RET

__get_active_task_location:
    MOVI T6 9088
    RET

_get_task_start:
    MOVI T6 9089
    RET

_get_task_len:
    CALL _get_task_len_pos
    LOADB T6 T6
    RET

_get_task_len_pos:
    MOVI T6 9149
    RET

_get_state_location:
    CALL _get_task_size
    MUL T1 T6       # get Offsets
    CALL _get_task_start
    ADD T1 T6       # get location
    SUBI T1 2
    RET

_get_state:          # T1 has the task number from 0
    CALL _get_state_location
    LOADB T1 T1     # LOAD into T1
    RET

_getReadPtr:
    MOVI O3 49152
    RET

_getWritePtr:
    MOVI O4 49153
    RET

_get_bitmap_start:
    MOVI O6 8192
    RET

_get_bitmap_end:
    MOVI O6 9087
    RET

_get_writeable_heap:
    MOVI O6 9629
    RET