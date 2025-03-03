package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes a structured JSON logger
func InitLogger() {
	// Create a JSON-based logger
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Example log entry to indicate logger is initialized
	Log.Info("Logger initialized")
}
