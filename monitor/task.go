package monitor

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"sync"
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
	Output    bytes.Buffer
	outputMu  sync.Mutex
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
	// Create pipes for output
	stdout, err := t.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := t.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start command
	if err := t.Cmd.Start(); err != nil {
		return err
	}

	// Handle output in background
	go t.handleOutput(stdout)
	go t.handleOutput(stderr)

	// Wait for completion
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

func (t *Task) handleOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t.outputMu.Lock()
		t.Output.WriteString(scanner.Text() + "\n")
		t.outputMu.Unlock()
	}
}

func (t *Task) GetOutput() string {
	t.outputMu.Lock()
	defer t.outputMu.Unlock()
	return t.Output.String()
}

func (t *Task) Stop() {
	if t.Cmd != nil && t.Cmd.Process != nil {
		t.Cmd.Process.Kill()
	}
}
