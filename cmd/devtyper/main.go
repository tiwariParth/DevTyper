package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

	cmd := strings.Join(args, " ")
	_, description := monitor.DetectCommand(cmd)

	task := monitor.NewTask(args[0], args[1:]...)
	if err := task.Start(); err != nil {
		fmt.Printf("Error starting task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting long-running task: %s\n", description)
	fmt.Println("Want to practice typing while waiting? [Y/n]")

	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "n" {
		g, err := game.New(task.Done, description)
		if err != nil {
			fmt.Printf("Error starting game: %v\n", err)
			os.Exit(1)
		}
		g.ForceExit = *forceExit
		g.Run()
	}

	<-task.Done
}
