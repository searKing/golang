package resillience

//go:generate stringer -type TaskState -trimprefix=TaskState
//go:generate jsonenums -type TaskState

type TaskState int

const (
	TaskStateNew               TaskState = iota // Task state for a task which has not yet started.
	TaskStateRunning                            // Task state for a running task. A task in the running state is executing in the Go routine but it may be waiting for other resources from the operating system such as processor.
	TaskStateDoneErrorHappened                  // Task state for a terminated state. The task has completed execution with some errors happened
	TaskStateDoneNormally                       // Task state for a terminated state. The task has completed execution normally
	TaskStateDormancy                           // Task state for a terminated state. The task has completed execution normally and will be started if New's called
	TaskStateDeath                              // Task state for a terminated state. The task has completed execution normally and will never be started again
	TaskStateButt
)
