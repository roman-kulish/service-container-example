package service

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"

	defaultLogFormat = LogFormatText
	defaultLogLevel  = zapcore.InfoLevel
)

// LogFormat represents a format of logs.
type LogFormat string

// LoggerOption represents a logger configuration option.
type LoggerOption interface {
	apply(*loggerConfig)
}

// loggerOptionFunc is a wrapper for configuration function, which satisfies
// LoggerOption interface.
type loggerOptionFunc func(*loggerConfig)

func (fn loggerOptionFunc) apply(cfg *loggerConfig) {
	fn(cfg)
}

// loggerConfig contains logger configuration.
type loggerConfig struct {
	format LogFormat
	level  zapcore.Level
	out    io.Writer
	opts   []zap.Option
}

// newLoggerConfig constructs logger configuration object with default configuration.
func newLoggerConfig() *loggerConfig {
	return &loggerConfig{
		format: defaultLogFormat,
		level:  defaultLogLevel,
	}
}

// WithLevel returns log level configuration option.
func WithLevel(l zapcore.Level) LoggerOption {
	return loggerOptionFunc(func(cfg *loggerConfig) {
		cfg.level = l
	})
}

// WithFormat returns log format configuration option.
func WithFormat(f LogFormat) LoggerOption {
	return loggerOptionFunc(func(cfg *loggerConfig) {
		cfg.format = f
	})
}

// WithOptions returns configuration option, which sets custom options
// to the logger.
//
// See: https://pkg.go.dev/go.uber.org/zap#Option
func WithOptions(opts ...zap.Option) LoggerOption {
	return loggerOptionFunc(func(cfg *loggerConfig) {
		cfg.opts = opts
	})
}

// WithOutput returns configuration option, which sets log output.
// Default is os.Stdout.
func WithOutput(w io.Writer) LoggerOption {
	return loggerOptionFunc(func(cfg *loggerConfig) {
		cfg.out = w
	})
}
