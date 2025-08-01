package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Init sets up the logger, creates a new log file,
// directs output to both console and file, and
// returns an error if initialization fails.
func Init() error {
	now := time.Now()
	yearMonth := now.Format("2006-01")
	day := now.Format("02")
	logsDir := filepath.Join(".", "logs", yearMonth, day)

	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		Errorf(`Failed to create logs directory "%s". Error: %v`, logsDir, err)

		return err
	}

	// Generate unique log file (./logs/2025-06/18/2025-06-18-12-58-54.938.log).
	filename := now.Format("2006-01-02-15-04-05.000") + ".log"
	logFilePath := filepath.Join(logsDir, filename)

	// Open log file.
	const permissions = 0644

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, permissions)
	if err != nil {
		Errorf(`Failed to open log file "%s". Error: %v`, logFilePath, err)

		return err
	}

	// Set log output to console and to file.
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	return nil
}

// NewLine logs a blank line without timestamp or prefix.
func NewLine() {
	oldFlags := log.Flags()

	log.SetFlags(0)
	log.Println()
	log.SetFlags(oldFlags)
}

// Fatalf logs a formatted fatal message with context.
func Fatalf(format string, args ...interface{}) {
	log.Printf("[FATAL] %s %s", fmt.Sprintf(format, args...), occurrence())
}

// Errorf logs a formatted error message with context.
func Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] %s %s", fmt.Sprintf(format, args...), occurrence())
}

// Warnf logs a formatted warning message with context.
func Warnf(format string, args ...interface{}) {
	log.Printf("[WARN]  %s %s", fmt.Sprintf(format, args...), occurrence())
}

// Infof logs a formatted informational message with context.
func Infof(format string, args ...interface{}) {
	log.Printf("[INFO]  %s %s", fmt.Sprintf(format, args...), occurrence())
}

// Debugf logs a formatted debug message with context.
func Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] %s %s", fmt.Sprintf(format, args...), occurrence())
}

// occurrence retrieves the caller’s file and line number
// for log entries.
func occurrence() string {
	const skip = 2

	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	return fmt.Sprintf("(%s:%d)\n", filepath.Base(file), line)
}
