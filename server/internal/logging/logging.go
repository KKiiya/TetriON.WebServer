package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	Cyan   = color.New(color.FgCyan).Add(color.Bold)
	Green  = color.New(color.FgGreen).Add(color.Bold)
	Yellow = color.New(color.FgYellow).Add(color.Bold)
	Red    = color.New(color.FgRed).Add(color.Bold)
	White  = color.New(color.FgWhite)
)

func LogLine(c *color.Color, a ...any) {
	c.Println(a...)
}

func LogInfo(msg string, a ...any) {
	LogWithTime(Green, "INFO", msg, a...)
}

func LogError(msg string, a ...any) {
	LogWithTime(Red, "ERROR", msg, a...)
}

func LogWarning(msg string, a ...any) {
	LogWithTime(Yellow, "WARNING", msg, a...)
}

func LogDebug(msg string, a ...any) {
	LogWithTime(White, "DEBUG", msg, a...)
}

// LogWithTime prints a timestamped message with the given color.
func LogWithTime(c *color.Color, level string, msg string, a ...any) {
	timestamp := time.Now().Format("15:04:05")

	// If there are arguments but the msg contains no formatting verbs, append a `%v`
	if len(a) > 0 && !strings.ContainsAny(msg, "%") {
		msg = msg + " %v"
	}

	formatted := fmt.Sprintf(msg, a...)
	// Print: [HH:MM:SS][LEVEL] EMOJI formatted-message
	c.Printf("[%s] [%s] %s\n", timestamp, level, formatted)
}