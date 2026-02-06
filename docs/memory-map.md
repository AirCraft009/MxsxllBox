## Memory Layout

## Memory Layout (64 KB)

| Segment             | Size  | Address Range        | Description                                |
|---------------------|-------|----------------------|--------------------------------------------|
| **Program**         | 8 KB  | `0x0000` – `0x1FFF`  | Code and instructions                      |
| **Heap**            | 16 KB | `0x2000` – `0x5FFF`  |                                            |
| ├─ Tasks            | 540 B | `0x2381` - `0x259C`  | Tasks for Scheduling                       |
| ├─ Bitmap           | 895 B | `0x2000` - `0x237F`  | Bitmap with 16 B blocks                    |
| └─ Writeable Heap   | 14 KB | `0x259C` - `0x5D9C`  | Dynamic memory allocation (heap)           |
| **Interrupt-Table** | 611 B | `0x5D9D` - `0x....`  | Interrupt first jump to here               |
| **Stack**           | 8 KB  | `0x6000` – `0x7FFF`  | Stack for function calls, grows downward   |
| **Video RAM**       | 16 KB | `0x8000` – `0xBFFF`  | Framebuffer for visual output              |
| **Reserved**        | 8 KB  | `0xC000` – `0xDFFF`  | Reserved for I/O, buffers, MMIO            |
| └─ Keyboard I/O     | ~30 B | `0xC000` – `0xC020`  | Ring buffer, read/write pointers           |
| **Data-Section**    | 1KB   | `0xE000`  - `0xE3FF` | Initialized global / static data           |
| **Bss-Section**     | 1KB   | `0xE400`-`0xE7FF`    | Un / zero initialized data                 |
| **Extra / Future**  | 6 KB  | `0xE800` – `0xFFFF`  | Expansion, paging tables, filesystem, etc. |