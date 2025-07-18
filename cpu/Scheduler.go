package cpu

type TaskState byte

const (
	ready TaskState = iota
	blocked
	keyboardBlocked
	timerBlocked
	terminated
	running
)

type Task struct {
	ID        int
	PC        uint16
	SP        uint16
	Registers [NumRegisters]uint16
	State     TaskState
}

func (task *Task) SaveTask(TaskExitReason TaskState, cpu *CPU) {
	cpu.Tasks = RemoveValueAtIndex(cpu.Tasks, int(cpu.ActiveTask))
	if TaskExitReason == terminated {
		return
	}
	task.State = TaskExitReason
	task.Registers = cpu.Registers
	cpu.Tasks = append(cpu.Tasks, task)
}

func CreateNewTask(cpu *CPU, id int, state TaskState) *Task {
	return &Task{
		ID:        id,
		PC:        cpu.PC,
		SP:        0,
		Registers: cpu.Registers,
		State:     state,
	}
}

func (cpu *CPU) ReturnToTask(task *Task) {
	task.State = running
	cpu.PC = task.PC
	if task.SP != 0 {
		cpu.SP = task.SP
	} else {
		task.SP = cpu.SP + instructionSizeShort
	}
	cpu.Registers = task.Registers
}

func RemoveValueAtIndex(input []*Task, index int) []*Task {
	inputBeforeIndex := input[:index]
	if index == len(input)-1 {
		return inputBeforeIndex
	}
	return append(inputBeforeIndex, input[index+1:]...)

}
