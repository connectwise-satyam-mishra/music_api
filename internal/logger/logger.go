package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

// Init initializes the logger by creating a log file and setting up the logger instance.
func Init() {
	logFile, log_err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if log_err != nil {
		log.Fatal("Failed to create log file:", log_err)
	}

	// Initialize the logger
	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

// Info logs an informational message.
func Info(message string) {
	Logger.Println(message)
}

// Error logs an error message and exits the program.
func Error(message string, err error) {
	Logger.Fatal(message, err)
}

// Warning logs a warning message.
func Warning(message string) {
	Logger.Println("warning ", message)
}
