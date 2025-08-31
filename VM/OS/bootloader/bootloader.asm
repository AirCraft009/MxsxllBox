_BOOT:
    STINTI 1
    YIELD
    CALL _get_video_char_table_start
    MOV VC O6       # table-start
    CALL _get_video_char_table_size
    MOV VS O6       # table-size
    MOV CY VC       # row 0
    MOVI SL 51234   # inputstring-lenght
    MOVI SS 51236   # inputstring-stringstart
    MOVA O1 _render_screen
    CALL _spawn
    MOVA O1 _console
    CALL _spawn
    JMP _init_scheduler