package log

import (
	"go.uber.org/zap"
	"strings"
)

type Logger interface {
	WithModule(module string) Logger

	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Write(p []byte) (n int, err error)
}

type logger struct {
	z *zap.Logger
}

func New(opts ...Option) (Logger, error) {
	c := NewConfig()
	opts = append(opts, CallSkipOption(1))

	for _, opt := range opts {
		opt.apply(c)
	}

	l, err := c.conf.Build(c.getZapOptions()...)
	if err != nil {
		return nil, err
	}

	return &logger{z: l}, nil
}

func (l *logger) WithModule(module string) Logger {
	lClone := *l
	lClone.z = lClone.z.Named(module)
	return &lClone
}

func (l *logger) Debug(s string, fs ...Field) {
	l.z.Debug(s, getZapFields(fs...)...)
}

func (l *logger) Info(s string, fs ...Field) {
	l.z.Info(s, getZapFields(fs...)...)
}

func (l *logger) Warn(s string, fs ...Field) {
	l.z.Warn(s, getZapFields(fs...)...)
}

func (l *logger) Error(s string, fs ...Field) {
	l.z.Error(s, getZapFields(fs...)...)
}

func (l *logger) Fatal(s string, fs ...Field) {
	l.z.Fatal(s, getZapFields(fs...)...)
}

func (l *logger) Write(b []byte) (n int, err error) {
	l.Info(strings.TrimRight(string(b), "\n"))
	return len(b), nil
}
