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
	Cmd        *exec.Cmd
	StartTime  time.Time
	State      TaskState
	Done       chan bool
	Output     bytes.Buffer
	outputMu   sync.Mutex
	ErrorChan  chan error
	err        error
	errMu      sync.Mutex
	pty        *os.File
	isComplete bool
	statusMu   sync.Mutex
	output     chan string
}

func NewTask(command string, args ...string) *Task {
	cmd := exec.Command(command, args...)
	return &Task{
		Cmd:        cmd,
		StartTime:  time.Now(),
		State:      TaskRunning,
		Done:       make(chan bool),
		ErrorChan:  make(chan error, 1),
		output:     make(chan string, 100),
	}
}

func (t *Task) Start() error {
	var err error
	t.pty, err = pty.Start(t.Cmd)
	if err != nil {
		return err
	}

	// Handle output in background with better buffer management
	go func() {
		buf := make([]byte, 1024) // Smaller buffer for more frequent updates
		for {
			n, err := t.pty.Read(buf)
			if err != nil {
				if !t.isComplete {
					t.setError(err)
				}
				break
			}

			output := string(buf[:n])
			t.outputMu.Lock()
			// Limit buffer size to prevent memory issues
			if t.Output.Len() > 100*1024 {
				// Clear half the buffer if it gets too large
				oldContent := t.Output.String()
				t.Output.Reset()
				t.Output.WriteString(oldContent[len(oldContent)/2:])
			}
			t.Output.WriteString(output)
			t.outputMu.Unlock()

			// Send to output channel with non-blocking write
			select {
			case t.output <- output:
			default:
				// Skip if channel is full
			}
		}
	}()

	// Wait for completion
	go func() {
		err := t.Cmd.Wait()
		t.statusMu.Lock()
		t.isComplete = true
		t.statusMu.Unlock()

		if err != nil {
			t.setError(err)
			t.State = TaskFailed
		} else {
			t.State = TaskCompleted
		}
		close(t.output)
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

func (t *Task) IsComplete() bool {
	t.statusMu.Lock()
	defer t.statusMu.Unlock()
	return t.isComplete
}

func (t *Task) GetOutputChannel() <-chan string {
	return t.output
}
