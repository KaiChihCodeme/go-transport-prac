package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger with additional functionality
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// Config holds logger configuration
type Config struct {
	Level       string `json:"level"`
	Format      string `json:"format"`
	OutputPaths string `json:"output_paths"`
	Development bool   `json:"development"`
}

// New creates a new logger with the given configuration
func New(cfg Config) (*Logger, error) {
	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(parseLevel(cfg.Level)),
		Development: cfg.Development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: cfg.Format,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      parseOutputPaths(cfg.OutputPaths),
		ErrorOutputPaths: []string{"stderr"},
	}

	// Adjust encoder config for development
	if cfg.Development {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	// Build the logger
	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

// NewDevelopment creates a development logger with sensible defaults
func NewDevelopment() (*Logger, error) {
	return New(Config{
		Level:       "debug",
		Format:      "console",
		OutputPaths: "stdout",
		Development: true,
	})
}

// NewProduction creates a production logger with sensible defaults
func NewProduction() (*Logger, error) {
	return New(Config{
		Level:       "info",
		Format:      "json",
		OutputPaths: "stdout",
		Development: false,
	})
}

// Sugar returns the sugared logger for easier logging
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		sugar:  l.Logger.With(fields...).Sugar(),
	}
}

// WithComponent adds a component field to the logger
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithFields(zap.String("component", component))
}

// WithRequestID adds a request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

// WithUserID adds a user ID field to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithFields(zap.String("user_id", userID))
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}

// LogHTTPRequest logs HTTP request information
func (l *Logger) LogHTTPRequest(method, path string, statusCode int, duration string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.String("duration", duration),
	}, fields...)

	l.Info("HTTP request", allFields...)
}

// LogGRPCRequest logs gRPC request information
func (l *Logger) LogGRPCRequest(method string, statusCode int, duration string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("grpc_method", method),
		zap.Int("status_code", statusCode),
		zap.String("duration", duration),
	}, fields...)

	l.Info("gRPC request", allFields...)
}

// LogDatabaseQuery logs database query information
func (l *Logger) LogDatabaseQuery(query string, duration string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("query", query),
		zap.String("duration", duration),
	}, fields...)

	l.Debug("Database query", allFields...)
}

// LogCacheOperation logs cache operation information
func (l *Logger) LogCacheOperation(operation, key string, hit bool, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("cache_operation", operation),
		zap.String("cache_key", key),
		zap.Bool("cache_hit", hit),
	}, fields...)

	l.Debug("Cache operation", allFields...)
}

// parseLevel parses string level to zapcore.Level
func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// parseOutputPaths parses output paths string to slice
func parseOutputPaths(paths string) []string {
	if paths == "" {
		return []string{"stdout"}
	}

	// Split by comma and trim whitespace
	pathList := strings.Split(paths, ",")
	for i, path := range pathList {
		pathList[i] = strings.TrimSpace(path)
	}

	return pathList
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Close closes the logger and flushes any buffered entries
func (l *Logger) Close() error {
	return l.Sync()
}

// Global logger instance
var globalLogger *Logger

// InitGlobal initializes the global logger
func InitGlobal(cfg Config) error {
	logger, err := New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize global logger: %w", err)
	}

	globalLogger = logger
	return nil
}

// InitGlobalDevelopment initializes the global logger for development
func InitGlobalDevelopment() error {
	logger, err := NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to initialize global development logger: %w", err)
	}

	globalLogger = logger
	return nil
}

// InitGlobalProduction initializes the global logger for production
func InitGlobalProduction() error {
	logger, err := NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize global production logger: %w", err)
	}

	globalLogger = logger
	return nil
}

// Global returns the global logger instance
func Global() *Logger {
	if globalLogger == nil {
		// Fallback to development logger if not initialized
		logger, err := NewDevelopment()
		if err != nil {
			// Last resort: create a basic logger
			zapLogger, _ := zap.NewDevelopment()
			globalLogger = &Logger{
				Logger: zapLogger,
				sugar:  zapLogger.Sugar(),
			}
		} else {
			globalLogger = logger
		}
	}
	return globalLogger
}

// Convenience functions using global logger
func Debug(msg string, fields ...zap.Field) {
	Global().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Global().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Global().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Global().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Global().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Global().Panic(msg, fields...)
}

// Sugar convenience functions
func Debugf(template string, args ...interface{}) {
	Global().Sugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Global().Sugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Global().Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Global().Sugar().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Global().Sugar().Fatalf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	Global().Sugar().Panicf(template, args...)
}