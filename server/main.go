package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var config = make(map[string]string)

func main() {
	fmt.Println("Starting server...")
	LoadEnv()
	LoadConfig()
	//InitRedis()
}

func LoadEnv() {
	fmt.Println("Loading .env file...")
	var file, err = os.OpenFile("..\\.env", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error loading .env file:", err)
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
		fmt.Println("Error reading .env file:", err)
		return
	}
	fmt.Println(".env file loaded successfully.")
}

func LoadConfig() {
	fmt.Println("Loading config.json file...")
	var file, err = os.OpenFile("..\\config.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error loading config.json file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding config.json file:", err)
		return
	}
	fmt.Println("config.json file loaded successfully.")
	fmt.Println("Config:", config)
}

