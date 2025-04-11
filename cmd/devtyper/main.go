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
	_, description := monitor.DetectCommand(cmdString)

	// Create task but don't start yet
	task := monitor.NewTask(args[0], args[1:]...)

	// Get user input before starting task
	fmt.Println("Want to practice typing while waiting? [Y/n]")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start task after user input
	if err := task.Start(); err != nil {
		fmt.Printf("Error starting task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting long-running task: %s\n", description)

	if strings.ToLower(response) != "n" {
		g, err := game.New(task.Done, description, task)
		if err != nil {
			fmt.Printf("Error starting game: %v\n", err)
			task.Stop() // Stop task if game fails
			os.Exit(1)
		}

		// Handle Ctrl+C while game is running
		go func() {
			<-sigChan
			g.Cleanup()
			task.Stop()
			os.Exit(0)
		}()

		g.ForceExit = *forceExit
		g.Run()

		// After game exits, check if we should wait for task
		if *keepAlive {
			fmt.Println("\nGame exited. Waiting for command to complete...")
			fmt.Println("Press Ctrl+C to force quit")
			select {
			case <-task.Done:
				fmt.Println("Command completed!")
			case <-sigChan:
				fmt.Println("\nForce quitting...")
				task.Stop()
			}
		} else {
			task.Stop()
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
