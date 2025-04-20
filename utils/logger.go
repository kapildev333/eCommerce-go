package utils

import (
	"log/slog"
	"os"
	"sync"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

// InitLogger initializes the global logger instance
func InitLogger() *slog.Logger {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler := slog.NewJSONHandler(os.Stdout, opts)
		globalLogger = slog.New(handler)
	})
	return globalLogger
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	if globalLogger == nil {
		return InitLogger()
	}
	return globalLogger
}
