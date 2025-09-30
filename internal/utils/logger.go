package utils

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Logger returns the logger instance
func Logger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		// Set log level
		logLevel := os.Getenv("LOG_LEVEL")
		switch logLevel {
		case "debug":
			logger.SetLevel(logrus.DebugLevel)
		case "info":
			logger.SetLevel(logrus.InfoLevel)
		case "warn":
			logger.SetLevel(logrus.WarnLevel)
		case "error":
			logger.SetLevel(logrus.ErrorLevel)
		default:
			logger.SetLevel(logrus.InfoLevel)
		}

		// Set formatter
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})

		// Set output
		logger.SetOutput(os.Stdout)
	}

	return logger
}

// WithFields creates a logger with fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return Logger().WithFields(fields)
}

// Debug logs a debug message
func Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Debug(msg)
	} else {
		Logger().Debug(msg)
	}
}

// Info logs an info message
func Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Info(msg)
	} else {
		Logger().Info(msg)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Warn(msg)
	} else {
		Logger().Warn(msg)
	}
}

// Error logs an error message
func Error(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Error(msg)
	} else {
		Logger().Error(msg)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Fatal(msg)
	} else {
		Logger().Fatal(msg)
	}
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Panic(msg)
	} else {
		Logger().Panic(msg)
	}
}

// Trace logs a trace message
func Trace(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		Logger().WithFields(fields[0]).Trace(msg)
	} else {
		Logger().Trace(msg)
	}
}
