# Design Decisions


## Memory

- Flat memory, programs can interface directly with hardware, but abstractions are provided
- Dynamic heap allocation is provided but not enforced
- Heap allocator uses first-fit policy for simplicity and speed
- Details to [Memory Management](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/memory-management.md)
- Details to the [Heap] (https://github.com/AirCraft009/MxsxllBox/blob/master/docs/abi.md#chapter-32---interacting-with-the-heap)

- Stack handles word-sized(16B) data
- Details to the [Stack] (https://github.com/AirCraft009/MxsxllBox/blob/master/docs/abi.md#chapter-31---interacting-with-the-stack)

- Frame Buffer is connected to the visual output
- Each pixel is scaled up x4
- 4 colors are avaliable
- Details to the [Frame Buffer](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/abi.md#chapter-33---interacting-with-the-framebuffer)

- Labels aren't assembled into byte code (only if debug is enabled)
- The pre-Assembled code is outputted in a .obj file
- The .obj file is then linked together with other .obj files into a  final .bin file


