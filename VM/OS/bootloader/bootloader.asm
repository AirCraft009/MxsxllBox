BOOT:
    CALL _get_video_char_table_start
    MOV VC O6
    CALL _get_video_char_table_size
    MOV VS O6