package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

func Listen() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	stopChan := make(chan struct{})

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "stop", "exit", "quit":
				fmt.Println("Stopping server...")
				close(stopChan)
				return
			case "status":
				fmt.Println("Server is currently running.")
			case "":
				continue
			default:
				fmt.Printf("Unknown command: %s\n", input)
			}
		}
	}()

	<-stopChan
	fmt.Println("Server has been stopped.")
}
