package log

import "go.uber.org/zap"

type Logger interface {
	WithPrefix(str string)
	Info(str string)
	Error(str string)
	Fatal(str string)
	Write(p []byte) (n int, err error)
}

type log struct {
	z *zap.Logger
}

func New() Logger {
	return log{z: &zap.Logger{}}
}

func (l log) WithPrefix(str string) {
	panic("implement me")
}

func (l log) Info(str string) {
	panic("implement me")
}

func (l log) Error(str string) {
	panic("implement me")
}

func (l log) Fatal(str string) {
	panic("implement me")
}

func (l log) Write(p []byte) (n int, err error) {
	panic("implement me")
}
