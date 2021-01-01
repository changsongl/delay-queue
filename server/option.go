package server

import (
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/vars"
)

type Option func(s *server)

func LoggerOption(l log.Logger) Option {
	return func(s *server) {
		s.l = l
	}
}

func EnvOption(env vars.Env) Option {
	return func(s *server) {
		s.env = env
	}
}

func BeforeStartEventOption(events ...Event) Option {
	return func(s *server) {
		s.beforeStart = append(s.beforeStart, events...)
	}
}

func AfterStopEventOption(events ...Event) Option {
	return func(s *server) {
		s.beforeStart = append(s.afterStop, events...)
	}
}
