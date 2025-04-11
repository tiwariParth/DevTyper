package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/parth/DevTyper/game"
	"github.com/parth/DevTyper/monitor"
)

// Add TaskContext struct to hold shared channels
type TaskContext struct {
	doneChan chan struct{}
	sigChan  chan os.Signal
	task     *monitor.Task
}

// Handle task completion based on user preference
func handleTask(ctx *TaskContext, keepAlive bool) {
	select {
	case <-ctx.task.Done:
		<-ctx.doneChan // Wait for output to finish
		if ctx.task.HasError() {
			fmt.Printf("\nTask failed: %s\n", ctx.task.GetError())
		} else {
			fmt.Println("\nTask completed successfully!")
		}
		if !keepAlive {
			ctx.task.Stop()
		}
	case <-ctx.sigChan:
		fmt.Println("\nStopping task...")
		ctx.task.Stop()
	}
}

func main() {
	forceExit := flag.Bool("force-exit", false, "Exit game immediately when task completes")
	keepAlive := flag.Bool("keep-alive", true, "Keep command running after exiting game")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: devtyper [-force-exit] [-keep-alive] <command>")
		os.Exit(1)
	}

	// Join all args to detect command type
	cmdString := strings.Join(args, " ")
	_, description, isInteractive, argExample := monitor.DetectCommand(cmdString)

	if isInteractive {
		fmt.Println("\nThis command requires interactive input.")
		fmt.Println("To skip interactive mode, try using arguments instead:")
		fmt.Printf("\n  %s\n\n", argExample)
		fmt.Println("Exiting. Please retry with arguments.")
		os.Exit(0)
	}

	// Setup signal handling
	ctx := &TaskContext{
		doneChan: make(chan struct{}),
		sigChan:  make(chan os.Signal, 1),
		task:     monitor.NewTask(args[0], args[1:]...),
	}
	signal.Notify(ctx.sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals for clean shutdown
	go func() {
		<-ctx.sigChan
		fmt.Print("\n") // New line after ^C
		ctx.task.Stop()
		fmt.Print("\033[?25h") // Show cursor
		os.Exit(0)
	}()

	// Get user input before starting task
	fmt.Println("Want to practice typing while waiting? [Y/n]")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	// Start output display goroutine
	go func() {
		for output := range ctx.task.GetOutputChannel() {
			fmt.Print(output)
		}
		ctx.doneChan <- struct{}{}
	}()

	// Start task after user input
	if err := ctx.task.Start(); err != nil {
		fmt.Printf("\nError starting task: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	fmt.Printf("\nStarting task: %s\n", description)

	if strings.ToLower(response) != "n" {
		g, err := game.New(ctx.task.Done, description, ctx.task)
		if err != nil {
			fmt.Printf("\nError starting game: %v\n", err)
			ctx.task.Stop()
			cleanup()
			os.Exit(1)
		}

		// Run game
		g.ForceExit = *forceExit
		g.Run()

		// Show status after game exits
		if ctx.task.IsComplete() {
			fmt.Println("\nTask completed while playing!")
			<-ctx.doneChan // Wait for output to finish
			cleanup()
			os.Exit(0)
		}

		// Wait for task if it's still running
		fmt.Println("\nTask is still running. Press Ctrl+C to stop.")
		handleTask(ctx, *keepAlive)
	}

	cleanup()
}

func cleanup() {
	fmt.Print("\033[?25h") // Show cursor
	fmt.Print("\033[2J\033[H") // Clear screen
}
