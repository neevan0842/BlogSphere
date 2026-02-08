package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

// Init initializes the global logger
func Init() {
	logger := zap.Must(zap.NewProduction())
	log = logger.Sugar()
}

// Get returns the global logger instance
func Get() *zap.SugaredLogger {
	if log == nil {
		Init()
	}
	return log
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
