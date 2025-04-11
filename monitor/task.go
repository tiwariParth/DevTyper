package monitor

import (
	"os"
	"os/exec"
	"strings"
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
	waitGroup sync.WaitGroup
}

func NewTask(command string, args ...string) *Task {
	// Create command with proper shell interpretation
	var cmd *exec.Cmd
	if strings.Contains(command, " ") || len(args) > 0 {
		// Use shell for complex commands
		shellCmd := append([]string{command}, args...)
		cmd = exec.Command("sh", "-c", strings.Join(shellCmd, " "))
	} else {
		// Simple command
		cmd = exec.Command(command, args...)
	}

	// Inherit parent's environment
	cmd.Env = os.Environ()

	// Setup stdio
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return &Task{
		Cmd:       cmd,
		StartTime: time.Now(),
		State:     TaskRunning,
		Done:      make(chan bool),
	}
}

func (t *Task) Start() error {
	// Prepare command but don't start yet
	if strings.Contains(t.Cmd.Path, " ") {
		shellCmd := []string{"-c", t.Cmd.Path + " " + strings.Join(t.Cmd.Args[1:], " ")}
		t.Cmd = exec.Command("sh", shellCmd...)
	}

	// Setup pipes
	t.Cmd.Stdout = os.Stdout
	t.Cmd.Stderr = os.Stderr
	t.Cmd.Stdin = os.Stdin

	// Start command after proper setup
	if err := t.Cmd.Start(); err != nil {
		return err
	}

	t.waitGroup.Add(1)
	go func() {
		defer t.waitGroup.Done()
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

func (t *Task) Stop() {
	if t.Cmd != nil && t.Cmd.Process != nil {
		t.Cmd.Process.Kill()
	}
	t.waitGroup.Wait()
}
