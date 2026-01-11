// Package logging provides application logging with file and console output.
package logging

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"hytale-launcher/internal/build"
	"hytale-launcher/internal/hytale"
	"hytale-launcher/internal/ioutil"
)

const (
	// logFileName is the name of the log file.
	logFileName = "hytale-launcher.log"

	// maxLogFileSize is the maximum size of the log file before it is removed (10MB).
	maxLogFileSize = 10 * 1024 * 1024
)

var (
	// logFile is the current open log file.
	logFile *os.File

	// initOnce ensures Init is only called once.
	initOnce sync.Once
)

// Init initializes the logging system.
// It creates a log file in the hytale storage directory and configures
// both the standard logger and slog to write to both the file and stdout.
func Init() error {
	var initErr error

	initOnce.Do(func() {
		initErr = doInit()
	})

	return initErr
}

func doInit() error {
	// Get the log file path in the storage directory.
	logPath := hytale.InStorageDir(logFileName)
	logDir := filepath.Dir(logPath)

	// Ensure the storage directory exists.
	if err := ioutil.MkdirAll(logDir); err != nil {
		return fmt.Errorf("unable to create storage directory: %w", err)
	}

	// Check if the log file exists and is too large.
	if info, err := os.Stat(logPath); err == nil {
		if info.Size() > maxLogFileSize {
			if err := os.Remove(logPath); err != nil {
				return fmt.Errorf("unable to remove oversized log file %s: %w", logPath, err)
			}
		}
	}

	// Open the log file for appending.
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open log file %s: %w", logPath, err)
	}
	logFile = f

	// Create a multi-writer that writes to both the file and stdout.
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Configure the standard logger.
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Determine the log level based on environment.
	var logLevel slog.Level
	if build.DebugLogging() {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	// Configure slog with a text handler.
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: logLevel,
	})
	slog.SetDefault(slog.New(handler))

	return nil
}

// Close closes the log file.
// It should be called when the application exits.
func Close() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}
