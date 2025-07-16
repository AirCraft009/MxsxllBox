main:
    MOVI R2 5
    ALLOC R2 O1
    MOVI R2 0
    STRING R2 O1 "Hello you lachs"
    PRINTSTR O1
    ALLOC R2 O2
    CALL _strcpy
    HALT