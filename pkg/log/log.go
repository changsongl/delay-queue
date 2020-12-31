package log

import (
	"github.com/changsongl/delay-queue/vars"
	"go.uber.org/zap"
)

type Logger interface {
	WithModule(module string)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Write(p []byte) (n int, err error)
}

type log struct {
	z      *zap.Logger
	module string
}

func New(env vars.Env) (Logger, error) {
	var z *zap.Logger
	var err error

	if env == vars.EnvRelease {
		z, err = zap.NewProduction(zap.AddCallerSkip(1))
	} else {
		z, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	}

	if err != nil {
		return nil, err
	}

	return log{z: z}, nil
}

func (l log) WithModule(module string) {
	l.module = module
}

func (l log) Debug(msg string, fields ...Field) {
	fs := getZapFields(fields...)
	l.z.Debug(msg, fs...)
}

func (l log) Info(msg string, fields ...Field) {
	fs := getZapFields(fields...)
	l.z.Info(msg, fs...)
}

func (l log) Error(msg string, fields ...Field) {
	fs := getZapFields(fields...)
	l.z.Error(msg, fs...)
}

func (l log) Fatal(msg string, fields ...Field) {
	fs := getZapFields(fields...)
	l.z.Fatal(msg, fs...)
}

func (l log) Write(b []byte) (n int, err error) {
	l.Info(string(b))
	return len(b), nil
}
