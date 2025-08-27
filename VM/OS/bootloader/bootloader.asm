_BOOT:
    YIELD
    CALL _get_video_char_table_start
    MOV VC O6
    CALL _get_video_char_table_size
    MOV VS O6
    MOV CY VC
    MOVA O1 _render_screen
    CALL _spawn
    MOVA O1 _read_input
    CALL _spawn
    JMP _init_scheduler