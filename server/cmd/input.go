package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"TetriON.WebServer/server/internal/logging"
)

var (
	commands = make(map[string]*Command)

	stopChan = make(chan struct{})
)

type Command struct {
	name        string
	description string
	usage       string
	aliases     []string
	run         func(arguments ...any)
	enabled     bool
}

func Listen() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	registerDefaultCommands()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			arguments := strings.Fields(input)[1:]

			if command, exists := commands[input]; exists && command.enabled {
				command.run(arguments)
			} else {
				fmt.Printf("Unknown command: %s. Type 'help' to see available commands\n", input)
			}
		}
	}()

	<-stopChan
	fmt.Println("Server has been stopped.")
}

func RegisterCommand(command string, alias []string, description string, usage string, action func(arguments ...any)) *Command {
	cmd := &Command{
		name:        command,
		aliases:     alias,
		description: description,
		usage:       usage,
		run:         action,
		enabled:     true,
	}
	commands[command] = cmd
	return cmd
}

func ToggleCommand(command string) {
	if command, exists := commands[command]; exists {
		command.enabled = !command.enabled
		var status string
		if command.enabled {
			status = "enabled"
		} else {
			status = "disabled"
		}
		logging.LogDebug("Command '%s' %s", command, status)
	} else {
		logging.LogDebug("Command '%s' doesn't exist")
	}
}

// PRIVATE FUNCTION
func registerDefaultCommands() {
	RegisterCommand("help", []string{"h", "?"}, "Show every command in console", "help", func(arguments ...any) {

	})

	RegisterCommand("stop", []string{"quit", "exit", "shut", "kill"}, "Stop the WebServer", "stop", func(arguments ...any) {
		fmt.Println("Stopping server...")
		close(stopChan)
	})

	RegisterCommand("status", []string{}, "Show the current server status", "status", func(arguments ...any) {
		fmt.Println("Server is currently running.")
	})
}
