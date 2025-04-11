package monitor

import (
	"bytes"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
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
	ErrorChan chan error
	err       error
	errMu     sync.Mutex
	pty       *os.File
}

func NewTask(command string, args ...string) *Task {
	cmd := exec.Command(command, args...)
	return &Task{
		Cmd:       cmd,
		StartTime: time.Now(),
		State:     TaskRunning,
		Done:      make(chan bool),
		ErrorChan: make(chan error, 1),
	}
}

func (t *Task) Start() error {
	var err error
	t.pty, err = pty.Start(t.Cmd)
	if err != nil {
		return err
	}

	// Handle output in background
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := t.pty.Read(buf)
			if err != nil {
				break
			}
			t.outputMu.Lock()
			t.Output.Write(buf[:n])
			t.outputMu.Unlock()
		}
	}()

	// Wait for completion
	go func() {
		err := t.Cmd.Wait()
		if err != nil {
			t.setError(err)
			t.ErrorChan <- err
			t.State = TaskFailed
		} else {
			t.State = TaskCompleted
		}
		t.pty.Close()
		t.Done <- true
	}()

	return nil
}

func (t *Task) Stop() {
	if t.Cmd != nil && t.Cmd.Process != nil {
		t.Cmd.Process.Signal(syscall.SIGTERM)
		time.Sleep(100 * time.Millisecond)
		t.Cmd.Process.Kill()
	}
	if t.pty != nil {
		t.pty.Close()
	}
}

func (t *Task) GetOutput() string {
	t.outputMu.Lock()
	defer t.outputMu.Unlock()
	return t.Output.String()
}

func (t *Task) setError(err error) {
	t.errMu.Lock()
	t.err = err
	t.errMu.Unlock()
}

func (t *Task) GetError() string {
	t.errMu.Lock()
	defer t.errMu.Unlock()
	if t.err != nil {
		return t.err.Error()
	}
	return ""
}

func (t *Task) HasError() bool {
	t.errMu.Lock()
	defer t.errMu.Unlock()
	return t.err != nil
}
