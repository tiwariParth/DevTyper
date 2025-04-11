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
	_, description, isInteractive := monitor.DetectCommand(cmdString)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if isInteractive {
		fmt.Println("This command requires user input. Please provide the required information:")
		task := monitor.NewTask(args[0], args[1:]...)
		if err := task.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Wait for command to complete without game
		select {
		case <-task.Done:
			fmt.Println("\nCommand completed!")
		case <-sigChan:
			fmt.Println("\nForce quitting...")
			task.Stop()
		}
		return
	}

	// Create task but don't start yet
	task := monitor.NewTask(args[0], args[1:]...)

	// Handle signals for clean shutdown
	go func() {
		<-sigChan
		fmt.Print("\n") // New line after ^C
		task.Stop()
		fmt.Print("\033[?25h") // Show cursor
		os.Exit(0)
	}()

	// Get user input before starting task
	fmt.Println("Want to practice typing while waiting? [Y/n]")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	// Start task after user input
	if err := task.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		// Reset terminal state
		fmt.Print("\033[?25h") // Show cursor
		fmt.Print("\033[2J\033[H") // Clear screen
		os.Exit(1)
	}

	// Setup error recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Print("\033[?25h") // Show cursor
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Printf("Error: %v\n", r)
		}
	}()

	fmt.Printf("Starting long-running task: %s\n", description)

	if strings.ToLower(response) != "n" {
		g, err := game.New(task.Done, description, task)
		if err != nil {
			fmt.Printf("Error starting game: %v\n", err)
			task.Stop() // Stop task if game fails
			os.Exit(1)
		}

		g.ForceExit = *forceExit
		g.Run()

		if *keepAlive {
			fmt.Println("\nGame exited. Command is still running...")
			select {
			case <-task.Done:
				fmt.Println("\nCommand completed!")
			case <-sigChan:
				fmt.Println("\nStopping command...")
				task.Stop()
			}
		} else {
			task.Stop()
			fmt.Println("\nCommand stopped.")
		}
	} else {
		// Wait for task if not playing
		select {
		case <-task.Done:
			fmt.Println("Command completed!")
		case <-sigChan:
			fmt.Println("\nForce quitting...")
			task.Stop()
		}
	}
}
