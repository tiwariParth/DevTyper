package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/parth/DevTyper/game"
	"github.com/parth/DevTyper/monitor"
)

func main() {
	forceExit := flag.Bool("force-exit", false, "Exit game immediately when task completes")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: devtyper [-force-exit] <command>")
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

	// Start task after user input
	if err := task.Start(); err != nil {
		fmt.Printf("Error starting task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting long-running task: %s\n", description)

	if strings.ToLower(response) != "n" {
		g, err := game.New(task.Done, description)
		if err != nil {
			fmt.Printf("Error starting game: %v\n", err)
			task.Stop() // Stop task if game fails
			os.Exit(1)
		}
		g.ForceExit = *forceExit
		g.Run()
	}

	// Wait for task completion
	select {
	case <-task.Done:
		fmt.Println("Task completed!")
	case <-time.After(time.Second * 1): // Add timeout for cleanup
		task.Stop()
	}
}
