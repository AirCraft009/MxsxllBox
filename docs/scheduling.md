# Scheduler

## Technical Information

- 9 Task slots
- `_spawn`: creates task at addr(O1)
- `_yield`: willingly gives up control
- `_init_scheduler`: gives control to the scheduler only used once at the beginning
- All CPU flags are saved
- Registers R0 - R17, O1 - O6 and T1 are saved


### Yield Table

| Code | description                                             |
|------|---------------------------------------------------------|
| 0    | `running`(only internal)                                |
| 1    | `ready`(is ready)                                       |
| 2    | `keyboard-blocked`: (is set to ready after input)       |
| 3    | `timer/blocked`: (waiting on the next timer interrupt)  |
| 4    | `unused`:                                               |
| 5    | `unused`:                                               |
| 6    | `unused`:                                               |
| 7    | `terminated`: (the tasks isn't supposed to run anymore) |



## Details
> To avoid a single  action like redrawing the screen | reading the keyboard-buffer "hogging" all resources \
> the scheduler can save the current context/state of the cpu meaning(regs, PC, SP and flags) to memory.\
> To then give another Task the oportunity to continue. \
> This can occur with the help of `_yield`: which willingly gives up control,\
> a code can be moved into O2(see `Yield Table`) this code confirms if the task is IO,\
> for example Keyboard blocke(waiting on input) or smth else code 1 is ready,\
> meaning that it can be chosen again if no other task is currently available.\
> The other possibility for changing the current task is an `interrupt`.\
> Either IO | hardware timer these are forced and interrupt the program wherever it is at the moment.\
> A task is created by first moving its addr into O1 using `MOVA REG LBL`: and then calling `_spawn`\
> After creating all tasks control can be given to the scheduler by calling `_init_scheduler`.\
> The scheduler will start with the last task added and work it's way down before returning to the last.\
