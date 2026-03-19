package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"TetriON.WebServer/server/internal/logging"
)

var (
	commands       = make(map[string]*Command)
	commandAliases = make(map[string]string)
	stopChan       = make(chan struct{})
	stopOnce       sync.Once
)

type Command struct {
	name        string
	description string
	usage       string
	aliases     []string
	run         func(arguments ...string)
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
			if input == "" {
				continue
			}

			parts := strings.Fields(strings.ToLower(input))
			if len(parts) == 0 {
				continue
			}

			commandName := parts[0]
			arguments := parts[1:]
			if canonical, exists := commandAliases[commandName]; exists {
				commandName = canonical
			}

			if command, exists := commands[commandName]; exists {
				if command.enabled {
					command.run(arguments...)
				} else {
					logging.White.Printf("Command '%s' is disabled.\n", commandName)
				}
			} else {
				logging.White.Printf("Unknown command: %s. Type 'help' to see available commands\n", commandName)
			}
		}
	}()

	<-stopChan
	logging.White.Println("Server has been stopped.")
}

func RegisterCommand(command string, alias []string, description string, usage string, action func(arguments ...string)) *Command {
	command = strings.ToLower(command)
	if _, exists := commands[command]; exists {
		logging.LogDebug("Command '%s' already exists", command)
		return commands[command]
	}
	cmd := &Command{
		name:        command,
		aliases:     alias,
		description: description,
		usage:       usage,
		run:         action,
		enabled:     true,
	}
	commands[command] = cmd
	for _, a := range alias {
		commandAliases[strings.ToLower(a)] = command
	}
	return cmd
}

func ToggleCommand(command string) {
	command = strings.ToLower(command)
	if canonical, exists := commandAliases[command]; exists {
		command = canonical
	}

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
		logging.LogDebug("Command '%s' doesn't exist", command)
	}
}

// PRIVATE FUNCTION
func registerDefaultCommands() {
	RegisterCommand("help", []string{"h", "?"}, "Show every command in console", "help", func(arguments ...string) {
		logging.Gray.Println("------------------------------------------------------")
		logging.White.Println("Available commands:")
		for name, command := range commands {
			logging.White.Printf("%s: %s\n", name, command.description)
		}
		logging.Gray.Println("------------------------------------------------------")
	})

	RegisterCommand("stop", []string{"quit", "exit", "shut", "kill"}, "Stop the WebServer", "stop", func(arguments ...string) {
		logging.White.Println("Stopping server...")
		stopOnce.Do(func() {
			close(stopChan)
		})
	})

	RegisterCommand("status", []string{}, "Show the current server status", "status", func(arguments ...string) {
		logging.White.Println("Server is currently running.")
	})
}
