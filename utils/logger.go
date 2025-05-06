package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// Log levels
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	// Default logger
	logger *Logger

	// Log level names
	levelNames = map[int]string{
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
	}
)

// Logger represents a custom logger
type Logger struct {
	level  int
	output io.Writer
}

// InitLogger initializes the default logger
func InitLogger(level int, logToFile bool) error {
	// Set default output to stdout
	output := io.Writer(os.Stdout)

	// If logging to file is enabled
	if logToFile {
		// Create logs directory if it doesn't exist
		if err := os.MkdirAll("logs", 0755); err != nil {
			return err
		}

		// Create log file with current date
		logFileName := fmt.Sprintf("logs/app_%s.log", time.Now().Format("2006-01-02"))
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		// Use both stdout and file for output
		output = io.MultiWriter(os.Stdout, logFile)
	}

	// Create default logger
	logger = &Logger{
		level:  level,
		output: output,
	}

	// Set standard log output to our logger
	log.SetOutput(output)

	return nil
}

// log formats and writes a log message
func (l *Logger) log(level int, format string, args ...interface{}) {
	// Skip if log level is lower than configured level
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	fileInfo := "???"
	if ok {
		fileInfo = filepath.Base(file)
	}

	// Format timestamp
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	// Format log message
	levelName := levelNames[level]
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] [%s] %s:%d - %s\n", timestamp, levelName, fileInfo, line, message)

	// Write to output
	_, _ = fmt.Fprint(l.output, logMessage)

	// If fatal, exit the program
	if level == LevelFatal {
		os.Exit(1)
	}
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	if logger == nil {
		InitLogger(LevelDebug, false)
	}
	logger.log(LevelDebug, format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	if logger == nil {
		InitLogger(LevelDebug, false)
	}
	logger.log(LevelInfo, format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	if logger == nil {
		InitLogger(LevelDebug, false)
	}
	logger.log(LevelWarn, format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	if logger == nil {
		InitLogger(LevelDebug, false)
	}
	logger.log(LevelError, format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(format string, args ...interface{}) {
	if logger == nil {
		InitLogger(LevelDebug, false)
	}
	logger.log(LevelFatal, format, args...)
}
