package monitor

import (
	"os/exec"
	"time"
)

type TaskState int

const (
	TaskRunning TaskState = iota
	TaskCompleted
	TaskFailed
)

type Task struct {
	Cmd       *exec.Cmd
	StartTime time.Time
	State     TaskState
	Done      chan bool
}

func NewTask(command string, args ...string) *Task {
	cmd := exec.Command(command, args...)
	return &Task{
		Cmd:       cmd,
		StartTime: time.Now(),
		State:     TaskRunning,
		Done:      make(chan bool),
	}
}

func (t *Task) Start() error {
	if err := t.Cmd.Start(); err != nil {
		return err
	}

	go func() {
		t.Cmd.Wait()
		if t.Cmd.ProcessState.Success() {
			t.State = TaskCompleted
		} else {
			t.State = TaskFailed
		}
		t.Done <- true
	}()

	return nil
}
