package logging

import (
	"fmt"
	"log"
	"os"
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
	Gray   = color.New(color.FgHiBlack)
)

var (
	logFile *os.File
	logger  *log.Logger
)

func Init() error {
	filename := fmt.Sprintf("logs/log_%s.txt", time.Now().Format("2006-01-02_15-04"))

	var err error
	logFile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	return nil
}

func Log(c *color.Color, msg string, a ...any) {
	c.Printf(msg+"\n", a...)
}

func LogLine(c *color.Color, a ...any) {
	c.Println(a...)
}

func LogInfo(msg string, a ...any) {
	LogWithTime(Green, "INFO", msg, a...)
}

func LogInfoC(c *color.Color, msg string, a ...any) {
	LogWithTime(c, "INFO", msg, a...)
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

	if logger != nil {
		plainLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, formatted)
		logger.Print(plainLine)
	}
}
