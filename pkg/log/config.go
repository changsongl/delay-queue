package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level log level
type Level int8

const (
	// LevelDebug debug
	LevelDebug Level = iota
	// LevelInfo info
	LevelInfo
	// LevelWarn warn
	LevelWarn
	// LevelError error
	LevelError
	// LevelFatal fatal
	LevelFatal
)

// levelMap level map zap level
var levelMap = map[Level]zapcore.Level{
	LevelDebug: zap.DebugLevel,
	LevelInfo:  zapcore.InfoLevel,
	LevelWarn:  zapcore.WarnLevel,
	LevelError: zapcore.ErrorLevel,
	LevelFatal: zapcore.FatalLevel,
}

// Config object
type Config struct {
	conf      zap.Config
	zapOption []zap.Option
}

// NewConfig create new log configuration
func NewConfig() *Config {
	c := &Config{conf: zap.NewProductionConfig()}
	c.conf.Encoding = "console"
	c.conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return c
}

// Option option interface
type Option interface {
	apply(*Config)
}

// option function
type optionFunc func(*Config)

// apply function
func (f optionFunc) apply(c *Config) {
	f(c)
}

// IsDevelopOption Set delopement option
func IsDevelopOption() Option {
	return optionFunc(func(c *Config) {
		c.conf.Development = true
	})
}

// LevelOption set log level
func LevelOption(level Level) Option {
	l := zap.DebugLevel
	zapL, exists := levelMap[level]
	if exists {
		l = zapL
	}

	return optionFunc(func(c *Config) {
		c.conf.Level = zap.NewAtomicLevelAt(l)
	})
}

// CallSkipOption skip caller
func CallSkipOption(stack int) Option {
	return optionFunc(func(c *Config) {
		c.zapOption = append(c.zapOption, zap.AddCallerSkip(stack))
	})
}

// get zap options
func (c *Config) getZapOptions() []zap.Option {
	return c.zapOption
}
