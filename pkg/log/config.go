package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelMap = map[Level]zapcore.Level{
	LevelDebug: zap.DebugLevel,
	LevelInfo:  zapcore.InfoLevel,
	LevelWarn:  zapcore.WarnLevel,
	LevelError: zapcore.ErrorLevel,
	LevelFatal: zapcore.FatalLevel,
}

type Config struct {
	conf      zap.Config
	zapOption []zap.Option
}

func NewConfig() *Config {
	c := &Config{conf: zap.NewProductionConfig()}
	c.conf.Encoding = "console"
	c.conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return c
}

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

func IsDevelopOption() Option {
	return optionFunc(func(c *Config) {
		c.conf.Development = true
	})
}

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

func CallSkipOption(stack int) Option {
	return optionFunc(func(c *Config) {
		c.zapOption = append(c.zapOption, zap.AddCallerSkip(stack))
	})
}

func (c *Config) getZapOptions() []zap.Option {
	return c.zapOption
}
