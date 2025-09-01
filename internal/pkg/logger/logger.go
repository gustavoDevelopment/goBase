// Package logger provides logging functionality for the application
package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Context key for storing/retrieving logger from context
const LoggerKey = "logger"

var (
	// Log is the global logger instance
	Log *zap.Logger
)

// InitLogger initializes the application logger
func InitLogger(debug bool) error {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	// Common configuration
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.ConsoleSeparator = " "

	// Create a new core with console encoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// Create the logger with our custom core
	Log = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	// Redirect standard logs
	zap.RedirectStdLog(Log)

	return nil
}

// Logger returns a new logger with the request ID field
func Logger() *zap.Logger {
	return Log
}

// WithRequestID adds a request ID to the logger context
func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}

// FromContext retrieves a logger from the context, or returns a new one if not found
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Log
	}
	if logger, ok := ctx.Value("logger").(*zap.Logger); ok && logger != nil {
		return logger
	}
	return Log
}

// Sync closes the log buffer
func Sync() error {
	return Log.Sync()
}
