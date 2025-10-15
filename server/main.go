package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var config = make(map[string]string)

var cyan = color.New(color.FgCyan)
var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)
var yellow = color.New(color.FgYellow)

func main() {
	cyan.Println("TetriON WebServer")
	green.Println("Starting server...")
	LoadEnv()
	LoadConfig()
	InitRedis()
}

func LoadEnv() {
	yellow.Println("Loading .env file...")
	var file, err = os.OpenFile("..\\.env", os.O_RDONLY, 0644)
	if err != nil {
		red.Println("Error loading .env file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	if err := scanner.Err(); err != nil {
		red.Println("Error reading .env file:", err)
		return
	}
	green.Println(".env file loaded successfully.")
}

func LoadConfig() {
	yellow.Println("Loading config.json file...")
	var file, err = os.OpenFile("..\\config.json", os.O_RDONLY, 0644)
	if err != nil {
		red.Println("Error loading config.json file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding config.json file:", err)
		return
	}
	green.Println("config.json file loaded successfully.")
	green.Println("Config:", config)
}

