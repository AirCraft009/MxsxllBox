## Dynamic Memory Allocation

- Heap size: 16 KB
- writeable Heap size: 14 KB
  - start: 0x259C
  - end: 0x5D9C
- Uses a **bitmap allocator** with block size of 16 bytes
- Metadata stored in the first word of an allocation block
- `_alloc`: requests block counts (multiples of 16 bytes)
- `_free`: returns blocks to the heap
- Bitmap is stored at the beginning of Heap after the tasks
- Allocating isn't enforced
- It is strongly recommended to only interact with the heap via ``_alloc``